package api

import (
	"encoding/json"
	"fmt"

	"github.com/StarsiegePlayers/DiscordBot/config"
	"github.com/StarsiegePlayers/DiscordBot/module"
	"github.com/StarsiegePlayers/DiscordBot/rpc"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"gopkg.in/yaml.v3"
)

func (s *Service) apiRequestLatestPubSubHandler(pubsub *module.PubSubInfo, msg *message.Message) (err error) {
	defer msg.Ack()

	if len(s.APIHistory) <= 0 {
		return
	}

	e, err := json.Marshal(s.APIHistory[len(s.APIHistory)-1])
	if err != nil {
		return fmt.Errorf("error marshalling APIHistory response %#v", err)
	}

	err = pubsub.Publish(rpc.APIRequestResponse, message.NewMessage(watermill.NewUUID(), e))
	if err != nil {
		return fmt.Errorf("error sending APIMessage %#v", err)
	}
	return
}

func (s *Service) configMessagePubSubHandler(pubsub *module.PubSubInfo, msg *message.Message) (err error) {
	defer msg.Ack()
	if s.config == nil {
		defer s.wg.Done()
	}

	s.Logln("received config")

	c := new(config.Config)

	err = yaml.Unmarshal(msg.Payload, c)
	if err != nil {
		s.Logln(err)
	}

	s.config = &c.API

	return
}
