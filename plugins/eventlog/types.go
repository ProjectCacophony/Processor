package eventlog

import (
	"strings"
)

// Actions
const (
	ActionTypeModDM        actionType = "cacophony_mod_dm"
	ActionTypeDiscordBan   actionType = "discord_ban"
	ActionTypeDiscordUnban actionType = "discord_unban"
	ActionTypeDiscordJoin  actionType = "discord_join"
	ActionTypeDiscordLeave actionType = "discord_leave"
)

func (t actionType) String() string {
	switch t {
	case ActionTypeModDM:
		return "Mod DM"
	}

	return titleify(string(t))
}

// Entities
const (
	EntityTypeUser    entityType = "discord_user"
	EntityTypeRole    entityType = "discord_role"
	EntityTypeGuild   entityType = "discord_guild"
	EntityTypeChannel entityType = "discord_channel"

	EntityTypeMessageCode entityType = "cacophony_message_code"
)

func (t entityType) String(value string) string {
	switch t {
	case EntityTypeUser:
		return "<@" + value + "> #" + value
	case EntityTypeRole:
		return "<@&" + value + "> #" + value
	case EntityTypeGuild:
		return "Server"
	case EntityTypeChannel:
		return "<#" + value + "> #" + value
	case EntityTypeMessageCode:
		return value
	}

	return titleify(string(t)) + ": #" + value
}

func titleify(input string) string {
	return strings.Title(strings.Replace(
		strings.TrimPrefix(strings.TrimPrefix(input, "cacophony_"), "discord_"),
		"_", " ", -1))
}
