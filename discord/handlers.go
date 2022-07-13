package discord

import (
	"strings"

	"golang.org/x/exp/maps"
)

func (s *Service) registerHandlers() {
	s.commands.Register(Command{
		Name:    "commands",
		Handler: s.commandsHandler,
	})
	s.commands.Register(Command{
		Name:    "help",
		Handler: s.commandsHandler,
	})
	s.commands.Register(Command{
		Name:    "ping",
		Handler: s.pingHandler,
	})

	s.registerAdminHandlers()
	s.registerStarsiegeHandlers()
}

func (s *Service) commandsHandler(d *Session, m *MessageCreate, payload string) {
	if len(payload) <= 1 {
		message := m.Author.Mention() + " Commands: " + strings.Join(maps.Keys(s.commands), ", ")
		_, err := s.session.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			s.Logln(err)
			return
		}
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
