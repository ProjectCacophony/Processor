package eventlog

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"github.com/pkg/errors"
	"gitlab.com/Cacophony/go-kit/discord"
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
	ActionTypeRoleCreate    actionType = "discord_role_create"
	ActionTypeGuildUpdate   actionType = "discord_guild_update"
	ActionTypeMemberUpdate  actionType = "discord_member_update"
	ActionTypeChannelUpdate actionType = "discord_channel_update"
	ActionTypeChannelDelete actionType = "discord_channel_delete"
	ActionTypeRoleUpdate    actionType = "discord_role_update"
	ActionTypeRoleDelete    actionType = "discord_role_delete"
	ActionTypeEmojiCreate   actionType = "discord_emoji_create"
	ActionTypeEmojiUpdate   actionType = "discord_emoji_update"
	ActionTypeEmojiDelete   actionType = "discord_emoji_delete"
	ActionTypeWebhookCreate actionType = "discord_webhook_create"
	ActionTypeWebhookUpdate actionType = "discord_webhook_update"
	ActionTypeWebhookDelete actionType = "discord_webhook_delete"
	ActionTypeInviteCreate  actionType = "discord_invite_create"
	ActionTypeInviteUpdate  actionType = "discord_invite_update"
	ActionTypeInviteDelete  actionType = "discord_invite_delete"
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
		ActionTypeDiscordUnban,
		ActionTypeGuildUpdate,
		ActionTypeChannelUpdate,
		ActionTypeChannelDelete,
		ActionTypeRoleUpdate:
		return true
	}

	return false
}

// Entities
const (
	EntityTypeUser                      entityType = "discord_user"
	EntityTypeRole                      entityType = "discord_role"
	EntityTypeGuild                     entityType = "discord_guild"
	EntityTypeChannel                   entityType = "discord_channel"
	EntityTypeEmoji                     entityType = "discord_emoji"
	EntityTypeWebhook                   entityType = "discord_webhook"
	EntityTypePermission                entityType = "discord_permission"
	EntityTypeColor                     entityType = "discord_color"
	EntityTypeChannelType               entityType = "discord_channel_type"
	EntityTypePermissionOverwrites      entityType = "discord_permission_overwrites"
	EntityTypeDiscordInvite             entityType = "discord_invite"
	EntityTypeGuildVerificationLevel    entityType = "discord_guild_verification_level"
	EntityTypeGuildExplicitContentLevel entityType = "discord_guild_explicit_content_level"
	EntityTypeGuildMfaLevel             entityType = "discord_guild_mfa_level"

	EntityTypeImageURL    entityType = "cacophony_image_url" // TODO: cache?
	EntityTypeMessageCode entityType = "cacophony_message_code"

	EntityTypeText    entityType = "text"
	EntityTypeNumber  entityType = "number"
	EntityTypeBool    entityType = "bool"
	EntityTypeSeconds entityType = "seconds" // TODO
)

func (t entityType) String(state *state.State, guildID, value string) string {
	switch t {
	case EntityTypeUser:
		return "<@" + value + "> #" + value
	case EntityTypeRole:
		values := strings.Split(value, ",")
		var result string
		for _, item := range values {
			result += "<@&" + item + "> #" + item
		}
		return strings.Trim(result, ", ")
	case EntityTypeGuild:
		return "Server"
	case EntityTypeChannel:
		text := "<#" + value + "> #" + value

		targetChannel, err := state.Channel(value)
		if err == nil && targetChannel.ParentID != "" {
			text = "<#" + targetChannel.ParentID + "> #" + targetChannel.ParentID + " / " + text
		}

		return text
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
	case EntityTypeColor:
		parsed, _ := strconv.Atoi(value)

		return "#" + discord.ColorCodeToHex(parsed)
	case EntityTypePermission:
		parsed, _ := strconv.Atoi(value)

		return permissionsText(parsed)
	case EntityTypeDiscordInvite:
		return "https://discord.gg/" + value
	case EntityTypeImageURL:
		return value
	case EntityTypeGuildVerificationLevel:
		parsed, _ := strconv.Atoi(value)

		switch discordgo.VerificationLevel(parsed) {
		case discordgo.VerificationLevelNone:
			return "None"
		case discordgo.VerificationLevelLow:
			return "Low"
		case discordgo.VerificationLevelMedium:
			return "Medium"
		case discordgo.VerificationLevelHigh:
			return "High"
		}
	case EntityTypeGuildExplicitContentLevel:
		parsed, _ := strconv.Atoi(value)

		switch discordgo.ExplicitContentFilterLevel(parsed) {
		case discordgo.ExplicitContentFilterDisabled:
			return "Disabled"
		case discordgo.ExplicitContentFilterMembersWithoutRoles:
			return "Members without roles"
		case discordgo.ExplicitContentFilterAllMembers:
			return "All member"
		}
	case EntityTypeGuildMfaLevel:
		parsed, _ := strconv.Atoi(value)

		switch discordgo.MfaLevel(parsed) {
		case discordgo.MfaLevelNone:
			return "None"
		case discordgo.MfaLevelElevated:
			return "Elevated"
		}
	case EntityTypeEmoji:
		guild, err := state.Guild(guildID)
		if err == nil {
			emoji := emojiSliceFindEmoji(value, guild.Emojis)
			if emoji != nil {
				return emoji.MessageFormat()
			}
		}

		return "Emoji #" + value
	case EntityTypeWebhook:
		webhook, err := state.Webhook(value)
		if err == nil {
			return "Webhook " + webhook.Name + " <#" + webhook.ChannelID + "> #" + webhook.ID
		}

		return "Webhook #" + value
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
		values := strings.Split(value, ",")
		var result string
		for _, item := range values {
			role, err := state.Role(guildID, item)
			if err != nil {
				role = &discordgo.Role{
					ID:   value,
					Name: "N/A",
				}
			}
			result += role.Name + " #" + value
		}
		return strings.Trim(result, ", ")
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
	}

	return t.String(state, guildID, value)
}

func titleify(input string) string {
	return strings.Title(strings.Replace(
		strings.TrimPrefix(strings.TrimPrefix(input, "cacophony_"), "discord_"),
		"_", " ", -1))
}
