package main

import (
	"fmt"
	"log"

	"github.com/aboglioli/big-brother/composition"
	implEvents "github.com/aboglioli/big-brother/impl/events"
	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/events"
)

type Context struct {
	eventMgr events.Manager
	repo     composition.Repository
	serv     composition.Service
}

func (c *Context) UpdateUses(comp *composition.Composition) error {
	path := "cmd/uses/main.Context.UpdateUses"

	fmt.Printf("# Updating uses of %s (%s): ", comp.Name, comp.ID.Hex())

	uses, err := c.serv.UpdateUses(comp)
	if err != nil {
		return errors.NewInternal("UPDATE_USES").SetPath(path).SetRef(err)
	}
	fmt.Printf("updated %d dependencies\n", len(uses))

	for _, u := range uses {
		fmt.Printf("- %s (%s)\n", u.Name, u.ID.Hex())
	}

	return nil
}

func main() {
	// Dendencies resolution
	eventMgr, err := implEvents.Rabbit()
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
		opts := &events.Options{"composition", "topic", "composition.updated", "uses"}
		msgs, err := eventMgr.Consume(opts)
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
