package discord

import (
	"fmt"
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
	s.registerHandlers()

	s.RPCSubscribe(rpc.DiscordMessageSendTopic, s.discordMessageSendRPCHandler)
	s.RPCSubscribe(rpc.APIRequestResponse, s.apiRequestResponseRPCHandler)
	s.RPCSubscribe(rpc.NewConfigLoadedTopic, s.configMessageRPCHandler)
}

func (s *Service) Start() (err error) {
	// wait for config
	s.wg.Wait()

	s.session, err = discordgo.New("Bot " + s.config.AuthToken)
	if err != nil {
		return err
	}

	s.session.AddHandler(s.initMessage)
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

func (s *Service) sendUsageMessage(d *Session, m *MessageCreate) {
	_, err := d.ChannelMessageMentionSend(m.ChannelID, m.Author, s.formatUsageMessage(d.GuildConfig.CommandPrefix, m.Command.Usage))
	if err != nil {
		s.Logf("(%s) error while sending usage message: %s", m.Guild.Name, err)
	}
}

func (s *Service) memberHasPermission(d *Session, guildID string, userID string, permission int64) (bool, error) {
	member, err := d.State.Member(guildID, userID)
	if err != nil {
		if member, err = d.GuildMember(guildID, userID); err != nil {
			return false, err
		}
	}

	for _, roleID := range member.Roles {
		role, err := d.State.Role(guildID, roleID)
		if err != nil {
			return false, err
		}
		if role.Permissions&permission != 0 {
			return true, nil
		}
	}

	return false, nil
}

func (s *Service) formatUsageMessage(commandPrefix string, usage string) string {
	return fmt.Sprintf("Usage: `%s%s`", commandPrefix, usage)
}

func (s *Service) messageDispatcher(d *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == d.State.User.ID {
		return
	}

	cfg, hasConfig := s.config.Guilds[m.GuildID]

	session := &Session{
		Session:     d,
		GuildConfig: cfg,
	}

	// did we receive a command?
	if (hasConfig && strings.HasPrefix(m.Content, s.config.Guilds[m.GuildID].CommandPrefix)) ||
		(!hasConfig && strings.HasPrefix(m.Content, DefaultCommandPrefix+"init")) ||
		(s.isMentioned(m.Mentions, d.State.User)) {

		content := m.Content

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
		if cmd, ok := s.commands[command]; ok {
			guild, _ := d.State.Guild(m.GuildID)
			channel, _ := d.State.Channel(m.ChannelID)
			member, _ := d.State.Member(m.GuildID, m.Author.ID)
			perms, _ := d.State.UserChannelPermissions(m.Author.ID, m.ChannelID)

			msg := &MessageCreate{
				Guild:         guild,
				Channel:       channel,
				Member:        member,
				Command:       cmd,
				Permissions:   perms,
				MessageCreate: m,
			}

			s.Logf("(%s) [#%s] <%s>: %s", guild.Name, channel.Name, m.Author.Username+"#"+m.Author.Discriminator, m.Content)

			go cmd.Handler(session, msg, content)
		}
	}
}

func (s *Service) initMessage(*discordgo.Session, *discordgo.Ready) {
	// wait for config
	s.wg.Wait()

	for name, id := range s.config.DebugUsers {
		s.Logf("sending init message to %s(%s)", name, id)
		dm, err := s.session.UserChannelCreate(id)
		if err != nil {
			s.Logln("error creating init user channel:", err)
			return
		}

		_, err = s.session.ChannelMessageSend(dm.ID, fmt.Sprintf("[%s] - bot is online", time.Now().Format(time.ANSIC)))
		if err != nil {
			s.Logln("error sending init message:", err)
			return
		}
	}

	for _, v := range s.session.State.Guilds {
		if cfg, ok := s.config.Guilds[v.ID]; ok {
			// we have a guild config, perform muzzle maintenance
			s.muzzleMaintenance(v.ID, cfg)
		}

		s.Logf("Guild Registered: %s(%s)", v.Name, v.ID)
		channels, err := s.session.GuildChannels(v.ID)
		if err != nil {
			s.Logln("error fetching channels", err)
		}

		out := make([]string, 0, len(channels))
		for _, c := range channels {
			if c.Type == discordgo.ChannelTypeGuildText {
				out = append(out, fmt.Sprintf("%s(%s)", c.Name, c.ID))
			}
		}

		s.Logf("Channels for %s: %s", v.Name, strings.Join(out, ", "))
	}

}
