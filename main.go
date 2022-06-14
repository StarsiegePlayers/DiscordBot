package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/StarsiegePlayers/DiscordBot/api"
	"github.com/StarsiegePlayers/DiscordBot/calendar"
	"github.com/StarsiegePlayers/DiscordBot/config"
	"github.com/StarsiegePlayers/DiscordBot/discord"
	logger "github.com/StarsiegePlayers/DiscordBot/log"
	"github.com/StarsiegePlayers/DiscordBot/module"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

func main() {
	log := logger.NewLogger("main")
	log.Println("startup")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGQUIT)
	defer stop()

	pubsub := gochannel.NewGoChannel(gochannel.Config{BlockPublishUntilSubscriberAck: true}, watermill.NopLogger{})

	services := make(map[string]module.Interface)
	services[api.ServiceName] = new(api.Service)
	services[discord.ServiceName] = new(discord.Service)
	services[calendar.ServiceName] = new(calendar.Service)

	for name, service := range services {
		log.Logf("init %s", name)
		service.BaseInit(ctx, pubsub, name)
		service.Init()
	}

	c := new(config.Service)
	c.BaseInit(ctx, pubsub, config.ServiceName)
	c.Init()
	_ = c.Start()

	for name, service := range services {
		log.Logf("Starting %s", name)
		if err := service.Start(); err != nil {
			panic(err)
		}
	}

	services[config.ServiceName] = c

	log.Println("running...")

	// block until done
	<-ctx.Done()
	log.Println("signal caught - shutting down...")
	stop()

	for _, m := range services {
		err := m.Stop()
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("graceful shutdown")
}
