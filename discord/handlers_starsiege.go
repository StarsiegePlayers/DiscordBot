package discord

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"text/template"

	"github.com/bwmarrin/discordgo"
)

func (s *Service) registerStarsiegeHandlers() {
	s.commands.Register(Command{
		Name:    "ls",
		Handler: s.lsHandler,
		Summary: "Lists running games in the current channel.",
		Usage:   "ls",
	})
	s.commands.Register(Command{
		Name:    "qc",
		Handler: s.qcHandler,
		Summary: "Sends a quick chat to the current channel.",
		Usage:   "qc [quick chat id]",
	})
	s.commands.Register(Command{
		Name:    "slap",
		Handler: s.slapHandler,
		Summary: "Slaps a user with a random Starsiege death message.",
		Usage:   "slap <@user>",
	})
}

func (s *Service) lsHandler(d *Session, m *MessageCreate, payload string) {

}

func (s *Service) qcHandler(d *Session, m *MessageCreate, payload string) {
	if qc, ok := s.quickChats[payload]; ok {
		qcFile, err := os.OpenFile("qc/"+qc.SoundFile, os.O_RDONLY, 0)
		if err != nil {
			s.Logf("(%s) error opening qc file %s: %s", m.Guild.Name, qc.SoundFile, err)
			return
		}

		defer qcFile.Close()

		qcFileData, err := io.ReadAll(qcFile)
		if err != nil {
			s.Logf("(%s) error reading qc file %s: %s", m.Guild.Name, qc.SoundFile, err)
			return
		}

		_, err = qcFile.Seek(0, 0)
		if err != nil {
			s.Logf("(%s) error seeking qc file %s: %s", m.Guild.Name, qc.SoundFile, err)
			return
		}

		msg := discordgo.MessageSend{
			Content: qc.Text,
			Files: []*discordgo.File{
				{
					ContentType: http.DetectContentType(qcFileData),
					Name:        qc.SoundFile,
					Reader:      qcFile,
				},
			},
		}

		_, err = d.ChannelMessageSendComplex(m.ChannelID, &msg)
		if err != nil {
			s.Logf("(%s) error sending qc message: %s", m.Guild.Name, err)
			return
		}
	}
}

type SlapTemplateData struct {
	User   string
	Target string
}

func (s *Service) slapHandler(d *Session, m *MessageCreate, payload string) {
	if len(m.Mentions) <= 0 {
		return
	}

	nBig, err := rand.Int(rand.Reader, big.NewInt(1024))
	if err != nil {
		s.Logf("(%s) error generating random number:", m.Guild.Name, err)
		return
	}

	mod := new(big.Int)
	section := mod.Mod(nBig, big.NewInt(3))
	item := new(big.Int)
	output := ""

	switch int(section.Int64()) {
	case 0:
		item.Mod(nBig, big.NewInt(int64(len(s.slap.Active))))
		output = s.slap.Active[item.Int64()]
	case 1:
		item.Mod(nBig, big.NewInt(int64(len(s.slap.Passive))))
		output = s.slap.Passive[item.Int64()]
	default:
		item.Mod(nBig, big.NewInt(int64(len(s.slap.Generic))))
		output = s.slap.Generic[item.Int64()]
	}

	t, err := template.New(fmt.Sprintf("qc_%d_%d", int(section.Int64()), int(item.Int64()))).Parse(output)
	if err != nil {
		s.Logf("(%s) error parsing template:", m.Guild.Name, err)
		return
	}

	outBuff := bytes.NewBufferString("")

	err = t.Execute(outBuff, &SlapTemplateData{
		User:   m.Author.Mention(),
		Target: m.Mentions[0].Mention(),
	})
	if err != nil {
		s.Logln("(%s) error executing template:", m.Guild.Name, err)
		return
	}

	_, err = s.session.ChannelMessageSend(m.ChannelID, outBuff.String())
	if err != nil {
		s.Logln("(%s) error sending message:", m.Guild.Name, err)
		return
	}
}
