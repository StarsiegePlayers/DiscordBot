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

type MessageHandler func(*Session, *discordgo.MessageCreate, string)
type Commands map[string]MessageHandler

func (c *Commands) Register(name string, fn MessageHandler) {
	(*c)[name] = fn
}

type Session struct {
	GuildConfig config.GuildConfig
	*discordgo.Session
}

type Service struct {
	module.Base

	wg       sync.WaitGroup
	commands Commands
	session  *discordgo.Session

	config *config.DiscordConfig

	quickChats map[string]quickchat
	slap       slap
}

func (s *Service) Init() {
	// queue waiting on config
	s.wg.Add(1)

	s.loadDataFiles()

	s.commands = make(Commands)
	s.commands.Register("init", s.messageLogger(s.initHandler))
	s.commands.Register("commands", s.messageLogger(s.commandsHandler))
	s.commands.Register("help", s.messageLogger(s.commandsHandler))
	s.commands.Register("ping", s.messageLogger(s.pingHandler))
	s.commands.Register("ls", s.messageLogger(s.lsHandler))
	s.commands.Register("qc", s.messageLogger(s.qcHandler))
	s.commands.Register("slap", s.messageLogger(s.slapHandler))
	s.commands.Register("move", s.messageLogger(s.roleCheck("Staff", s.moveHandler)))

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

func (s *Service) messageLogger(fn MessageHandler) MessageHandler {
	return func(d *Session, m *discordgo.MessageCreate, payload string) {
		guild, _ := d.State.Guild(m.GuildID)
		channel, _ := d.State.Channel(m.ChannelID)

		s.Logf("(%s) [#%s] <%s>: %s", guild.Name, channel.Name, m.Author.Username+"#"+m.Author.Discriminator, m.Content)

		fn(d, m, payload)
	}
}

func (s *Service) roleCheck(namedRole string, fn MessageHandler) MessageHandler {
	return func(d *Session, m *discordgo.MessageCreate, payload string) {
		guild, _ := d.State.Guild(m.GuildID)
		channel, _ := d.State.Channel(m.ChannelID)
		member, _ := d.State.Member(m.GuildID, m.Author.ID)

		if roleID, ok := d.GuildConfig.NamedRoles[namedRole]; !ok || !slices.Contains(member.Roles, roleID) {
			s.Logf("(%s) [#%s] %s does not have the %s (%s) role - roles possessed: {%s}", guild.Name, channel.Name, m.Author.Username+"#"+m.Author.Discriminator, namedRole, roleID, strings.Join(member.Roles, ", "))
			return
		}

		fn(d, m, payload)
	}
}

func (s *Service) permissionCheck(permission int64, fn MessageHandler) MessageHandler {
	return func(d *Session, m *discordgo.MessageCreate, payload string) {
		guild, _ := d.State.Guild(m.GuildID)
		channel, _ := d.State.Channel(m.ChannelID)
		perms, _ := d.State.UserChannelPermissions(m.Author.ID, m.ChannelID)

		if perms&permission == 0 {
			s.Logf("(%s) [#%s] %s does not have the %d permission flag {permissions: %d}", guild.Name, channel.Name, m.Author.Username+"#"+m.Author.Discriminator, permission, perms)
			return
		}

		fn(d, m, payload)
	}
}

func (s *Service) messageDispatcher(d *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == d.State.User.ID {
		return
	}

	// should we enforce moderation?
	if s.PerformModeration(d, m) {
		// moderation was performed
		return
	}

	// did we receive a command?
	if cfg, ok := s.config.Guilds[m.GuildID]; (ok && strings.HasPrefix(m.Content, s.config.Guilds[m.GuildID].CommandPrefix)) ||
		(s.isMentioned(m.Mentions, d.State.User)) ||
		(!ok && strings.HasPrefix(m.Content, DefaultCommandPrefix+"init")) {

		content := m.Content
		session := &Session{
			Session:     d,
			GuildConfig: cfg,
		}

		if s.isMentioned(m.Mentions, d.State.User) {
			content = strings.ReplaceAll(content, d.State.User.Mention(), "")
			content = strings.TrimSpace(content)
		}

		content = strings.TrimSpace(content)

		command := strings.SplitN(content, " ", 2)[0]
		command = strings.ToLower(command)

		content = strings.TrimPrefix(content, command)
		content = strings.TrimSpace(content)

		if strings.HasPrefix(command, s.config.Guilds[m.GuildID].CommandPrefix) {
			command = strings.TrimPrefix(command, s.config.Guilds[m.GuildID].CommandPrefix)
			command = strings.TrimSpace(command)
		}

		// dispatch message to correct function, if registered
		if fn, ok := s.commands[command]; ok {
			go fn(session, m, content)
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
