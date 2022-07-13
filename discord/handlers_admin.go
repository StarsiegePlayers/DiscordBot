package discord

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/png"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

func (s *Service) registerAdminHandlers() {
	s.commands.Register(Command{
		Name:       "init",
		Handler:    s.mixinPermissionCheck(s.initHandler),
		Permission: discordgo.PermissionAdministrator,
	})
	s.commands.Register(Command{
		Name:    "move",
		Handler: s.mixinRoleCheck(s.moveHandler),
		Roles:   []string{"Staff"},
	})
	s.commands.Register(Command{
		Name:    "timeout",
		Handler: s.mixinRoleCheck(s.muzzleHandler),
		Roles:   []string{"Staff"},
	})
	s.commands.Register(Command{
		Name:    "resetwebhooks",
		Handler: s.mixinRoleCheck(s.resetWebhooksHandler),
		Roles:   []string{"Staff"},
	})
}

func (s *Service) initHandler(d *Session, m *MessageCreate, payload string) {

}

func (s *Service) resetWebhooksHandler(d *Session, m *MessageCreate, payload string) {
	channels, err := d.GuildChannels(m.GuildID)
	if err != nil {
		s.Logln("error fetching channels:", err)
		return
	}

	for _, channel := range channels {
		webhooks, err := d.ChannelWebhooks(channel.ID)
		if err != nil {
			s.Logf("error fetching webhooks for %s", channel.Name, err)
			continue
		}

		for _, w := range webhooks {
			if w.ApplicationID == s.config.ApplicationID {
				err = d.WebhookDelete(w.ID)
				if err != nil {
					s.Logln("error deleting webhook:", err)
					continue
				}
			}
		}
	}

	_, err = d.ChannelMessageSend(m.ChannelID, "Webhooks belonging to this bot have been deleted.")
	if err != nil {
		s.Logln("error while sending confirmation message:", err)
	}
}

func (s *Service) muzzleHandler(d *Session, m *MessageCreate, payload string) {

}

func (s *Service) moveHandler(d *Session, m *MessageCreate, payload string) {
	if len(payload) == 0 {
		_, err := s.session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s, Usage: in a reply to the message you wish to move, `!move <channel>`", m.Author.Mention()))
		if err != nil {
			s.Logln("error sending syntax message:", err)
			return
		}
		return
	}

	if m.ReferencedMessage == nil {
		_, err := s.session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s, you must reference or reply to a message to move it.", m.Author.Mention()))
		if err != nil {
			s.Logln("error sending error message:", err)
		}
		return
	}

	channel := payload
	if len(channel) <= 3 || channel[0:2] != "<#" || channel[len(channel)-1:] != ">" {
		_, err := s.session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s, you must specify a channel to move the message to.", m.Author.Mention()))
		if err != nil {
			s.Logln("error sending error message:", err)
		}
		return
	}

	channel = channel[2 : len(channel)-1]
	prevMessage := m.Message.ReferencedMessage

	member, err := d.GuildMember(m.GuildID, prevMessage.Author.ID)
	if err != nil {
		s.Logf("Unable to find guild member %s", prevMessage.Author.ID)
		return
	}

	var hook *discordgo.Webhook

	webhooks, err := d.ChannelWebhooks(channel)
	if err != nil {
		s.Logln("error fetching webhooks:", err)
	}

	for _, w := range webhooks {
		if w.ApplicationID == s.config.ApplicationID {
			hook = w
		}
	}

	if hook == nil {
		avatarImg, err := d.UserAvatarDecode(d.State.User)
		if err != nil {
			s.Logf("Unable to decode avatar for %s", prevMessage.Author.ID)
			return
		}

		avatarPng := new(bytes.Buffer)

		err = png.Encode(avatarPng, avatarImg)
		if err != nil {
			s.Logf("Unable to encode avatar for %s", prevMessage.Author.ID)
			return
		}

		avatarBase64 := fmt.Sprintf("%s%s", "data:image/png;base64,", base64.StdEncoding.EncodeToString(avatarPng.Bytes()))

		hook, err = d.WebhookCreate(channel, d.State.User.Username, avatarBase64)
		if err != nil {
			s.Logf("Unable to create webhook in %s (%s)", channel, err)
			return
		}
	}

	nick := member.Nick
	if nick == "" {
		nick = member.User.Username
	}

	var files []*discordgo.File
	if len(prevMessage.Attachments) > 0 {
		for _, v := range prevMessage.Attachments {
			body, err := http.Get(v.URL)
			if err != nil {
				s.Logln("error getting attachment:", err)
				continue
			}

			files = append(files, &discordgo.File{
				Name:        v.Filename,
				ContentType: v.ContentType,
				Reader:      body.Body,
			})
		}
	}

	_, err = d.WebhookExecute(hook.ID, hook.Token, true, &discordgo.WebhookParams{
		Content:    prevMessage.Content,
		Username:   nick,
		Files:      files,
		Components: prevMessage.Components,
		Embeds:     prevMessage.Embeds,
	})
	if err != nil {
		s.Logln("unable to execute webhook?!", err)
		return
	}

	err = d.ChannelMessageDelete(prevMessage.ChannelID, prevMessage.ID)
	if err != nil {
		s.Logln("unable to delete previous message?!", err)
		// do not return, we want to delete the message
	}

	err = d.ChannelMessageDelete(m.ChannelID, m.Message.ID)
	if err != nil {
		s.Logln("unable to delete trigger message?!", err)
		return
	}
}
