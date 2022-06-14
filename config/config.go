package config

import (
	"io"
	"os"

	"github.com/StarsiegePlayers/DiscordBot/module"
	"github.com/StarsiegePlayers/DiscordBot/rpc"
	
	"gopkg.in/yaml.v3"
)

const (
	ServiceName = "config"
	FileName    = "config.yaml"
)

type Service struct {
	module.Base

	config Config
}

func (s *Service) Init() {
	f, err := os.OpenFile(FileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	c, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(c, &s.config)
	if err != nil {
		s.Logln(err)
	}
}

func (s *Service) Start() error {
	c, err := yaml.Marshal(s.config)
	if err != nil {
		s.Logln(err)
		return err
	}

	err = s.PubSubPublish(rpc.NewConfigLoadedTopic, c)
	if err != nil {
		s.Logln(err)
		return err
	}

	return err
}

func (s *Service) Stop() error {
	return nil
}
