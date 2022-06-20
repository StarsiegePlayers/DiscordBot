package discord

import (
	"encoding/json"
	"io"
	"os"
)

type quickchat struct {
	Text      string `json:"text"`
	SoundFile string `json:"soundFile"`
}

type slap struct {
	Active  []string `json:"active"`
	Passive []string `json:"passive"`
	Generic []string `json:"generic"`
}

func (s *Service) loadDataFiles() {
	s.loadQuickChats()
	s.loadSlaps()
}

func (s *Service) loadQuickChats() {
	s.Logln("loading quickchats")

	f, err := os.OpenFile("json/qc.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		s.Logln(err)
		return
	}

	defer f.Close()

	qc, err := io.ReadAll(f)
	if err != nil {
		s.Logln(err)
		return
	}

	err = json.Unmarshal(qc, &s.quickChats)
	if err != nil {
		s.Logln(err)
		return
	}

	s.Logf("loaded %d quickchats", len(s.quickChats))
}

func (s *Service) loadSlaps() {
	s.Logln("loading slaps")

	f, err := os.OpenFile("json/slap.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		s.Logln(err)
		return
	}

	defer f.Close()

	slaps, err := io.ReadAll(f)
	if err != nil {
		s.Logln(err)
		return
	}

	err = json.Unmarshal(slaps, &s.slap)
	if err != nil {
		s.Logln(err)
		return
	}

	s.Logf("loaded %d slaps", len(s.slap.Active)+len(s.slap.Passive)+len(s.slap.Generic))
}
