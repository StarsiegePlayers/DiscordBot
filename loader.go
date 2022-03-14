package main

import (
	"context"
	"log"

	"github.com/StarsiegePlayers/DiscordBot/plugin"
)

type pluginData struct {
	Filename   string
	Instance   *plugin.Interface
	IPC        plugin.IPCNode
	Context    context.Context
	CancelFunc context.CancelFunc
}

var plugins map[string]pluginData

func init() {
	plugins = make(map[string]pluginData)

	plugins["api"] = pluginData{
		Filename: "./plugins/api.so",
	}
	plugins["discord"] = pluginData{
		Filename: "./plugins/discord.so",
	}
	plugins["calendar"] = pluginData{
		Filename: "./plugins/calendar.so",
	}
	plugins["quickchat"] = pluginData{
		Filename: "./plugins/quickchat.so",
	}
}

func LoadPlugins() {
	for k, v := range plugins {
		instance, err := plugin.Load(v.Filename)
		if err != nil {
			log.Printf("error loading plugin [%s]: %s\n", k, err)
			continue
		}

		v.Instance = instance
		v.IPC.TX = make(chan plugin.IPCMessage, plugin.MIN_BUFFER)
		v.IPC.RX = (*instance).Init(v.IPC.TX)
		v.Context, v.CancelFunc = context.WithCancel(context.Background())
		(*instance).Attach(v.Context)
	}
}
