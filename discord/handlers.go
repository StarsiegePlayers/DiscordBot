package discord

import (
	"fmt"
	"sort"
	"strings"

	"golang.org/x/exp/slices"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

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
		Handler: s.owoHandler,
		Summary: "Turns ywour text intwo swomwething owo-ifwied.",
		Usage:   "owo <text>",
	})
	s.commands.Register(Command{
		Name:    "uwu",
		Handler: s.owoHandler,
		Summary: "Tuwns ywouw text intwo swomwefwing owo-ifwied(o･ω･o).",
		Usage:   "uwu <text>",
	})
	s.commands.Register(Command{
		Name:    "uvu",
		Handler: s.owoHandler,
		Summary: "Tuwns ywowouw text indwowo swowomwefwing owowowo-ifwieduwu.",
		Usage:   "uvu <text>",
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

		_, err := d.ChannelMessageSendEmbed(m.ChannelID, embed)
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

	_, err = d.ChannelMessageSendEmbed(m.ChannelID, embed)
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

func (s *Service) owoHandler(d *Session, m *MessageCreate, payload string) {
	var (
		owoness owo.Owoness
		output  string
	)

	switch m.Command.Name {
	case "owo":
		owoness = owo.Owo
	case "uwu":
		owoness = owo.Uwu
	case "uvu":
		owoness = owo.Uvu
	}

	if len(payload) == 0 {
		payload = fmt.Sprintf("error: please specify some text to %s", m.Command.Name)
	}

	output = owo.Owoify(payload, owoness)

	_, err := d.ChannelMessageMentionSend(m.ChannelID, m.Author, output)
	if err != nil {
		s.Logln("error sending message:", err)
	}
}
