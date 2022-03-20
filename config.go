package main

type DiscordConfig struct {
	AuthToken      string
	CommandPrefix  string
	GuildWhitelist []string
}

type GoogleAuthConfig struct {
	TokenData string
}

type Config struct {
	Plugins map[string]string

	DiscordCalendar struct {
		DiscordConfig
	}
	Calendar struct {
		CalendarURI string
		GoogleAuthConfig
	}

	DiscordIRC struct {
		DiscordConfig
	}
	IRC struct {
		Server      string
		Port        int
		WebIRCToken string
	}
}
