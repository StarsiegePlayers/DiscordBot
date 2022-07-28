package discord

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/StarsiegePlayers/DiscordBot/config"

	"github.com/bwmarrin/discordgo"
)

func (s *Service) registerAdminHandlers() {
	s.commands.Register(Command{
		Name:       "init",
		Handler:    s.mixinPermissionCheck(s.initHandler),
		Summary:    "Initialize the bot. This command is only available to administrators.",
		Usage:      "init",
		Permission: discordgo.PermissionAdministrator,
	})
	s.commands.Register(Command{
		Name:    "move",
		Handler: s.mixinRoleCheck(s.moveHandler),
		Summary: "Move a channel message to a different channel.",
		Usage:   "move <#channel> (must be in a reply to a message)",
		Roles:   []string{"Staff"},
	})
	s.commands.Register(Command{
		Name:    "massmove",
		Handler: s.mixinRoleCheck(s.moveHandler),
		Summary: "Move a channel message to a different channel.",
		Usage:   "move <#channel> <message ID> <message ID...>",
		Roles:   []string{"Staff"},
	})
	s.commands.Register(Command{
		Name:    "timeout",
		Handler: s.mixinRoleCheck(s.muzzleHandler),
		Summary: "Muzzle a user for a specified duration.",
		Usage:   "timeout <@user> <duration h|m|s>",
		Roles:   []string{"Staff"},
	})
	s.commands.Register(Command{
		Name:    "resetwebhooks",
		Handler: s.mixinRoleCheck(s.resetWebhooksHandler),
		Summary: "Reset all webhooks belonging to this bot.",
		Usage:   "resetwebhooks",
		Roles:   []string{"Staff"},
	})
}

func (s *Service) initHandler(d *Session, m *MessageCreate, payload string) {

}

func (s *Service) resetWebhooksHandler(d *Session, m *MessageCreate, payload string) {
	webhooks, err := d.GuildWebhooks(m.GuildID)
	if err != nil {
		s.Logln("(%s) error fetching webhooks: %s", m.Guild.Name, err)

		_, err = s.ChannelMessageMentionSend(m.ChannelID, m.Author, "Error: Webhooks for this bot have not been deleted.")
		if err != nil {
			s.Logf("(%s) error while sending confirmation message: %s", m.Guild.Name, err)
		}
	}

	for _, hook := range webhooks {
		if hook.ApplicationID == s.config.ApplicationID {
			err = d.WebhookDelete(hook.ID)
			if err != nil {
				s.Logln("(%s) error deleting webhook:", m.Guild.Name, err)
				continue
			}
		}
	}

	_, err = s.ChannelMessageMentionSend(m.ChannelID, m.Author, "Webhooks belonging to this bot have been deleted.")
	if err != nil {
		s.Logf("(%s) error while sending confirmation message: %s", m.Guild.Name, err)
	}
}

func (s *Service) muzzleHandler(d *Session, m *MessageCreate, payload string) {
	if _, ok := d.GuildConfig.NamedRoles["Muzzle"]; !ok {
		_, err := s.ChannelMessageMentionSend(m.ChannelID, m.Author, "Error: No timeout role has been set.")
		if err != nil {
			s.Logln("(%s) error while sending error message: %s", m.Guild.Name, err)
		}
		return
	}

	if len(payload) == 0 {
		s.sendUsageMessage(d, m)
		return
	}

	args := strings.Split(payload, " ")
	if len(args) != 2 {
		s.sendUsageMessage(d, m)
		return
	}

	userID := args[0]
	userID = userID[2 : len(userID)-1]

	member, err := s.session.GuildMember(m.GuildID, userID)
	if err != nil {
		s.Logf("(%s) error fetching member: %s", m.Guild.Name, err)

		_, err := s.ChannelMessageMentionSend(m.ChannelID, m.Author, fmt.Sprintf("Error: %s is not a valid user ID.", userID))
		if err != nil {
			s.Logf("(%s) error while sending error message: %s", m.Guild.Name, err)
		}

		return
	}

	nick := member.Nick
	if nick == "" {
		nick = member.User.Username
	}

	duration, err := time.ParseDuration(args[1])
	if err != nil {
		s.sendUsageMessage(d, m)
		return
	}

	if _, ok := d.GuildConfig.MuzzledUsers[member.User.ID]; ok {
		_, err := s.ChannelMessageMentionSend(m.ChannelID, m.Author, fmt.Sprintf("Error: %s is already muzzled.", nick))
		if err != nil {
			s.Logf("(%s) error while sending error message: %s", m.Guild.Name, err)
		}
		return
	}

	d.GuildConfig.MuzzledUsers[member.User.ID] = time.Now().Add(duration).Unix()

	err = s.session.GuildMemberRoleAdd(m.GuildID, member.User.ID, d.GuildConfig.NamedRoles["Muzzle"])
	if err != nil {
		s.Logf("(%s) error adding role: %s", m.Guild.Name, err)
		return
	}

	s.sendRPCConfigUpdate()

	_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s has been timed out until %s", nick, time.Unix(d.GuildConfig.MuzzledUsers[member.User.ID], 0).Format(time.RFC822)))
	if err != nil {
		s.Logf("(%s) error while sending confirmation message: %s", m.Guild.Name, err)
	}
}

