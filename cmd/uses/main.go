package main

import (
	"fmt"
	"log"

	"github.com/aboglioli/big-brother/composition"
	"github.com/aboglioli/big-brother/infrastructure/events"
)

func main() {
	// Dendencies resolution
	eventMgr, err := events.GetManager()
	if err != nil {
		log.Fatal(err)
	}

	compositionRepository, rawErr := composition.NewRepository()
	if rawErr != nil {
		log.Fatal(err)
	}

	compositionService := composition.NewService(compositionRepository, eventMgr)

	forever := make(chan bool)

	go func() {
		msgs, err := eventMgr.Consume("composition", "topic", "composition.updated")
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("[Listening for composition updates]")
		for msg := range msgs {
			evt, err := composition.EventFromBytes(msg.Body())
			if err != nil {
				fmt.Println(err)
				continue
			}

			if evt.Type == "CompositionUpdatedManually" {
				fmt.Printf("Updating uses of %s (%s)\n", evt.Composition.Name, evt.Composition.ID.Hex())

				err := compositionService.UpdateUses(evt.Composition)
				if err != nil {
					fmt.Println("# ERROR #")
					fmt.Println(err)
					continue
				}

				fmt.Println("Updated")
			}

			msg.Ack()
		}
	}()

	<-forever
}
