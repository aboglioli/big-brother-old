package main

import (
	"fmt"

	"github.com/aboglioli/big-brother/composition"
	infrEvents "github.com/aboglioli/big-brother/infrastructure/events"
	"github.com/aboglioli/big-brother/pkg/events"
)

func main() {
	eventMgr, err := infrEvents.GetManager()
	if err != nil {
		fmt.Println(err)
		return
	}

	forever := make(chan bool)

	go func() {
		opts := &events.Options{"composition", "topic", "composition.*", ""}
		msgs, err := eventMgr.Consume(opts)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("[Waiting for events on topic: 'composition.*']")
		for msg := range msgs {
			eventType := msg.Type()

			fmt.Printf("# NEW EVENT: %s\n", msg.Type())

			switch eventType {
			case "CompositionCreated", "CompositionUpdatedManually", "CompositionUsesUpdatedSinceLastChange":
				var event composition.CompositionChangedEvent
				if err := msg.Decode(&event); err != nil {
					fmt.Println(err)
					continue
				}
				comp := event.Composition
				fmt.Printf("- Composition: %s (%s)\n", comp.Name, comp.ID.Hex())
			case "CompositionsUpdatedAutomatically":
				var event composition.CompositionsUpdatedAutomaticallyEvent
				if err := msg.Decode(&event); err != nil {
					fmt.Println(err)
					continue
				}
				comps := event.Compositions
				fmt.Println("- Compositions:")
				for _, c := range comps {
					fmt.Printf("-- %s (%s)\n", c.Name, c.ID.Hex())
				}
			default:
				fmt.Println("- Unknown event")
			}

			msg.Ack()
		}
	}()

	<-forever
}