func (s *Service) muzzleMaintenance(guildID string, cfg config.GuildConfig) {
	sendUpdate := false

	g, err := s.session.State.Guild(guildID)
	if err != nil {
		s.Logln("error fetching guild info:", err)
	}

	for user, t := range cfg.MuzzledUsers {
		if time.Now().After(time.Unix(t, 0)) {
			delete(cfg.MuzzledUsers, user)
			if muzzleID, ok := cfg.NamedRoles["Muzzle"]; ok {
				m, err := s.session.GuildMember(guildID, user)
				if err != nil {
					s.Logf("(%s) error fetching member info:", g.Name, err)
				}

				nick := m.Nick
				if nick == "" {
					nick = m.User.Username
				}

				err = s.session.GuildMemberRoleRemove(guildID, user, muzzleID)
				if err != nil {
					s.Logf("(%s) error removing muzzle role:", g.Name, err)
				}

				s.Logf("(%s) removed muzzle role from %s", g.Name, nick)
			}
			sendUpdate = true
		}
	}

	if sendUpdate {
		s.sendRPCConfigUpdate()
	}
}

func (s *Service) moveHandler(d *Session, m *MessageCreate, payload string) {
	if len(payload) == 0 {
		s.sendUsageMessage(d, m)
		return
	}

	if m.ReferencedMessage == nil {
		s.sendUsageMessage(d, m)
		return
	}

	channel := payload
	if len(channel) <= 3 || channel[0:2] != "<#" || channel[len(channel)-1:] != ">" {
		s.sendUsageMessage(d, m)
		return
	}

	channel = channel[2 : len(channel)-1]
	prevMessage := m.Message.ReferencedMessage

	member, err := d.GuildMember(m.GuildID, prevMessage.Author.ID)
	if err != nil {
		s.Logf("(%s) unable to find guild member %s", m.Guild.Name, prevMessage.Author.ID)
		return
	}

	hook, err := d.CreateOrFetchChannelWebhook(m.GuildID, channel)
	if err != nil {
		s.Logf("(%s) unable to fetch or create webhook: %s", m.Guild.Name, err)
		return
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
				s.Logf("(%s) error getting attachment: %s", m.Guild.Name, err)
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
		AvatarURL:  prevMessage.Author.AvatarURL(""),
		Components: prevMessage.Components,
		Embeds:     prevMessage.Embeds,
	})
	if err != nil {
		s.Logf("(%s) unable to execute webhook?!: %s", m.Guild.Name, err)
		return
	}

	err = d.ChannelMessageDelete(prevMessage.ChannelID, prevMessage.ID)
	if err != nil {
		s.Logf("(%s) unable to delete previous message?!: %s", m.Guild.Name, err)
		// do not return, we want to delete the message
	}

	err = d.ChannelMessageDelete(m.ChannelID, m.Message.ID)
	if err != nil {
		s.Logf("(%s) unable to delete trigger message?!: %s", m.Guild.Name, err)
		return
	}
}
