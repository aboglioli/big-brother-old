package main

import (
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
		msgs, err := eventMgr.Consume("composition", "topic", "", "composition.*")
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("[Waiting for events on topic: 'composition.*']")
		for msg := range msgs {
			t := &events.Type{}
			if err := t.FromBytes(msg.Body()); err != nil {
				fmt.Println(err)
				continue
			}

			fmt.Printf("# NEW EVENT: %s\n", t.Type)

			switch t.Type {
			case "CompositionCreated", "CompositionUpdatedManually", "CompositionUsesUpdatedSinceLastChange":
				event := &composition.CompositionChangedEvent{}
				if err := event.FromBytes(msg.Body()); err != nil {

				}
				comp := event.Composition
				fmt.Printf("- Composition: %s (%s)\n", comp.Name, comp.ID.Hex())
			case "CompositionsUpdatedAutomatically":
				event := &composition.CompositionUpdatedAutomaticallyEvent{}
				if err := event.FromBytes(msg.Body()); err != nil {
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
