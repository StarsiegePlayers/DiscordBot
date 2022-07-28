package config

import (
	"fmt"
	"strings"

	"golang.org/x/exp/maps"
)

type DiscordConfig struct {
	DebugUsers    map[string]string      `yaml:"DebugUsers"`
	Guilds        map[string]GuildConfig `yaml:"Guilds"`
	AuthToken     string                 `yaml:"AuthToken"`
	ApplicationID string                 `yaml:"ApplicationID"`
}

type GuildConfig struct {
	BotLogChannelID          string            `yaml:"BotLogChannelID"`
	VoiceChannelID           string            `yaml:"VoiceChannelID"`
	CommandPrefix            string            `yaml:"CommandPrefix"`
	APIAnnouncementChannelID string            `yaml:"APIAnnouncementChannelID"`
	NamedRoles               map[string]string `yaml:"NamedRoles"`
	MuzzledUsers             map[string]int64  `yaml:"MuzzledUsers"`
	Webhooks                 map[string]string `yaml:"Webhooks"`
}

func (g GuildConfig) String() string {
	return fmt.Sprintf("log: %s | vc: %s | pretix: %s | announcement: %s | namedroles: {%s} | MuzzledUsers: {%s}", g.BotLogChannelID, g.VoiceChannelID, g.CommandPrefix, g.APIAnnouncementChannelID, strings.Join(maps.Keys(g.NamedRoles), ", "), strings.Join(maps.Keys(g.MuzzledUsers), ", "))
}

type CalendarConfig struct {
	CalendarID        string `yaml:"CalendarID"`
	NumEventLookAhead int64  `yaml:"NumEventLookAhead"`
	AuthToken         string `yaml:"AuthToken"`
}

type IRCConfig struct {
	Server    string `yaml:"Server"`
	Port      int    `yaml:"Port"`
	AuthToken string `yaml:"AuthToken"`
}

type APIConfig struct {
	URL             string `yaml:"URL"`
	PollTimeMinutes int    `yaml:"PollTimeMinutes"`
}

type Config struct {
	Discord  DiscordConfig  `yaml:"Discord"`
	Calendar CalendarConfig `yaml:"Calendar"`
	IRC      IRCConfig      `yaml:"IRC"`
	API      APIConfig      `yaml:"API"`
}
