package main

import (
	"fmt"

	"github.com/aboglioli/big-brother/composition"
	"github.com/aboglioli/big-brother/infrastructure/events"
)

func main() {
	eventMgr, err := events.GetManager()
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

		fmt.Println("[Waiting for events]")
		for msg := range msgs {
			evt, err := composition.EventFromBytes(msg.Body())
			if err != nil {
				fmt.Println(err)
				continue
			}

			fmt.Println("# New event:")
			fmt.Println(evt)

			msg.Ack()
		}
	}()

	<-forever
}
