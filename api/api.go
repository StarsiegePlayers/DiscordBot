package api

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/StarsiegePlayers/DiscordBot/config"
	"github.com/StarsiegePlayers/DiscordBot/module"
	"github.com/StarsiegePlayers/DiscordBot/rpc"
)

const (
	ServiceName = "api"
)

type Service struct {
	module.Base

	wg     sync.WaitGroup
	config *config.APIConfig

	APIHistory []rpc.ServerListData
}

func (s *Service) Init() {
	// queue waiting on config
	s.wg.Add(1)

	s.APIHistory = make([]rpc.ServerListData, 0)
	s.RPCSubscribe(rpc.NewConfigLoadedTopic, s.configMessageRPCHandler)
}

func (s *Service) Start() error {
	// wait for config
	s.wg.Wait()

	s.alarmCallback(s.Base, func() {})
	s.RPCSubscribe(rpc.APIRequestLatest, s.apiRequestLatestRPCHandler)

	return nil
}

func (s *Service) Stop() error {
	return nil
}

func (s *Service) requestServerList() (list rpc.ServerListData, err error) {
	client := &http.Client{}
	client.Timeout = 5 * time.Second

	req, err := http.NewRequest(http.MethodGet, s.config.URL, nil)
	if err != nil {
		return
	}

	res, err := client.Do(req)
	if err != nil {
		return
	}

	defer res.Body.Close()

	lastpudated := res.Header.Get("Last-Modified")
	s.Logf("Last-Modified: %s", lastpudated)

	APILastUpdated, err := time.Parse(time.RFC1123, lastpudated)
	if err != nil {
		s.Logln("unable to parse last updated time, setting to now()")
		APILastUpdated = time.Now()
	}

	s.Base.NewAlarm(s.Base, APILastUpdated, time.Duration(s.config.PollTimeMinutes)*time.Minute, s.alarmCallback)

	body := json.NewDecoder(res.Body)

	err = body.Decode(&list)
	if err != nil {
		return
	}

	if len(s.APIHistory) > 0 {
		s.APIHistory = s.APIHistory[0:]
	}

	s.APIHistory = append(s.APIHistory, list)

	return
}

func (s *Service) alarmCallback(ctx context.Context, cancelfn context.CancelFunc) {
	list, err := s.requestServerList()
	if err != nil {
		s.Logln(err)
		return
	}

	s.Log.Println(list)
}
