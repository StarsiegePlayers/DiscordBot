package discord

import (
	"golang.org/x/exp/slices"
	"time"

	"github.com/bwmarrin/discordgo"
)

func (s *Service) PerformModeration(d *discordgo.Session, m *discordgo.MessageCreate) bool {
	// does this server have a guild config?
	if _, ok := s.config.Guilds[m.GuildID]; !ok {
		return false
	}

	// was the message sent in an exempt channel?
	guildConfig := s.config.Guilds[m.GuildID]
	if len(guildConfig.TimeoutConfig.ExemptChannels) <= 0 || slices.Contains(guildConfig.TimeoutConfig.ExemptChannels, m.ChannelID) {
		return false
	}

	s.Logln("guild config:", guildConfig)

	// does the member currently have the timeout role?
	if !slices.Contains(m.Member.Roles, guildConfig.TimeoutConfig.TimeoutRoleID) {
		return false
	}

	// perform the timeout
	timePlusDuration := time.Now().Add(time.Duration(guildConfig.TimeoutConfig.TimeoutTTL) * time.Second)
	s.Logf("timing out %s(%s) until %s", m.Author.Username+m.Author.Discriminator, m.Author.ID, timePlusDuration.String())

	err := s.session.GuildMemberTimeout(m.GuildID, m.Author.ID, &timePlusDuration)
	if err != nil {
		s.Logln(err)
	}

	return true
}
