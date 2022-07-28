package discord

import (
	"fmt"
	"golang.org/x/exp/slices"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"sort"
	"strings"

	"github.com/Neo-Desktop/emojipasta"
	"github.com/bwmarrin/discordgo"
	owo "github.com/deadshot465/owoify-go/v2"
)

const DefaultColor = 0xff88ff

func (s *Service) registerHandlers() {
	s.commands.Register(Command{
		Name:    "commands",
		Handler: s.commandsHandler,
		Summary: "Describes commands available to the current user.",
		Usage:   "commands [command]",
	})
	s.commands.Register(Command{
		Name:    "help",
		Handler: s.commandsHandler,
		Summary: "Describes commands available to the current user.",
		Usage:   "help [command]",
	})
	s.commands.Register(Command{
		Name:    "ping",
		Handler: s.pingHandler,
		Summary: "Sends a DM to the user with a message validating the bot is online.",
		Usage:   "ping",
	})
	s.commands.Register(Command{
		Name:    "owo",
		Handler: s.textTransformHandler,
		Summary: "Turns ywour text intwo swomwething owo-ifwied.",
		Usage:   "owo [-me] <text>",
	})
	s.commands.Register(Command{
		Name:    "uwu",
		Handler: s.textTransformHandler,
		Summary: "Tuwns ywouw text intwo swomwefwing owo-ifwied(o･ω･o).",
		Usage:   "uwu [-me] <text>",
	})
	s.commands.Register(Command{
		Name:    "uvu",
		Handler: s.textTransformHandler,
		Summary: "Tuwns ywowouw text indwowo swowomwefwing owowowo-ifwieduwu.",
		Usage:   "uvu [-me] <text>",
	})
	s.commands.Register(Command{
		Name:    "emojipasta",
		Handler: s.textTransformHandler,
		Summary: "Turns😖 your👅👺⛓ text📱😤 into🔍👉 something💝🚫👤 emojipasta-ified.",
		Usage:   "emojipasta [-me] <text>",
	})
	s.commands.Register(Command{
		Name:    "sarcasm",
		Handler: s.textTransformHandler,
		Summary: "TuRnS YoUr tExT InTo sOmEtHiNg sArCaStIc-iFiEd",
		Usage:   "sarcasm [-me] <text>",
	})

	s.registerAdminHandlers()
	s.registerStarsiegeHandlers()
}

func (s *Service) commandsHandler(d *Session, m *MessageCreate, payload string) {
	var (
		err                  error
		fields               []*discordgo.MessageEmbedField
		titleCaseCommand     = []byte(m.Command.Name)
		titleCaseTransformer = cases.Title(language.English)
		options              = strings.Split(payload, " ")
		nick                 = m.Member.Nick
	)

	if nick == "" {
		nick = m.Author.Username
	}

	_, _, err = titleCaseTransformer.Transform(titleCaseCommand, []byte(m.Command.Name), true)
	if err != nil {
		s.Logln("error transforming command name to title case:", err)
	}

	if len(options) == 1 && options[0] == "" {
		var commands []string
		for k, v := range s.commands {
			if v.Permission == 0 && len(v.Roles) == 0 {
				commands = append(commands, k)
				continue
			}

			if ok, err := s.memberHasPermission(d, m.GuildID, m.Author.ID, v.Permission); ok && err != nil {
				commands = append(commands, k)
				continue
			} else if err != nil {
				s.Logln("error checking permissions:", err)
				continue
			}

			if len(v.Roles) > 0 && len(m.Member.Roles) > 0 {
				for _, roleName := range v.Roles {
					if roleID, ok := d.GuildConfig.NamedRoles[roleName]; ok {
						if slices.Contains(m.Member.Roles, roleID) {
							commands = append(commands, k)
							break
						}
					} else {
						s.Logf("(%s)[%s] error: role %s not defined", m.Guild.Name, m.GuildID, roleName)
					}
				}
			}
		}

		sort.Strings(commands)
		for _, v := range commands {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:  fmt.Sprintf("%s%s", d.GuildConfig.CommandPrefix, s.commands[v].Usage),
				Value: s.commands[v].Summary,
			})
		}

		embed := &discordgo.MessageEmbed{
			Title:  fmt.Sprintf("%s | %s", titleCaseCommand, nick),
			Color:  DefaultColor,
			Fields: fields,
		}

		_, err := d.ChannelMessageSendEmbed(m.ChannelID, embed)
		if err != nil {
			s.Logln("error sending commands message:", err)
			return
		}
		return
	}

	if cmd, exists := s.commands[options[0]]; exists {
		fields = []*discordgo.MessageEmbedField{
			{
				Name:  fmt.Sprintf("%s%s", d.GuildConfig.CommandPrefix, cmd.Usage),
				Value: cmd.Summary,
			},
		}

		embed := &discordgo.MessageEmbed{
			Title:  fmt.Sprintf("%s | %s", titleCaseCommand, nick),
			Color:  DefaultColor,
			Fields: fields,
		}

		_, err := s.ChannelMessageSend(m.ChannelID, embed)
		if err != nil {
			s.Logln("error sending commands message:", err)
			return
		}

		return
	}

	fields = []*discordgo.MessageEmbedField{
		{
			Name:  fmt.Sprintf("Unknown command: %s", options[0]),
			Value: fmt.Sprintf("Try `%s%s` for a list of commands.", d.GuildConfig.CommandPrefix, m.Command.Name),
		},
	}

	embed := &discordgo.MessageEmbed{
		Title:  fmt.Sprintf("%s | %s", titleCaseCommand, nick),
		Color:  DefaultColor,
		Fields: fields,
	}

	_, err = s.ChannelMessageSend(m.ChannelID, embed)
	if err != nil {
		s.Logln("error sending commands message:", err)
		return
	}
}

