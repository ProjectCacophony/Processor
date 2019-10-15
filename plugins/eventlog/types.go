package eventlog

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/state"
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

	EntityTypeText entityType = "text"
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
		// TODO: look up parent
		return "<#" + value + "> #" + value
	case EntityTypeMessageCode:
		return value
	case EntityTypeText:
		return value
	}

	return titleify(string(t)) + ": #" + value
}

func (t entityType) StringWithoutMention(state *state.State, guildID, value string) string {
	switch t {
	case EntityTypeUser:
		user, err := state.User(value)
		if err != nil {
			user = &discordgo.User{
				ID:       value,
				Username: "N/A",
			}
		}
		return user.String() + " #" + value
	case EntityTypeRole:
		role, err := state.Role(guildID, value)
		if err != nil {
			role = &discordgo.Role{
				ID:   value,
				Name: "N/A",
			}
		}
		return role.Name + " #" + value
	case EntityTypeGuild:
		return "Server"
	case EntityTypeChannel:
		channel, err := state.Channel(value)
		if err != nil {
			channel = &discordgo.Channel{
				ID:   value,
				Name: "N/A",
			}
		}
		// TODO: look up parent
		return channel.Name + " #" + value
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
