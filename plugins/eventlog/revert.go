package eventlog

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/permissions"
)

func (i *Item) Revert(event *events.Event) error {
	revertingUser, err := event.State().User(event.UserID)
	if err != nil {
		return err
	}

	reason := event.Translate("eventlog.revert.reason", "item", i, "user", revertingUser)

	switch i.ActionType {
	case ActionTypeDiscordBan:
		if !permissions.DiscordBanMembers.Match(
			event.State(),
			event.DB(),
			event.BotUserID,
			event.ChannelID,
			false,
		) {
			return events.NewUserError("I have insufficient permissions")
		}

		if i.TargetType != EntityTypeUser {
			return fmt.Errorf("invalid target type: %s, expected: %s", i.TargetType, EntityTypeUser)
		}

		err = event.Discord().Client.GuildBanDelete(event.GuildID, i.TargetValue)
		if err != nil {
			return err
		}
	case ActionTypeDiscordUnban:
		if !permissions.DiscordBanMembers.Match(
			event.State(),
			event.DB(),
			event.BotUserID,
			event.ChannelID,
			false,
		) {
			return events.NewUserError("I have insufficient permissions")
		}

		if i.TargetType != EntityTypeUser {
			return fmt.Errorf("invalid target type: %s, expected: %s", i.TargetType, EntityTypeUser)
		}

		err = event.Discord().Client.GuildBanCreateWithReason(event.GuildID, i.TargetValue, reason, 0)
		if err != nil {
			return err
		}
	case ActionTypeGuildUpdate:
		if !permissions.DiscordManageServer.Match(
			event.State(),
			event.DB(),
			event.BotUserID,
			event.ChannelID,
			false,
		) {
			return events.NewUserError("I have insufficient permissions")
		}

		var guildParams discordgo.GuildParams
		var edited bool

		for _, option := range i.Options {
			switch option.Key {
			case "name":
				guildParams.Name = option.PreviousValue
				edited = true
			case "region":
				guildParams.Region = option.PreviousValue
				edited = true
			case "verification_level":
				value, err := strconv.Atoi(option.PreviousValue)
				if err != nil {
					return err
				}

				level := discordgo.VerificationLevel(value)
				guildParams.VerificationLevel = &level
				edited = true
			case "default_message_notifications":
				value, err := strconv.Atoi(option.PreviousValue)
				if err != nil {
					return err
				}

				guildParams.DefaultMessageNotifications = value
				edited = true
			case "afk_channel":
				guildParams.AfkChannelID = option.PreviousValue
				edited = true
			case "afk_timeout":
				value, err := strconv.Atoi(option.PreviousValue)
				if err != nil {
					return err
				}

				guildParams.AfkTimeout = value
				edited = true
			}
		}

		// TODO: revert icon
		// TODO: revert splash

		if !edited {
			return events.NewUserError("no revertable value found")
		}

		_, err = event.Discord().Client.GuildEdit(i.GuildID, guildParams)
		if err != nil {
			return err
		}
	case ActionTypeChannelUpdate:
		if !permissions.DiscordManageChannels.Match(
			event.State(),
			event.DB(),
			event.BotUserID,
			event.ChannelID,
			false,
		) {
			return events.NewUserError("I have insufficient permissions")
		}

		var channelParams discordgo.ChannelEdit
		var edited bool

		for _, option := range i.Options {
			switch option.Key {
			case "name":
				channelParams.Name = option.PreviousValue
				edited = true
			case "topic":
				channelParams.Topic = option.PreviousValue
				edited = true
			case "bitrate":
				value, err := strconv.Atoi(option.PreviousValue)
				if err != nil {
					return err
				}

				channelParams.Bitrate = value
				edited = true
			case "user_limit":
				value, err := strconv.Atoi(option.PreviousValue)
				if err != nil {
					return err
				}

				channelParams.UserLimit = value
				edited = true
			case "rate_limit_per_user":
				value, err := strconv.Atoi(option.PreviousValue)
				if err != nil {
					return err
				}

				channelParams.RateLimitPerUser = value
				edited = true
			case "parent":
				channelParams.ParentID = option.PreviousValue
				edited = true
			case "permissions":
				var permissions []*discordgo.PermissionOverwrite
				err := json.Unmarshal([]byte(option.PreviousValue), &permissions)
				if err != nil {
					return errors.Wrap(err, "unable to parse permission overwrites")
				}

				channelParams.PermissionOverwrites = permissions
				edited = true
			}
		}

		if !edited {
			return events.NewUserError("no revertable value found")
		}

		targetChannel, err := event.State().Channel(i.TargetValue)
		if err != nil {
			return err
		}
		channelParams.Position = targetChannel.Position // prevent resetting position

		_, err = event.Discord().Client.ChannelEditComplex(i.TargetValue, &channelParams)
		if err != nil {
			return err
		}

	default:
		return events.NewUserError("action not revertable")
	}

	err = markItemAsReverted(event.DB(), nil, event.GuildID, i.ID)
	if err != nil {
		return err
	}

	return CreateOptionForItem(event.DB(), event.Publisher(), i.ID, i.GuildID, &ItemOption{
		Key:           "reverted_by",
		PreviousValue: "",
		NewValue:      event.UserID,
		Type:          EntityTypeUser,
	})
}
