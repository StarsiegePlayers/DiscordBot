package module

import (
	"context"
	"errors"
	"fmt"

	"github.com/StarsiegePlayers/DiscordBot/log"

	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

var (
	ErrUnimplementedFunction = errors.New("unimplemented function")
)

type Interface interface {
	BaseInit(context.Context, *gochannel.GoChannel, string)
	Init()
	Start() error
	Stop() error
}

type Base struct {
	Name           string
	pubsub         *gochannel.GoChannel
	PubSubChannels map[string]*PubSubInfo

	context.Context
	log.Log
	Interface
}

func (b *Base) BaseInit(ctx context.Context, pubsub *gochannel.GoChannel, name string) {
	b.Context = ctx
	b.Name = name
	b.Log = log.NewLogger(b.Name)

	b.pubsub = pubsub
	b.PubSubChannels = make(map[string]*PubSubInfo)
}

func (b *Base) Init() {}

func (b *Base) Start() error {
	return errors.Unwrap(fmt.Errorf("%w: Start()", ErrUnimplementedFunction))
}

func (b *Base) Stop() error {
	return errors.Unwrap(fmt.Errorf("%w: Stop()", ErrUnimplementedFunction))
}
