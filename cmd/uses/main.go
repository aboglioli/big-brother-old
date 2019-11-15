package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aboglioli/big-brother/composition"
	"github.com/aboglioli/big-brother/events"
	infrEvents "github.com/aboglioli/big-brother/infrastructure/events"
)

type Context struct {
	eventMgr events.Manager
	repo     composition.Repository
	serv     composition.Service
}

func (c *Context) UpdateUses(comp *composition.Composition) {
	fmt.Printf("# Updating uses of %s (%s): ", comp.Name, comp.ID.Hex())

	uses, err := c.serv.UpdateUses(comp)
	if err != nil {
		fmt.Printf("[ERROR] %s\n", err)
	}
	fmt.Printf("updated %d dependencies\n", len(uses))

	for _, u := range uses {
		fmt.Printf("- %s (%s)\n", u.Name, u.ID.Hex())
	}

	// Update composition to set UsesUpdatedSinceLastChange
	comp.UsesUpdatedSinceLastChange = true
	if err := c.repo.Update(comp); err != nil {
		fmt.Println(err)
	}

	c.Publish("CompositionUsesUpdatedSinceLastChange", comp)
}

func (c *Context) UpdateUsesSinceLastChange() {
	// Update dependencies from last changes
	comps, err := c.repo.FindByUsesUpdatedSinceLastChange(false)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, comp := range comps {
		c.UpdateUses(comp)
	}
}

func (c *Context) Publish(event string, comp *composition.Composition) {
	evt := events.NewEvent(event, comp)
	body, err := evt.ToBytes()
	if err != nil {
		fmt.Println(err)
	}
	if err := c.eventMgr.Publish("composition", "topic", "composition.updated", body); err != nil {
		fmt.Println(err)
	}
}

func main() {
	// Dendencies resolution
	eventMgr, err := infrEvents.GetManager()
	if err != nil {
		log.Fatal(err)
	}

	compositionRepository, err := composition.NewRepository()
	if err != nil {
		log.Fatal(err)
	}

	compositionService := composition.NewService(compositionRepository, eventMgr)

	ctx := &Context{
		eventMgr: eventMgr,
		repo:     compositionRepository,
		serv:     compositionService,
	}

	forever := make(chan bool)

	go func() {
		msgs, err := eventMgr.Consume("composition", "topic", "uses", "composition.updated")
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("[Listening for composition updates]")
		for msg := range msgs {
			evt, err := events.FromBytes(msg.Body())
			if err != nil {
				fmt.Println(err)
			}

			if evt.Type == "CompositionUpdatedManually" {
				comp, err := payloadToComposition(evt.Payload)
				if err != nil {
					fmt.Println(err)
				}

				ctx.UpdateUses(comp)
			}
			msg.Ack()
		}
	}()

	<-forever
}

func payloadToComposition(payload interface{}) (*composition.Composition, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	var comp composition.Composition
	err = json.Unmarshal(b, &comp)
	if err != nil {
		return nil, err
	}

	return &comp, nil
}
