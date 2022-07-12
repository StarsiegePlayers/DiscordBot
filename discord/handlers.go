package discord

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"image/png"
	"io"
	"math/big"
	"net/http"
	"os"
	"strings"
	"text/template"

	"golang.org/x/exp/maps"

	"github.com/bwmarrin/discordgo"
)

type SlapTemplateData struct {
	User   string
	Target string
}

func (s *Service) initHandler(d *Session, m *discordgo.MessageCreate, payload string) {

}

func (s *Service) commandsHandler(d *Session, m *discordgo.MessageCreate, payload string) {
	if len(payload) <= 1 {
		message := m.Author.Mention() + " Commands: " + strings.Join(maps.Keys(s.commands), ", ")
		_, err := s.session.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			s.Logln(err)
			return
		}
	}
}

func (s *Service) pingHandler(d *Session, m *discordgo.MessageCreate, payload string) {
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

func (s *Service) lsHandler(d *Session, m *discordgo.MessageCreate, payload string) {

}

func (s *Service) qcHandler(d *Session, m *discordgo.MessageCreate, payload string) {
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

func (s *Service) slapHandler(d *Session, m *discordgo.MessageCreate, payload string) {
	if len(m.Mentions) <= 0 {
		return
	}

	nBig, err := rand.Int(rand.Reader, big.NewInt(1024))
	if err != nil {
		s.Logln("error generating random number:", err)
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
		s.Logln("error parsing template:", err)
		return
	}

	outBuff := bytes.NewBufferString("")

	err = t.Execute(outBuff, &SlapTemplateData{
		User:   m.Author.Mention(),
		Target: m.Mentions[0].Mention(),
	})
	if err != nil {
		s.Logln("error executing template:", err)
		return
	}

	_, err = s.session.ChannelMessageSend(m.ChannelID, outBuff.String())
	if err != nil {
		s.Logln("error sending message:", err)
		return
	}
}

func (s *Service) moveHandler(d *Session, m *discordgo.MessageCreate, payload string) {
	if len(payload) == 0 {
		_, err := s.session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@!%s>, Usage: in a reply to the message you wish to move, `!move <channel>`", m.Author.ID))
		if err != nil {
			s.Logln("error sending syntax message:", err)
			return
		}
		return
	}

	if m.ReferencedMessage == nil {
		_, err := s.session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@!%s>, you must reference or reply to a message to move it.", m.Author.ID))
		if err != nil {
			s.Logln("error sending error message:", err)
		}
		return
	}

	channel := payload
	if len(channel) <= 3 || channel[0:2] != "<#" || channel[len(channel)-1:] != ">" {
		_, err := s.session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@!%s>, you must specify a channel to move the message to.", m.Author.ID))
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

	avatarImg, err := d.UserAvatarDecode(prevMessage.Author)
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

	avatarBase64 := base64.StdEncoding.EncodeToString(avatarPng.Bytes())
	avatarBase64 = "data:image/png;base64," + avatarBase64

	nick := member.Nick
	if nick == "" {
		nick = member.User.Username
	}

	hook, err := d.WebhookCreate(channel, nick, avatarBase64)
	if err != nil {
		s.Logf("Unable to create webhook in %s for %s (%s)", channel, prevMessage.Author.ID, err)
		return
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

	err = d.WebhookDelete(hook.ID)
	if err != nil {
		s.Logln("unable to delete webhook?!", err)
		// do not return, we want to delete the message
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
