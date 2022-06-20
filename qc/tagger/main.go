package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"

	id3v2 "github.com/bogem/id3v2/v2"
)

type quickchat struct {
	Text      string `json:"text"`
	SoundFile string `json:"soundFile"`
}

func main() {
	quickChats := make(map[string]quickchat)

	log.Println("loading quickchats")

	f, err := os.OpenFile("..\\..\\json\\qc.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Println(err)
		return
	}

	defer f.Close()

	qc, err := io.ReadAll(f)
	if err != nil {
		log.Println(err)
		return
	}

	err = json.Unmarshal(qc, &quickChats)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("loaded %d quickchats", len(quickChats))

	for k, v := range quickChats {
		tag, err := id3v2.Open("..\\"+v.SoundFile, id3v2.Options{Parse: true})
		if err != nil {
			log.Printf("error opening quickchat %s - %s (%s)", k, v.SoundFile, err)
			continue
		}
		title := strings.Split(v.SoundFile, ".")
		tag.SetTitle(title[0])
		tag.AddTextFrame(tag.CommonID("Encoded by"), tag.DefaultEncoding(), "")
		tag.AddTextFrame(tag.CommonID("Software/Hardware and settings used for encoding"), tag.DefaultEncoding(), "Sydney <3")
		tag.AddTextFrame(tag.CommonID("Track number/Position in set"), tag.DefaultEncoding(), k)
		comment := id3v2.CommentFrame{
			Encoding:    tag.DefaultEncoding(),
			Language:    "eng",
			Description: "QuickChat Text",
			Text:        v.Text,
		}
		tag.AddCommentFrame(comment)
		tag.Save()
		tag.Close()
	}

}
