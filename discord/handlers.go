package discord

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"strings"
	"text/template"

	embed "github.com/Clinet/discordgo-embed"
	"github.com/bwmarrin/discordgo"
)

type SlapTemplateData struct {
	User   string
	Target string
}

func (s *Service) initHandler(d *discordgo.Session, m *discordgo.MessageCreate) {

}

func (s *Service) pingHandler(d *discordgo.Session, m *discordgo.MessageCreate) {
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

func (s *Service) lsHandler(d *discordgo.Session, m *discordgo.MessageCreate) {

}

func (s *Service) qcHandler(d *discordgo.Session, m *discordgo.MessageCreate) {
	payload := strings.SplitN(m.Content, " ", 2)[1]

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

func (s *Service) slapHandler(d *discordgo.Session, m *discordgo.MessageCreate) {

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

func (s *Service) eledoreHandler(d *discordgo.Session, m *discordgo.MessageCreate) {
	eledore := `ATTENTION//NOTE: 
The "Renewed Exemplar SecT" >REST<, has CHOSEN//DECIDED to expend its MEMBERS//UNITS by one.
The >REST< has Assigned '<@132328159726141442>' to join the Examplars in upholding the 'Core Directives' of this Server\Hub for all <<NEXT>> & <<HUMANS>>.
Reminder! Failure to uphold/follow the 'Core Directives' is ill advised!
ACKNOWLEDGE//SUBMIT!!`
	// _, _ = d.ChannelMessageSendEmbed("938596269880926279", embed.NewEmbed().SetDescription(eledore).SetColor(0x00a5e4).MessageEmbed)
	_, _ = d.ChannelMessageSendEmbed(m.ChannelID, embed.NewEmbed().SetDescription(eledore).SetColor(0x00a5e4).MessageEmbed)
}
