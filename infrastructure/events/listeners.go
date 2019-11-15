package events

import (
	"fmt"
	"sync"

	"github.com/aboglioli/big-brother/auth"
	"github.com/aboglioli/big-brother/composition"
	"github.com/aboglioli/big-brother/events"
)

func StartListeners(eventMgr events.Manager, serv composition.Service) {
	var wg sync.WaitGroup

	// Listen logout
	go func() {
		wg.Add(1)

		msgs, err := eventMgr.Consume("auth", "fanout", "", "")
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("[Listening for 'logout']")
		for msg := range msgs {
			t := &events.Type{}
			if err := t.FromBytes(msg.Body()); err != nil {
				fmt.Println(err)
				continue
			}

			if t.Type == "logout" {
				event := &auth.LogoutEvent{}
				if err := event.FromBytes(msg.Body()); err != nil {
					fmt.Println(err)
					continue
				}

				auth.Invalidate(event.Message)
			}
			msg.Ack()
		}

		wg.Done()
	}()

	// Listen article-exits
	go func() {
		wg.Add(1)

		msgs, err := eventMgr.Consume("auth", "direct", "", "")
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("[Listening for 'article-exists']")
		for msg := range msgs {
			t := &events.Type{}
			if err := t.FromBytes(msg.Body()); err != nil {
				fmt.Println(err)
				continue
			}

			if t.Type == "article-exist" {
				event := &composition.ArticleExistsEventResponse{}
				if err := event.FromBytes(msg.Body()); err != nil {
					fmt.Println(err)
					continue
				}

				if event.Message.Valid {
					serv.Validate(event.Message.ArticleID)
				}
			}
			msg.Ack()
		}

		wg.Done()
	}()

	wg.Wait()
}
