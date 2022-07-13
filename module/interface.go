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
	Name        string
	rpc         *gochannel.GoChannel
	RPCChannels map[string]*RPCInfo

	context.Context
	log.Log
	Interface
}

func (b *Base) BaseInit(ctx context.Context, rpc *gochannel.GoChannel, name string) {
	b.Context = ctx
	b.Name = name
	b.Log = log.NewLogger(b.Name)

	b.rpc = rpc
	b.RPCChannels = make(map[string]*RPCInfo)
}

func (b *Base) Init() {}

func (b *Base) Start() error {
	return errors.Unwrap(fmt.Errorf("%w: Start()", ErrUnimplementedFunction))
}

func (b *Base) Stop() error {
	return errors.Unwrap(fmt.Errorf("%w: Stop()", ErrUnimplementedFunction))
}
