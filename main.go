package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/StarsiegePlayers/DiscordBot/plugin"
)

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())

	hub := plugin.NewIPCHub(ctx)
	go hub.Process(hubHandler)

	log.Println("Running...")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	cancelFunc()
	log.Println("Graceful shutdown")

}

func hubHandler(message plugin.IPCMessage) {

}
