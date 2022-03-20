package main

import (
	"context"
	"fmt"
	"github.com/StarsiegePlayers/DiscordBot/plugin"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGQUIT)
	defer stop()

	hub := plugin.NewIPCHub(ctx)
	go hub.Process(hubHandler)

	log.Println("[main] Running...")

	LoadPlugins(hub)

	// block until canceled
	<-ctx.Done()
	stop()

	fmt.Println("[main] graceful shutdown")
}

func hubHandler(message plugin.IPCMessage) {
	if message.Destination == plugin.HUB {
		log.Printf("[main-hub] %s", message)
	}
}
