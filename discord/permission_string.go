// Code generated by "stringer -type=Permission"; DO NOT EDIT.

package discord

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[CREATE_INSTANT_INVITE-1]
	_ = x[KICK_MEMBERS-2]
	_ = x[BAN_MEMBERS-4]
	_ = x[ADMINISTRATOR-8]
	_ = x[MANAGE_CHANNELS-16]
	_ = x[MANAGE_GUILD-32]
	_ = x[ADD_REACTIONS-64]
	_ = x[VIEW_AUDIT_LOG-128]
	_ = x[PRIORITY_SPEAKER-256]
	_ = x[STREAM-512]
	_ = x[VIEW_CHANNEL-1024]
	_ = x[SEND_MESSAGES-2048]
	_ = x[SEND_TTS_MESSAGES-4096]
	_ = x[MANAGE_MESSAGES-8192]
	_ = x[EMBED_LINKS-16384]
	_ = x[ATTACH_FILES-32768]
	_ = x[READ_MESSAGE_HISTORY-65536]
	_ = x[MENTION_EVERYONE-131072]
	_ = x[USE_EXTERNAL_EMOJIS-262144]
	_ = x[VIEW_GUILD_INSIGHTS-524288]
	_ = x[CONNECT-1048576]
	_ = x[SPEAK-2097152]
	_ = x[MUTE_MEMBERS-4194304]
	_ = x[DEAFEN_MEMBERS-8388608]
	_ = x[MOVE_MEMBERS-16777216]
	_ = x[USE_VAD-33554432]
	_ = x[CHANGE_NICKNAME-67108864]
	_ = x[MANAGE_NICKNAMES-134217728]
	_ = x[MANAGE_ROLES-268435456]
	_ = x[MANAGE_WEBHOOKS-536870912]
	_ = x[MANAGE_EMOJIS_AND_STICKERS-1073741824]
	_ = x[USE_APPLICATION_COMMANDS-2147483648]
	_ = x[REQUEST_TO_SPEAK-4294967296]
	_ = x[MANAGE_EVENTS-8589934592]
	_ = x[MANAGE_THREADS-17179869184]
	_ = x[CREATE_PUBLIC_THREADS-34359738368]
	_ = x[CREATE_PRIVATE_THREADS-68719476736]
	_ = x[USE_EXTERNAL_STICKERS-137438953472]
	_ = x[SEND_MESSAGES_IN_THREADS-274877906944]
	_ = x[USE_EMBEDDED_ACTIVITIES-549755813888]
	_ = x[MODERATE_MEMBERS-1099511627776]
}

const _Permission_name = "CREATE_INSTANT_INVITEKICK_MEMBERSBAN_MEMBERSADMINISTRATORMANAGE_CHANNELSMANAGE_GUILDADD_REACTIONSVIEW_AUDIT_LOGPRIORITY_SPEAKERSTREAMVIEW_CHANNELSEND_MESSAGESSEND_TTS_MESSAGESMANAGE_MESSAGESEMBED_LINKSATTACH_FILESREAD_MESSAGE_HISTORYMENTION_EVERYONEUSE_EXTERNAL_EMOJISVIEW_GUILD_INSIGHTSCONNECTSPEAKMUTE_MEMBERSDEAFEN_MEMBERSMOVE_MEMBERSUSE_VADCHANGE_NICKNAMEMANAGE_NICKNAMESMANAGE_ROLESMANAGE_WEBHOOKSMANAGE_EMOJIS_AND_STICKERSUSE_APPLICATION_COMMANDSREQUEST_TO_SPEAKMANAGE_EVENTSMANAGE_THREADSCREATE_PUBLIC_THREADSCREATE_PRIVATE_THREADSUSE_EXTERNAL_STICKERSSEND_MESSAGES_IN_THREADSUSE_EMBEDDED_ACTIVITIESMODERATE_MEMBERS"

var _Permission_map = map[Permission]string{
	1:             _Permission_name[0:21],
	2:             _Permission_name[21:33],
	4:             _Permission_name[33:44],
	8:             _Permission_name[44:57],
	16:            _Permission_name[57:72],
	32:            _Permission_name[72:84],
	64:            _Permission_name[84:97],
	128:           _Permission_name[97:111],
	256:           _Permission_name[111:127],
	512:           _Permission_name[127:133],
	1024:          _Permission_name[133:145],
	2048:          _Permission_name[145:158],
	4096:          _Permission_name[158:175],
	8192:          _Permission_name[175:190],
	16384:         _Permission_name[190:201],
	32768:         _Permission_name[201:213],
	65536:         _Permission_name[213:233],
	131072:        _Permission_name[233:249],
	262144:        _Permission_name[249:268],
	524288:        _Permission_name[268:287],
	1048576:       _Permission_name[287:294],
	2097152:       _Permission_name[294:299],
	4194304:       _Permission_name[299:311],
	8388608:       _Permission_name[311:325],
	16777216:      _Permission_name[325:337],
	33554432:      _Permission_name[337:344],
	67108864:      _Permission_name[344:359],
	134217728:     _Permission_name[359:375],
	268435456:     _Permission_name[375:387],
	536870912:     _Permission_name[387:402],
	1073741824:    _Permission_name[402:428],
	2147483648:    _Permission_name[428:452],
	4294967296:    _Permission_name[452:468],
	8589934592:    _Permission_name[468:481],
	17179869184:   _Permission_name[481:495],
	34359738368:   _Permission_name[495:516],
	68719476736:   _Permission_name[516:538],
	137438953472:  _Permission_name[538:559],
	274877906944:  _Permission_name[559:583],
	549755813888:  _Permission_name[583:606],
	1099511627776: _Permission_name[606:622],
}

func (i Permission) String() string {
	if str, ok := _Permission_map[i]; ok {
		return str
	}
	return "Permission(" + strconv.FormatInt(int64(i), 10) + ")"
}