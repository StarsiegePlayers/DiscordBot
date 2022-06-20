package config

import (
	"fmt"
	"strings"
)

type DiscordConfig struct {
	DebugUsers map[string]string      `yaml:"DebugUsers"`
	Guilds     map[string]GuildConfig `yaml:"Guilds"`
	AuthToken  string                 `yaml:"AuthToken"`
}

type GuildConfig struct {
	VoiceChannelID string                    `yaml:"VoiceChannelID"`
	CommandPrefix  string                    `yaml:"CommandPrefix"`
	TimeoutConfig  DiscordGuildTimeoutConfig `yaml:"TimeoutConfig"`
}

func (g GuildConfig) String() string {
	return fmt.Sprintf("vc: %s | pretix: %s | timeoutconfig: %s", g.VoiceChannelID, g.CommandPrefix, g.TimeoutConfig)
}

type DiscordGuildTimeoutConfig struct {
	TimeoutRoleID  string   `yaml:"TimeoutRoleID"`
	TimeoutTTL     int      `yaml:"TimeoutTTL"`
	ExemptChannels []string `yaml:"ExemptChannels"`
}

func (d DiscordGuildTimeoutConfig) String() string {
	return fmt.Sprintf("role: %s | ttl: %d | exempt: %s", d.TimeoutRoleID, d.TimeoutTTL, strings.Join(d.ExemptChannels, ", "))
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
