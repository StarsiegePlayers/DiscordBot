package main

import (
	"context"
	"log"

	"github.com/StarsiegePlayers/DiscordBot/plugin"
)

var this = &API{
	plugin.Base{
		Name: "API",
	},
}

type API struct {
	plugin.Base
}

func (a *API) Attach(ctx context.Context) {
	a.RegisterCommand("ls", a.performLS)
	a.RegisterCommand("oof", a.performOof)

	a.Base.Attach(ctx)
}

func (a *API) performLS(m plugin.IPCMessage) {
	a.Logf("ls")
}

func (a *API) performOof(m plugin.IPCMessage) {
	log.Printf("[%s]: oof", a.Name)
}

func Export() plugin.Interface {
	return this
}
