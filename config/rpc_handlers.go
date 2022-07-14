package config

import (
	"io/ioutil"

	"github.com/StarsiegePlayers/DiscordBot/module"

	"github.com/ThreeDotsLabs/watermill/message"
	"gopkg.in/yaml.v3"
)

func (s *Service) discordConfigUpdate(rpcInfo *module.RPCInfo, msg *message.Message) (err error) {
	defer msg.Ack()

	c := new(DiscordConfig)

	err = yaml.Unmarshal(msg.Payload, c)
	if err != nil {
		s.Logln(err)
	}

	s.config.Discord = *c

	out, err := yaml.Marshal(s.config)
	if err != nil {
		s.Logln("unable to marshal config", err)
	}

	err = ioutil.WriteFile(FileName, out, 0644)
	if err != nil {
		s.Logln("unable to write config", err)
	}

	return
}
