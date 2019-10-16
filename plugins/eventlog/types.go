package eventlog

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"github.com/pkg/errors"
	"gitlab.com/Cacophony/go-kit/state"
)

// Actions
const (
	ActionTypeModDM         actionType = "cacophony_mod_dm"
	ActionTypeDiscordBan    actionType = "discord_ban"
	ActionTypeDiscordUnban  actionType = "discord_unban"
	ActionTypeDiscordJoin   actionType = "discord_join"
	ActionTypeDiscordLeave  actionType = "discord_leave"
	ActionTypeChannelCreate actionType = "discord_channel_create"
)

func (t actionType) String() string {
	switch t {
	case ActionTypeModDM:
		return "Mod DM"
	}

	return titleify(string(t))
}

func (t actionType) Destructive() bool {
	switch t {
	case ActionTypeDiscordBan,
		ActionTypeDiscordLeave:
		return true
	}

	return false
}

func (t actionType) Revertable() bool {
	switch t {
	case ActionTypeDiscordBan,
		ActionTypeDiscordUnban:
		return true
	}

	return false
}

// Entities
const (
	EntityTypeUser    entityType = "discord_user"
	EntityTypeRole    entityType = "discord_role"
	EntityTypeGuild   entityType = "discord_guild"
	EntityTypeChannel entityType = "discord_channel"

	EntityTypeChannelType          entityType = "discord_channel_type"          // TODO: implement formatting
	EntityTypePermissionOverwrites entityType = "discord_permission_overwrites" // TODO: implement formatting

	EntityTypeMessageCode entityType = "cacophony_message_code"

	EntityTypeText   entityType = "text"
	EntityTypeNumber entityType = "number"
	EntityTypeBool   entityType = "bool"
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
	case EntityTypeBool:
		parsed, _ := strconv.ParseBool(value)
		if parsed {
			return "Yes"
		}
		return "No"
	case EntityTypeNumber:
		parsed, _ := strconv.Atoi(value)

		return humanize.Comma(int64(parsed))
	case EntityTypeChannelType:
		parsed, _ := strconv.Atoi(value)
		switch discordgo.ChannelType(parsed) {
		case discordgo.ChannelTypeGuildText:
			return "Text"
		case discordgo.ChannelTypeGuildVoice:
			return "Voice"
		case discordgo.ChannelTypeGuildCategory:
			return "Category"
		case discordgo.ChannelTypeGuildNews:
			return "News"
		case discordgo.ChannelTypeGuildStore:
			return "Store"
		}
	case EntityTypePermissionOverwrites:
		var permissions []*discordgo.PermissionOverwrite
		err := json.Unmarshal([]byte(value), &permissions)
		if err != nil {
			return errors.Wrap(err, "unable to parse").Error()
		}
		var text string
		for _, permission := range permissions {
			if permission.Allow == 0 && permission.Deny == 0 {
				continue
			}

			switch permission.Type {
			case "role":
				text += "<@&" + permission.ID + "> "
			case "member":
				text += "<@" + permission.ID + "> "
			default:
				text += permission.Type + " #" + permission.ID
			}

			if permission.Allow > 0 {
				text += "Allow " + permissionsText(permission.Allow)
			}
			if permission.Deny > 0 {
				text += "Deny " + permissionsText(permission.Deny)
			}

			text += "; "
		}
		return strings.Trim(text, "; ")
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
