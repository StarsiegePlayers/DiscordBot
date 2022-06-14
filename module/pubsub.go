package module

import (
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

type PubSubChannel <-chan *message.Message
type PubSubHandler func(*PubSubInfo, *message.Message) error

type PubSubInfo struct {
	Topic string
	*gochannel.GoChannel
	PubSubChannel
	PubSubHandler
	Cancel func()
}

func (b *Base) PubSubSubscribe(topic string, handlerFn PubSubHandler) {
	ctx, cancelFn := context.WithCancel(b)

	c, err := b.pubsub.Subscribe(ctx, topic)
	if err != nil {
		b.Logln(fmt.Errorf("error subscribing to: %s", topic))
	}

	ps := &PubSubInfo{
		Topic:         topic,
		GoChannel:     b.pubsub,
		PubSubChannel: c,
		PubSubHandler: handlerFn,
		Cancel:        cancelFn,
	}

	b.PubSubChannels[topic] = ps

	go func() {
		for {
			select {
			case msg := <-c:
				err = handlerFn(ps, msg)
				if err != nil {
					b.Logln(err)
				}

			case <-ctx.Done():
				return
			}
		}
	}()
}

func (b *Base) PubSubPublish(topic string, payload []byte) error {
	return b.pubsub.Publish(topic, message.NewMessage(watermill.NewUUID(), payload))
}
