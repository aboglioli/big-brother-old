package main

import (
	"encoding/json"
	"fmt"

	"github.com/aboglioli/big-brother/composition"
	"github.com/aboglioli/big-brother/events"
	infrEvents "github.com/aboglioli/big-brother/infrastructure/events"
)

func main() {
	eventMgr, err := infrEvents.GetManager()
	if err != nil {
		fmt.Println(err)
		return
	}

	forever := make(chan bool)

	go func() {
		msgs, err := eventMgr.Consume("composition", "topic", "composition.*")
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("[Waiting for events on topic: 'composition.*']")
		for msg := range msgs {
			evt, err := events.FromBytes(msg.Body())
			if err != nil {
				fmt.Println(err)
				continue
			}

			fmt.Println("# New event:")

			switch evt.Type {
			case "CompositionUpdatedManually":
				comp, err := payloadToComposition(evt.Payload)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Printf("- Type: %s\n- Payload %+v\n", evt.Type, comp)
			case "CompositionsUpdatedAutomatically":
				comps, err := payloadToCompositions(evt.Payload)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Printf("- Type: %s\n- Payload %+v\n", evt.Type, comps)
			default:
				fmt.Printf("- Type: %s\n- Payload %s\n", evt.Type, evt.Payload)
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

func payloadToCompositions(payload interface{}) ([]*composition.Composition, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	var comps []*composition.Composition
	err = json.Unmarshal(b, &comps)
	if err != nil {
		return nil, err
	}

	return comps, nil
}
