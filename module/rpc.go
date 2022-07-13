package module

import (
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

type RPCChannel <-chan *message.Message
type RPCHandler func(*RPCInfo, *message.Message) error

type RPCInfo struct {
	Topic string
	*gochannel.GoChannel
	RPCChannel
	RPCHandler
	Cancel func()
}

func (b *Base) RPCSubscribe(topic string, handlerFn RPCHandler) {
	ctx, cancelFn := context.WithCancel(b)

	c, err := b.rpc.Subscribe(ctx, topic)
	if err != nil {
		b.Logln(fmt.Errorf("error subscribing to: %s", topic))
	}

	ps := &RPCInfo{
		Topic:      topic,
		GoChannel:  b.rpc,
		RPCChannel: c,
		RPCHandler: handlerFn,
		Cancel:     cancelFn,
	}

	b.RPCChannels[topic] = ps

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
	return b.rpc.Publish(topic, message.NewMessage(watermill.NewUUID(), payload))
}
