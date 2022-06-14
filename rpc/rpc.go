package rpc

import (
	"fmt"
	"time"
)

const (
	APIRequestLatest        = "api.request.latest"
	APIRequestResponse      = "api.request.response"
	NewConfigLoadedTopic    = "config.NewConfigLoaded"
	DiscordMessageSendTopic = "discord.message.send"
)

type ServerListMaster struct {
	Address     string `json:"Address"`
	CommonName  string `json:"CommonName"`
	Motd        string `json:"MOTD"`
	ServerCount int    `json:"ServerCount"`
	Ping        int    `json:"Ping"`
}

type ServerListGame struct {
	GameMode    int                  `json:"GameMode"`
	GameName    string               `json:"GameName"`
	GameVersion string               `json:"GameVersion"`
	GameStatus  ServerListGameStatus `json:"GameStatus"`
	PlayerCount int                  `json:"PlayerCount"`
	MaxPlayers  int                  `json:"MaxPlayers"`
	Name        string               `json:"Name"`
	Address     string               `json:"Address"`
	Ping        int                  `json:"Ping"`
}

type ServerListGameStatus struct {
	Protected       bool `json:"Protected"`
	Dedicated       bool `json:"Dedicated"`
	AllowOldClients bool `json:"AllowOldClients"`
	Started         bool `json:"Started"`
	Dynamix         bool `json:"Dynamix"`
	Won             bool `json:"WON"`
	Unknown2        bool `json:"Unknown2"`
	Unknown3        bool `json:"Unknown3"`
}

func (s ServerListGame) String() string {
	return fmt.Sprintf("[%d/%d] %s (starsiege://%s) [%s]", s.PlayerCount, s.MaxPlayers, s.Name, s.Address, time.Duration(s.Ping))
}

type ServerListData struct {
	RequestTime time.Time          `json:"RequestTime"`
	Masters     []ServerListMaster `json:"Masters"`
	Games       []ServerListGame   `json:"Games"`
	Errors      []string           `json:"Errors"`
}

func (s ServerListData) String() string {
	activeGames := 0
	playerCount := 0

	for _, v := range s.Games {
		if v.PlayerCount > 0 {
			activeGames++

			playerCount += v.PlayerCount
		}
	}

	return fmt.Sprintf("[%d/%d] active games [%d player(s)]", activeGames, len(s.Games), playerCount)
}

func (s ServerListData) GetActiveGames() (out []ServerListGame) {
	for _, v := range s.Games {
		if v.PlayerCount > 0 {
			out = append(out, v)
		}
	}

	return
}
