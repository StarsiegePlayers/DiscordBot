package discord

import (
	"strings"

	"golang.org/x/exp/maps"
)

//go:generate stringer -type=Permission
type Permission int64

var (
	_Permission_keys   = maps.Keys(_Permission_map)
	_Permission_values = maps.Values(_Permission_map)
)

const (
	CREATE_INSTANT_INVITE Permission = 1 << iota
	KICK_MEMBERS
	BAN_MEMBERS
	ADMINISTRATOR
	MANAGE_CHANNELS
	MANAGE_GUILD
	ADD_REACTIONS
	VIEW_AUDIT_LOG
	PRIORITY_SPEAKER
	STREAM
	VIEW_CHANNEL
	SEND_MESSAGES
	SEND_TTS_MESSAGES
	MANAGE_MESSAGES
	EMBED_LINKS
	ATTACH_FILES
	READ_MESSAGE_HISTORY
	MENTION_EVERYONE
	USE_EXTERNAL_EMOJIS
	VIEW_GUILD_INSIGHTS
	CONNECT
	SPEAK
	MUTE_MEMBERS
	DEAFEN_MEMBERS
	MOVE_MEMBERS
	USE_VAD
	CHANGE_NICKNAME
	MANAGE_NICKNAMES
	MANAGE_ROLES
	MANAGE_WEBHOOKS
	MANAGE_EMOJIS_AND_STICKERS
	USE_APPLICATION_COMMANDS
	REQUEST_TO_SPEAK
	MANAGE_EVENTS
	MANAGE_THREADS
	CREATE_PUBLIC_THREADS
	CREATE_PRIVATE_THREADS
	USE_EXTERNAL_STICKERS
	SEND_MESSAGES_IN_THREADS
	USE_EMBEDDED_ACTIVITIES
	MODERATE_MEMBERS
)

func (p Permission) Has(permission Permission) bool {
	return p&permission == permission
}

func (p Permission) MaskString() string {
	var (
		list      []string
		remainder = p
	)

	for _, perm := range _Permission_keys {
		if p.Has(perm) {
			list = append(list, perm.String())
			remainder -= perm
		}
	}

	if remainder > 0 {
		list = append(list, remainder.String())
	}

	return strings.Join(list, ", ")
}
