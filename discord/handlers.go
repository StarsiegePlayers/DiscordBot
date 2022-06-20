package discord

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"golang.org/x/exp/maps"
	"io"
	"math/big"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/bwmarrin/discordgo"
)

type SlapTemplateData struct {
	User   string
	Target string
}

func (s *Service) initHandler(d *discordgo.Session, m *discordgo.MessageCreate, payload string) {

}

func (s *Service) commandsHandler(d *discordgo.Session, m *discordgo.MessageCreate, payload string) {
	if len(payload) <= 1 {
		message := m.Author.Mention() + " Commands: " + strings.Join(maps.Keys(s.commands), ", ")
		_, err := s.session.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			s.Logln(err)
			return
		}
	}
}

func (s *Service) pingHandler(d *discordgo.Session, m *discordgo.MessageCreate, payload string) {
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

func (s *Service) lsHandler(d *discordgo.Session, m *discordgo.MessageCreate, payload string) {

}

func (s *Service) qcHandler(d *discordgo.Session, m *discordgo.MessageCreate, payload string) {
	if qc, ok := s.quickChats[payload]; ok {
		qcFile, err := os.OpenFile("qc/"+qc.SoundFile, os.O_RDONLY, 0)
		if err != nil {
			s.Logln(err)
			return
		}

		defer qcFile.Close()

		qcFileData, err := io.ReadAll(qcFile)
		if err != nil {
			s.Logln(err)
			return
		}

		_, err = qcFile.Seek(0, 0)
		if err != nil {
			s.Logln(err)
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
			s.Logln(err)
			return
		}
	}

}

func (s *Service) slapHandler(d *discordgo.Session, m *discordgo.MessageCreate, payload string) {
	if len(m.Mentions) <= 0 {
		return
	}

	nBig, err := rand.Int(rand.Reader, big.NewInt(1024))
	if err != nil {
		s.Logln(err)
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
		s.Logln(err)
		return
	}

	outBuff := bytes.NewBufferString("")
	err = t.Execute(outBuff, &SlapTemplateData{
		User:   m.Author.Mention(),
		Target: m.Mentions[0].Mention(),
	})

	_, err = s.session.ChannelMessageSend(m.ChannelID, outBuff.String())
	if err != nil {
		s.Logln(err)
		return
	}
}
