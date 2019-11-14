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
			evt, err := events.FromBytes(msg.Body())
			if err != nil {
				fmt.Println(err)
				continue
			}

			if evt.Type == "CompositionUpdatedManually" {
				comp, err := payloadToComposition(evt.Payload)
				if err != nil {
					fmt.Println(err)
					continue
				}

				fmt.Printf("# Updating uses of %s (%s): ", comp.Name, comp.ID.Hex())

				comps, err := compositionService.UpdateUses(comp)
				if err != nil {
					fmt.Printf("[ERROR] %s\n", err)
					continue
				}
				fmt.Printf("updated %d dependencies\n", len(comps))

				for _, c := range comps {
					fmt.Printf("> %+v\n", c)
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
