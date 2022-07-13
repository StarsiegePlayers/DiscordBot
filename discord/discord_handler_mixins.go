package discord

import (
	"strings"

	"golang.org/x/exp/slices"
)

func (s *Service) mixinRoleCheck(fn MessageHandler) MessageHandler {
	return func(d *Session, m *MessageCreate, payload string) {
		roleFound := false

		for _, role := range m.Command.Roles {
			if roleID, ok := d.GuildConfig.NamedRoles[role]; ok && slices.Contains(m.Member.Roles, roleID) {
				roleFound = true
				break
			}
		}

		if !roleFound {
			s.Logf("(%s) [#%s] %s does not have any {%s} role - roles possessed: {%s}", m.Guild.Name, m.Channel.Name, m.Author.Username+"#"+m.Author.Discriminator, strings.Join(m.Command.Roles, ", "), strings.Join(m.Member.Roles, ", "))
			return
		}

		fn(d, m, payload)
	}
}

func (s *Service) mixinPermissionCheck(fn MessageHandler) MessageHandler {
	return func(d *Session, m *MessageCreate, payload string) {
		if m.Permissions&m.Command.Permission == 0 {
			s.Logf("(%s) [#%s] %s does not have the %d permission flag {permissions: %d}", m.Guild.Name, m.Channel.Name, m.Author.Username+"#"+m.Author.Discriminator, m.Command.Permission, m.Permissions)
			return
		}

		fn(d, m, payload)
	}
}
