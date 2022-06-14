package discord

import (
	"fmt"
	"golang.org/x/exp/slices"
	"strings"
	"sync"
	"time"

	"github.com/StarsiegePlayers/DiscordBot/config"
	"github.com/StarsiegePlayers/DiscordBot/module"
	"github.com/StarsiegePlayers/DiscordBot/rpc"

	"github.com/bwmarrin/discordgo"
)

const (
	ServiceName          = "discord"
	DefaultCommandPrefix = "!"
)

type MessageHandler func(*discordgo.Session, *discordgo.MessageCreate)

type Service struct {
	module.Base

	wg       sync.WaitGroup
	commands map[string]MessageHandler
	session  *discordgo.Session

	config *config.DiscordConfig

	quickChats map[string]quickchat
	slap       slap
}

func (s *Service) Init() {
	// queue waiting on config
	s.wg.Add(1)

	s.loadFiles()

	s.commands = make(map[string]MessageHandler)
	s.commands["init"] = s.initHandler
	s.commands["ping"] = s.pingHandler
	s.commands["ls"] = s.lsHandler
	s.commands["qc"] = s.qcHandler
	s.commands["slap"] = s.slapHandler
	s.commands["eledore"] = s.eledoreHandler

	s.PubSubSubscribe(rpc.DiscordMessageSendTopic, s.discordMessageSendPubSubHandler)
	s.PubSubSubscribe(rpc.APIRequestResponse, s.APIRequestResponsePubSubHandler)
	s.PubSubSubscribe(rpc.NewConfigLoadedTopic, s.configMessagePubSubHandler)
}

func (s *Service) Start() (err error) {
	// wait for config
	s.wg.Wait()

	s.session, err = discordgo.New("Bot " + s.config.AuthToken)
	if err != nil {
		return err
	}

	s.session.AddHandler(s.initMessageSender)
	s.session.AddHandler(s.messageDispatcher)

	s.session.Identify.Intents = discordgo.IntentsAll

	err = s.PubSubPublish(rpc.APIRequestLatest, []byte{})
	if err != nil {
		return
	}

	err = s.session.Open()
	if err != nil {
		s.Log.Printf("cannot open the session: %v", err)
		return
	}

	return
}

func (s *Service) Stop() (err error) {
	if s.session != nil {
		err = s.session.Close()
	}

	return
}

func (s *Service) isMentioned(input []*discordgo.User, compare *discordgo.User) bool {
	for _, v := range input {
		if v.ID == compare.ID {
			return true
		}
	}

	return false
}

func (s *Service) messageDispatcher(d *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == d.State.User.ID {
		return
	}

	if slices.Contains(m.Member.Roles, s.config.Guilds[m.GuildID].TimeoutConfig.TimeoutRoleID) {
		timePlusDuration := time.Now().Add(time.Duration(s.config.Guilds[m.GuildID].TimeoutConfig.TimeoutTTL) * time.Second)
		err := s.session.GuildMemberTimeout(m.GuildID, m.Author.ID, &timePlusDuration)

		if err != nil {
			s.Logln(err)
		}

		return
	}

	// did we receive a command?
	if _, ok := s.config.Guilds[m.GuildID]; (ok && strings.HasPrefix(m.Content, s.config.Guilds[m.GuildID].CommandPrefix)) ||
		(s.isMentioned(m.Mentions, d.State.User)) ||
		(!ok && strings.HasPrefix(m.Content, DefaultCommandPrefix+"init")) {
		command := strings.SplitN(m.Content, " ", 2)[0]
		command = strings.ToLower(command)[1:]

		// dispatch message to correct function, if registered
		if fn, ok := s.commands[command]; ok {
			go fn(d, m)
		}
	}
}

func (s *Service) initMessageSender(*discordgo.Session, *discordgo.Ready) {
	// wait for config
	s.wg.Wait()

	for name, id := range s.config.DebugUsers {
		s.Logf("sending init message to %s(%s)", name, id)
		dm, err := s.session.UserChannelCreate(id)
		if err != nil {
			return
		}

		_, err = s.session.ChannelMessageSend(dm.ID, fmt.Sprintf("[%s] - bot is online", time.Now().Format(time.ANSIC)))
		if err != nil {
			return
		}
	}

	for _, v := range s.session.State.Guilds {
		s.Logf("Guild Registered: %s(%s)", v.Name, v.ID)
	}

	return
}
