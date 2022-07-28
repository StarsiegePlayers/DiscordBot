package discord

import (
	"encoding/json"
	"time"

	"github.com/StarsiegePlayers/DiscordBot/config"
	"github.com/StarsiegePlayers/DiscordBot/module"
	"github.com/StarsiegePlayers/DiscordBot/rpc"
	"github.com/ThreeDotsLabs/watermill/message"
	"gopkg.in/yaml.v3"
)

func (s *Service) sendRPCConfigUpdate() {
	cfg, err := yaml.Marshal(s.config)
	if err != nil {
		s.Logln("error marshalling config", err)
	}

	err = s.RPCPublish(rpc.ConfigUpdatedFromDiscord, cfg)
	if err != nil {
		s.Logln("error publishing config", err)
	}
}

func (s *Service) discordMessageSendRPCHandler(rpcInfo *module.RPCInfo, msg *message.Message) (err error) {
	defer msg.Ack()
	s.Logf("received %s", string(msg.Payload))

	return
}

func (s *Service) configMessageRPCHandler(rpcInfo *module.RPCInfo, msg *message.Message) (err error) {
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

	s.config = &c.Discord

	return
}

func (s *Service) apiResponseRPCHandler(rpcInfo *module.RPCInfo, msg *message.Message) (err error) {
	msg.Ack()

	s.Logln("received new api response")

	if s.session == nil {
		time.Sleep(time.Second * 10)
	}

	// if we're still in init while receiving this message, we can wait until connected to process it
	s.wg.Wait()

	var apiResponse []rpc.ServerListData

	err = json.Unmarshal(msg.Payload, &apiResponse)
	if err != nil {
		s.Logln("error unmarshaling apiResponse JSON", err)
		return
	}

	err = s.session.UpdateGameStatus(0, apiResponse[0].String())
	if err != nil {
		s.Logln("error updating bot status", err)
		return
	}

	return
}
