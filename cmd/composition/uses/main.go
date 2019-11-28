package main

import (
	"fmt"
	"log"

	"github.com/aboglioli/big-brother/composition"
	"github.com/aboglioli/big-brother/errors"
	"github.com/aboglioli/big-brother/events"
	infrEvents "github.com/aboglioli/big-brother/infrastructure/events"
)

type Context struct {
	eventMgr events.Manager
	repo     composition.Repository
	serv     composition.Service
}

func (c *Context) UpdateUses(comp *composition.Composition) errors.Error {
	errGen := errors.NewInternal().SetPath("cmd/uses/main.Context.UpdateUses")

	fmt.Printf("# Updating uses of %s (%s): ", comp.Name, comp.ID.Hex())

	uses, err := c.serv.UpdateUses(comp)
	if err != nil {
		return errGen.SetCode("UPDATE_USES").SetMessage(err.Error())
	}
	fmt.Printf("updated %d dependencies\n", len(uses))

	for _, u := range uses {
		fmt.Printf("- %s (%s)\n", u.Name, u.ID.Hex())
	}

	// Update composition to set UsesUpdatedSinceLastChange
	comp.UsesUpdatedSinceLastChange = true
	if err := c.repo.Update(comp); err != nil {
		return errGen.SetCode("UPDATE_UsesUpdatedSinceLastChange").SetMessage(err.Error())
	}

	if err := c.Publish("CompositionUsesUpdatedSinceLastChange", comp); err != nil {
		return errGen.SetCode("PUBLISH_CompositionUsesUpdatedSinceLastChange").SetMessage(err.Error())
	}

	return nil
}

func (c *Context) Publish(eventType string, comp *composition.Composition) errors.Error {
	event := &composition.CompositionChangedEvent{events.Event{eventType}, comp}
	if err := c.eventMgr.Publish("composition", "topic", "composition.updated", event); err != nil {
		return err
	}

	return nil
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
			if msg.Type() == "CompositionUpdatedManually" {
				var event composition.CompositionChangedEvent
				if err := msg.Decode(&event); err != nil {
					fmt.Println(err)
					continue
				}
				comp := event.Composition

				if err := ctx.UpdateUses(comp); err != nil {
					fmt.Println(comp.ID.Hex(), err)
				}
			}
			msg.Ack()
		}
	}()

	<-forever
}
