package main

import (
	"fmt"

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
			fmt.Printf("- Type: %s\n- Payload %s\n", evt.Type, evt.Payload)

			msg.Ack()
		}
	}()

	<-forever
}
