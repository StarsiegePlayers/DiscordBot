package discord

import (
	"fmt"

	"github.com/StarsiegePlayers/DiscordBot/config"

	"github.com/bwmarrin/discordgo"
)

type Session struct {
	GuildConfig config.GuildConfig
	*discordgo.Session
}

func (d *Session) ChannelMessageMentionSend(channelID string, user *discordgo.User, message string) (*discordgo.Message, error) {
	return d.ChannelMessageSend(channelID, fmt.Sprintf("%s, %s", user.Mention(), message))
}

type MessageCreate struct {
	Guild       *discordgo.Guild
	Channel     *discordgo.Channel
	Member      *discordgo.Member
	Permissions int64
	Command     Command
	*discordgo.MessageCreate
}

type MessageHandler func(*Session, *MessageCreate, string)

type Command struct {
	Name       string
	Handler    MessageHandler
	Summary    string
	Usage      string
	Roles      []string
	Permission int64
}

type Commands map[string]Command

func (c *Commands) Register(cmd Command) {
	(*c)[cmd.Name] = cmd
}
