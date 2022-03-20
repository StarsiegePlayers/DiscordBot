package main

import (
	"context"
	"github.com/StarsiegePlayers/DiscordBot/plugin"
)

var this = &Discord{
	plugin.Base{
		Name: "Discord",
	},
}

type Discord struct {
	plugin.Base
}

func (d *Discord) Init(rxIn chan plugin.IPCMessage) chan plugin.IPCMessage {
	d.IPCNode.RX = rxIn
	d.IPCNode.TX = make(chan plugin.IPCMessage, plugin.MIN_BUFFER)

	return d.IPCNode.TX
}

func (d *Discord) Attach(ctx context.Context) {
	d.Base.Attach(ctx)
}

func Export() plugin.Interface {
	return this
}
