package config

type DiscordConfig struct {
	DebugUsers map[string]string
	Guilds     map[string]GuildConfig
	AuthToken  string
}

type GuildConfig struct {
	VoiceChannelID string
	CommandPrefix  string
	TimeoutConfig  DiscordGuildTimeoutConfig
}

type DiscordGuildTimeoutConfig struct {
	TimeoutRoleID  string
	TimeoutTTL     int
	ExemptChannels []string
}

type CalendarConfig struct {
	CalendarID        string
	NumEventLookAhead int64
	AuthToken         string
}

type IRCConfig struct {
	Server    string
	Port      int
	AuthToken string
}

type APIConfig struct {
	URL             string
	PollTimeMinutes int
}

type Config struct {
	Discord  DiscordConfig
	Calendar CalendarConfig
	IRC      IRCConfig
	API      APIConfig
}
