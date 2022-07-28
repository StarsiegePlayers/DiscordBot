package discord

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/png"

	"github.com/StarsiegePlayers/DiscordBot/config"

	"github.com/bwmarrin/discordgo"
)

type Session struct {
	*discordgo.Session

	GuildConfig config.GuildConfig
	service     *Service
}

func (s *Session) SetService(service *Service) {
	s.service = service
}

func (s *Session) CreateOrFetchChannelWebhook(guildID string, channelID string) (*discordgo.Webhook, error) {
	var hook *discordgo.Webhook
	guild, err := s.State.Guild(guildID)
	if err != nil {
		s.service.Logln("error fetching guild")
		return nil, err
	}

	channel, err := s.State.Channel(channelID)
	if err != nil {
		s.service.Logln("error fetching channel")
		return nil, err
	}

	if h, ok := s.GuildConfig.Webhooks[channelID]; ok {
		return s.Webhook(h)
	} else {
		img, err := s.UserAvatar(s.State.User.ID)
		if err != nil {
			s.service.Logln("error decoding User Avatar", err)
			return nil, err
		}

		avatarPng := new(bytes.Buffer)

		err = png.Encode(avatarPng, img)
		if err != nil {
			s.service.Logf("(%s) unable to encode default avatar: %s", guild.Name, err)
			return nil, err
		}

		avatarBase64 := fmt.Sprintf("%s%s", "data:image/png;base64,", base64.StdEncoding.EncodeToString(avatarPng.Bytes()))

		hook, err = s.WebhookCreate(channelID, s.State.User.Username, avatarBase64)
		if err != nil {
			s.service.Logf("(%s) unable to create webhook in %s: %s", guild.Name, channel.Name, err)
			return nil, err
		}

		s.GuildConfig.Webhooks[channelID] = hook.ID
		s.service.sendRPCConfigUpdate()
		return hook, err
	}
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
