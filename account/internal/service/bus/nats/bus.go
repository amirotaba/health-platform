package nats

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"

	"git.paygear.ir/giftino/account/internal/account/event"
)

type bus struct {
	conn   *nats.Conn
	topics []string
}

func New(conn *nats.Conn) event.Bus {
	return &bus{
		conn:   conn,
		topics: make([]string, 0),
	}
}

func (b *bus) Publish(events []event.Event) error {
	// Simple Publisher for publish new event in distributed
	//c, err := nats.NewEncodedConn(b.conn, nats.JSON_ENCODER)
	//if err != nil {
	//	return err
	//}

	for _, evt := range events {
		switch e := evt.(type) {
		default:
			_ = e
			//err = c.Publish(evt.GetAggregateType(), evt)
			//if err != nil {
			//	return err
			//}

		}
	}
	return nil
}

func (b *bus) Register(topic string, handler func(event event.BaseEvent)) error {
	// Use a WaitGroup to wait for a message to arrive
	//wg := sync.WaitGroup{}
	//wg.Add(1)

	// Subscribe
	if _, err := b.conn.Subscribe(topic, func(m *nats.Msg) {
		var baseEvent event.BaseEvent
		err := json.Unmarshal(m.Data, &baseEvent)
		if err != nil {
			log.Println(err)
			return
		}
		handler(baseEvent)
		//wg.Done()
	}); err != nil {
		return err
	}

	// Wait for a message to come in
	//wg.Wait()
	fmt.Println(topic, "  registered")
	return nil
}

func (b *bus) ListTopic() []string {
	return b.topics
}
