package discord

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

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

	s.registerAdminHandlers()
	s.registerStarsiegeHandlers()
}

func (s *Service) commandsHandler(d *Session, m *MessageCreate, payload string) {
	options := strings.Split(payload, " ")
	if len(options) == 1 && options[0] == "" {
		var commands []string
		for k, v := range s.commands {
			if v.Permission != 0 && m.Member.Permissions&v.Permission != 0 {
				commands = append(commands, k)
				continue
			}

			if len(v.Roles) > 0 && len(m.Member.Roles) > 0 {
				for _, role := range v.Roles {
					if slices.Contains(m.Member.Roles, role) {
						commands = append(commands, k)
						break
					}
				}
			}
		}

		_, err := d.ChannelMessageMentionSend(m.ChannelID, m.Author, fmt.Sprintf("Available commands: %s", strings.Join(commands, ", ")))
		if err != nil {
			s.Logln("error sending commands message:", err)
			return
		}
		return
	}

	if cmd, exists := s.commands[options[0]]; exists {
		_, err := d.ChannelMessageMentionSend(m.ChannelID, m.Author, fmt.Sprintf("[%s] %s", cmd.Name, cmd.Summary))
		if err != nil {
			s.Logln("error sending specific command summary message:", err)
			return
		}
		_, err = d.ChannelMessageMentionSend(m.ChannelID, m.Author, fmt.Sprintf("[%s] %s", cmd.Name, s.formatUsageMessage(d.GuildConfig.CommandPrefix, cmd.Usage)))
		if err != nil {
			s.Logln("error sending specific command usage message:", err)
			return
		}
		return
	}

	_, err := d.ChannelMessageMentionSend(m.ChannelID, m.Author, fmt.Sprintf("Unknown command: %s", options[0]))
	if err != nil {
		s.Logln("error sending unknown command message:", err)
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
