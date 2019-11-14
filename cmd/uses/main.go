package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aboglioli/big-brother/composition"
	"github.com/aboglioli/big-brother/events"
	infrEvents "github.com/aboglioli/big-brother/infrastructure/events"
)

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

	forever := make(chan bool)

	go func() {
		msgs, err := eventMgr.Consume("composition", "topic", "composition.updated")
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

				fmt.Printf("# Updating uses of %s (%s): ", comp.Name, comp.ID.Hex())

				comps, err := compositionService.UpdateUses(comp)
				if err != nil {
					fmt.Printf("[ERROR] %s\n", err)
				}
				fmt.Printf("updated %d dependencies\n", len(comps))

				for _, c := range comps {
					fmt.Printf("- %s (%s)\n", c.Name, c.ID.Hex())
				}

				// Update composition to set UsesUpdatedSinceLastChange
				comp.UsesUpdatedSinceLastChange = true
				if err := compositionRepository.Update(comp); err != nil {
					fmt.Println(err)
				}

				evt := events.NewEvent("CompositionUsesUpdatedSinceLastChange", comp)
				body, err := evt.ToBytes()
				if err != nil {
					fmt.Println(err)
				}
				if err := eventMgr.Publish("composition", "topic", "composition.updated", body); err != nil {
					fmt.Println(err)
				}
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
