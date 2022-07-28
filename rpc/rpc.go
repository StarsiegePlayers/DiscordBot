package rpc

import (
	"bytes"
	"fmt"
	"image/color"
	"time"

	"gopkg.in/fogleman/gg.v1"
)

const (
	APIResponse              = "api.response"
	NewConfigLoadedTopic     = "config.NewConfigLoaded"
	ConfigUpdatedFromDiscord = "config.Update.Discord"
	DiscordMessageSendTopic  = "discord.message.send"
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

type RGBA struct {
	color.RGBA
}

func (c *RGBA) RGBAFloat() (float64, float64, float64, float64) {
	return float64(c.R) / 255, float64(c.G) / 255, float64(c.B) / 255, float64(c.A) / 255
}

var (
	Border     = RGBA{color.RGBA{R: 0x00, G: 0x39, B: 0x77, A: 0xff}}
	Blue       = RGBA{color.RGBA{R: 0x18, G: 0x30, B: 0x3f, A: 0xff}}
	Yellow     = RGBA{color.RGBA{R: 0x42, G: 0x36, B: 0x04, A: 0xff}}
	TextGreen  = RGBA{color.RGBA{R: 0x0a, G: 0x73, B: 0x00, A: 0xff}}
	TextYellow = RGBA{color.RGBA{R: 0xd5, G: 0xab, B: 0x00, A: 0xff}}
	NavBlue    = RGBA{color.RGBA{R: 0x18, G: 0x30, B: 0x3f, A: 0x80}}
	NavYellow  = RGBA{color.RGBA{R: 0x42, G: 0x36, B: 0x04, A: 0x80}}
)

func (s ServerListData) GetImage() interface{} {
	dc := gg.NewContext(800, 600)
	dc.SetRGBA(0x00, 0x00, 0x00, 0xff)
	dc.Fill()

	dc.SetRGBA(Border.RGBAFloat())
	dc.DrawRectangle(0, 0, 800, 600)
	dc.SetLineWidth(4)
	dc.Stroke()

	output := new(bytes.Buffer)
	_ = dc.SavePNG("out.png")
	_ = dc.EncodePNG(output)

	return output
}