func (s *Service) pingHandler(d *Session, m *MessageCreate, payload string) {
	dm, err := d.UserChannelCreate(m.Author.ID)
	if err != nil {
		s.Logln("error creating channel:", err)
		return
	}

	_, err = d.ChannelMessageSend(dm.ID, "Pong!")
	if err != nil {
		s.Logln("error sending DM message:", err)
	}
}

func (s *Service) textTransformOwo(payload string, owoness owo.Owoness) string {
	return owo.Owoify(payload, owoness)
}

func (s *Service) textTransformSarcasm(payload string) string {
	var (
		upper    = strings.ToUpper(payload)
		lower    = strings.ToLower(payload)
		bPayload = []byte(payload)
	)

	for k := range bPayload {
		if k%2 == 0 {
			bPayload[k] = upper[k]
		} else {
			bPayload[k] = lower[k]
		}
	}

	return string(bPayload)
}

func (s *Service) textTransformEmojipasta(payload string) string {
	return emojipasta.Generate(strings.Split(payload, " "), 3, false)
}

func (s *Service) textTransformHandler(d *Session, m *MessageCreate, payload string) {
	var (
		owoness owo.Owoness
		flagMe  bool
	)

	member, err := d.State.Member(m.GuildID, m.Author.ID)
	if err != nil {
		s.Logln("error getting member:", err)
		return
	}

	if m.ReferencedMessage == nil {
		p := strings.Split(payload, " ")

		for _, v := range p {
			if len(v) == 0 || v[0] != '-' {
				// stop parsing if we hit a non-flag
				break
			}

			switch v[1:] {
			case "me":
				flagMe = true
				p = p[1:]
			}
		}

		payload = strings.Join(p, " ")
		payload = strings.TrimSpace(payload)
	} else {
		payload = m.ReferencedMessage.Content
	}

	switch m.Command.Name {
	case "emojipasta":
		if len(payload) == 0 {
			title := "Usage✔"

			payload = fmt.Sprintf("error: please specify some text to %s-ify", m.Command.Name)
			payload = s.textTransformEmojipasta(payload)

			_, err = s.ChannelMessageSend(m.ChannelID, s.BuildDefaultEmbed(title, payload))
			if err != nil {
				s.Logln("error sending usage message:", err)
			}

			return
		}

		payload = s.textTransformEmojipasta(payload)

	case "sarcasm":
		if len(payload) == 0 {
			title := "Usage"
			title = s.textTransformSarcasm(title)

			payload = fmt.Sprintf("error: please specify some text to %s-ify", m.Command.Name)
			payload = s.textTransformSarcasm(payload)

			_, err = s.ChannelMessageSend(m.ChannelID, s.BuildDefaultEmbed(title, payload))
			if err != nil {
				s.Logln("error sending usage message:", err)
			}

			return
		}

		payload = s.textTransformSarcasm(payload)

	case "owo":
		fallthrough
	case "uwu":
		fallthrough
	case "uvu":
		{
			switch m.Command.Name {
			case "owo":
				owoness = owo.Owo
			case "uwu":
				owoness = owo.Uwu
			case "uvu":
				owoness = owo.Uvu
			}

			if len(payload) == 0 {
				title := "Uwusawage"
				title = s.textTransformOwo(title, owoness)

				payload = fmt.Sprintf("error: please specify some text to %s-ify", m.Command.Name)
				payload = s.textTransformOwo(payload, owoness)

				_, err = s.ChannelMessageSend(m.ChannelID, s.BuildDefaultEmbed(title, payload))
				if err != nil {
					s.Logln("error sending usage message:", err)
				}

				return
			}

			payload = s.textTransformOwo(payload, owoness)
		}
	}

	if !flagMe {
		if len(payload) >= 1999 {
			var payloads []string
			for len(payload) > 1999 {
				payloads = append(payloads, payload[:1999])
				payload = payload[1999:]
			}
			payloads = append(payloads, payload)

			for _, p := range payloads {
				_, err = d.Session.ChannelMessageSend(m.ChannelID, p)
				if err != nil {
					s.Logln("error sending message:", err)
				}
			}
			return
		}

		_, err = d.Session.ChannelMessageSend(m.ChannelID, payload)
		if err != nil {
			s.Logln("error sending message:", err)
		}

		return
	}

	webhook, err := d.CreateOrFetchChannelWebhook(m.GuildID, m.ChannelID)
	if err != nil {
		s.Logln("error creating webhook:", err)
		return
	}

	err = d.ChannelMessageDelete(m.ChannelID, m.ID)
	if err != nil {
		s.Logln("error deleting message:", err)
		return
	}

	nick := member.Nick
	if nick == "" {
		nick = member.User.Username
	}

	if len(payload) >= 1999 {
		var payloads []string
		for len(payload) > 1999 {
			payloads = append(payloads, payload[:1999])
			payload = payload[1999:]
		}
		payloads = append(payloads, payload)

		for _, p := range payloads {
			_, err = d.WebhookExecute(webhook.ID, webhook.Token, false, &discordgo.WebhookParams{
				Content:   p,
				Username:  nick,
				AvatarURL: member.AvatarURL(""),
			})
			if err != nil {
				s.Logln("error sending message:", err)
			}
		}
		return
	}

	_, err = d.WebhookExecute(webhook.ID, webhook.Token, false, &discordgo.WebhookParams{
		Content:   payload,
		Username:  nick,
		AvatarURL: member.AvatarURL(""),
	})
	if err != nil {
		s.Logln("error sending message:", err)
	}
}
