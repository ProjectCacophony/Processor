package eventlog

import (
	"encoding/json"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func optionsForChannel(channel *discordgo.Channel) []ItemOption {
	var options []ItemOption
	if channel == nil {
		return options
	}

	if channel.Name != "" {
		options = append(options, ItemOption{
			Key:      "name",
			NewValue: channel.Name,
			Type:     EntityTypeText,
		})
	}
	if channel.Topic != "" {
		options = append(options, ItemOption{
			Key:      "topic",
			NewValue: channel.Topic,
			Type:     EntityTypeText,
		})
	}
	options = append(options, ItemOption{
		Key:      "type",
		NewValue: strconv.Itoa(int(channel.Type)),
		Type:     EntityTypeChannelType,
	})
	if channel.NSFW {
		options = append(options, ItemOption{
			Key:      "nsfw",
			NewValue: strconv.FormatBool(channel.NSFW),
			Type:     EntityTypeBool,
		})
	}
	if channel.Bitrate > 0 {
		options = append(options, ItemOption{
			Key:      "bitrate",
			NewValue: strconv.Itoa(channel.Bitrate),
			Type:     EntityTypeNumber,
		})
	}
	if len(channel.PermissionOverwrites) > 0 {
		permissionOverwrites, err := json.Marshal(channel.PermissionOverwrites)
		if err == nil {
			options = append(options, ItemOption{
				Key:      "permissions",
				NewValue: string(permissionOverwrites),
				Type:     EntityTypePermissionOverwrites,
			})
		}
	}
	if channel.UserLimit > 0 {
		options = append(options, ItemOption{
			Key:      "user_limit",
			NewValue: strconv.Itoa(channel.UserLimit),
			Type:     EntityTypeNumber,
		})
	}
	if channel.ParentID != "" {
		options = append(options, ItemOption{
			Key:      "parent",
			NewValue: channel.ParentID,
			Type:     EntityTypeChannel,
		})
	}
	if channel.RateLimitPerUser > 0 {
		options = append(options, ItemOption{
			Key:      "parent",
			NewValue: strconv.Itoa(channel.RateLimitPerUser),
			Type:     EntityTypeNumber,
		})
	}

	return options
}
