package eventlog

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func optionsForChannel(old, new *discordgo.Channel) []ItemOption {
	var options []ItemOption

	if (old != nil && old.Name != "") || (old != nil && new != nil && old.Name != new.Name) || (old == nil && new != nil && new.Name != "") {
		option := ItemOption{
			Key:  "name",
			Type: EntityTypeText,
		}
		if old != nil {
			option.PreviousValue = old.Name
		}
		if new != nil {
			option.NewValue = new.Name
		}
		options = append(options, option)
	}
	if (old != nil && old.Topic != "") || (old != nil && new != nil && old.Topic != new.Topic) || (old == nil && new != nil && new.Topic != "") {
		option := ItemOption{
			Key:  "topic",
			Type: EntityTypeText,
		}
		if old != nil {
			option.PreviousValue = old.Topic
		}
		if new != nil {
			option.NewValue = new.Topic
		}
		options = append(options, option)
	}
	if (old != nil && new == nil) || (old != nil && new != nil && old.Type != new.Type) || (old == nil && new != nil) {
		option := ItemOption{
			Key:  "type",
			Type: EntityTypeChannelType,
		}
		if old != nil {
			option.PreviousValue = strconv.Itoa(int(old.Type))
		}
		if new != nil {
			option.NewValue = strconv.Itoa(int(new.Type))
		}
		options = append(options, option)
	}
	if (old != nil && old.NSFW) || (old != nil && new != nil && old.NSFW != new.NSFW) || (old == nil && new != nil && new.NSFW) {
		option := ItemOption{
			Key:  "nsfw",
			Type: EntityTypeBool,
		}
		if old != nil {
			option.PreviousValue = strconv.FormatBool(old.NSFW)
		}
		if new != nil {
			option.NewValue = strconv.FormatBool(new.NSFW)
		}
		options = append(options, option)
	}
	if (old != nil && old.Bitrate > 0) || (old != nil && new != nil && old.Bitrate != new.Bitrate) || (old == nil && new != nil && new.Bitrate > 0) {
		option := ItemOption{
			Key:  "bitrate",
			Type: EntityTypeNumber,
		}
		if old != nil {
			option.PreviousValue = strconv.Itoa(old.Bitrate)
		}
		if new != nil {
			option.NewValue = strconv.Itoa(new.Bitrate)
		}
		options = append(options, option)
	}
	if (old != nil && len(old.PermissionOverwrites) > 0) || (old != nil && new != nil && len(old.PermissionOverwrites) != len(new.PermissionOverwrites)) || (old == nil && new != nil && len(new.PermissionOverwrites) > 0) {
		var oldPermissionOverwrites, newPermissionOverwrites []byte
		if old != nil {
			oldPermissionOverwrites, _ = json.Marshal(old.PermissionOverwrites)
		}
		if new != nil {
			newPermissionOverwrites, _ = json.Marshal(new.PermissionOverwrites)
		}
		if len(oldPermissionOverwrites) > 0 || len(newPermissionOverwrites) > 0 {
			option := ItemOption{
				Key:  "permissions",
				Type: EntityTypePermissionOverwrites,
			}
			if old != nil {
				option.PreviousValue = string(oldPermissionOverwrites)
			}
			if new != nil {
				option.NewValue = string(newPermissionOverwrites)
			}
			options = append(options, option)
		}
	}
	if (old != nil && old.UserLimit > 0) || (old != nil && new != nil && old.UserLimit != new.UserLimit) || (old == nil && new != nil && new.UserLimit > 0) {
		option := ItemOption{
			Key:  "user_limit",
			Type: EntityTypeNumber,
		}
		if old != nil {
			option.PreviousValue = strconv.Itoa(old.UserLimit)
		}
		if new != nil {
			option.NewValue = strconv.Itoa(new.UserLimit)
		}
		options = append(options, option)
	}
	if (old != nil && old.ParentID != "") || (old != nil && new != nil && old.ParentID != new.ParentID) || (old == nil && new != nil && new.ParentID != "") {
		option := ItemOption{
			Key:  "parent",
			Type: EntityTypeChannel,
		}
		if old != nil {
			option.PreviousValue = old.ParentID
		}
		if new != nil {
			option.NewValue = new.ParentID
		}
		options = append(options, option)
	}
	if (old != nil && old.RateLimitPerUser > 0) || (old != nil && new != nil && old.RateLimitPerUser != new.RateLimitPerUser) || (old == nil && new != nil && new.RateLimitPerUser > 0) {
		option := ItemOption{
			Key:  "parent",
			Type: EntityTypeNumber,
		}
		if old != nil {
			option.PreviousValue = strconv.Itoa(old.RateLimitPerUser)
		}
		if new != nil {
			option.NewValue = strconv.Itoa(new.RateLimitPerUser)
		}
		options = append(options, option)
	}

	return options
}

func optionsForRole(old, new *discordgo.Role) []ItemOption {
	var options []ItemOption

	if (old != nil && old.Name != "") || (old != nil && new != nil && old.Name != new.Name) || (old == nil && new != nil && new.Name != "") {
		option := ItemOption{
			Key:  "name",
			Type: EntityTypeText,
		}
		if old != nil {
			option.PreviousValue = old.Name
		}
		if new != nil {
			option.NewValue = new.Name
		}
		options = append(options, option)
	}
	if (old != nil && old.Managed) || (old != nil && new != nil && old.Managed != new.Managed) || (old == nil && new != nil && new.Managed) {
		option := ItemOption{
			Key:  "managed",
			Type: EntityTypeBool,
		}
		if old != nil {
			option.PreviousValue = strconv.FormatBool(old.Managed)
		}
		if new != nil {
			option.NewValue = strconv.FormatBool(new.Managed)
		}
		options = append(options, option)
	}
	if (old != nil && old.Mentionable) || (old != nil && new != nil && old.Mentionable != new.Mentionable) || (old == nil && new != nil && new.Mentionable) {
		option := ItemOption{
			Key:  "mentionable",
			Type: EntityTypeBool,
		}
		if old != nil {
			option.PreviousValue = strconv.FormatBool(old.Mentionable)
		}
		if new != nil {
			option.NewValue = strconv.FormatBool(new.Mentionable)
		}
		options = append(options, option)
	}
	if (old != nil && old.Hoist) || (old != nil && new != nil && old.Hoist != new.Hoist) || (old == nil && new != nil && new.Hoist) {
		option := ItemOption{
			Key:  "hoist",
			Type: EntityTypeBool,
		}
		if old != nil {
			option.PreviousValue = strconv.FormatBool(old.Hoist)
		}
		if new != nil {
			option.NewValue = strconv.FormatBool(new.Hoist)
		}
		options = append(options, option)
	}
	if (old != nil && old.Color > 0) || (old != nil && new != nil && old.Color != new.Color) || (old == nil && new != nil && new.Color > 0) {
		option := ItemOption{
			Key:           "color",
			PreviousValue: strconv.Itoa(old.Color),
			Type:          EntityTypeColor,
		}
		if old != nil {
			option.PreviousValue = strconv.Itoa(old.Color)
		}
		if new != nil {
			option.NewValue = strconv.Itoa(new.Color)
		}
		options = append(options, option)
	}
	if (old != nil && old.Permissions > 0) || (old != nil && new != nil && old.Permissions != new.Permissions) || (old == nil && new != nil && new.Permissions > 0) {
		option := ItemOption{
			Key:  "permission",
			Type: EntityTypePermission,
		}
		if old != nil {
			option.PreviousValue = strconv.Itoa(old.Permissions)
		}
		if new != nil {
			option.NewValue = strconv.Itoa(new.Permissions)
		}
		options = append(options, option)
	}

	return options
}

func optionsForGuild(old, new *discordgo.Guild) []ItemOption {
	var options []ItemOption
	if old == nil || new == nil {
		return options
	}

	if old.Name != new.Name {
		options = append(options, ItemOption{
			Key:           "name",
			PreviousValue: old.Name,
			NewValue:      new.Name,
			Type:          EntityTypeText,
		})
	}
	if old.Icon != new.Icon {
		var oldIcon, newIcon string
		if old.Icon != "" {
			oldIcon = old.IconURL()
		}
		if new.Icon != "" {
			newIcon = new.IconURL()
		}
		options = append(options, ItemOption{
			Key:           "icon",
			PreviousValue: oldIcon,
			NewValue:      newIcon,
			Type:          EntityTypeImageURL,
		})
	}
	if old.Region != new.Region {
		options = append(options, ItemOption{
			Key:           "region",
			PreviousValue: old.Region,
			NewValue:      new.Region,
			Type:          EntityTypeText,
		})
	}
	if old.AfkChannelID != new.AfkChannelID {
		options = append(options, ItemOption{
			Key:           "afk_channel",
			PreviousValue: old.AfkChannelID,
			NewValue:      new.AfkChannelID,
			Type:          EntityTypeChannel,
		})
	}
	// TODO: not sent with normal state data?
	// if old.EmbedChannelID != new.EmbedChannelID {
	// 	options = append(options, ItemOption{
	// 		Key:           "embed_channel",
	// 		PreviousValue: old.EmbedChannelID,
	// 		NewValue:      new.EmbedChannelID,
	// 		Type:          EntityTypeChannel,
	// 	})
	// }
	if old.OwnerID != new.OwnerID {
		options = append(options, ItemOption{
			Key:           "owner",
			PreviousValue: old.OwnerID,
			NewValue:      new.OwnerID,
			Type:          EntityTypeUser,
		})
	}
	if old.Splash != new.Splash {
		var oldSplash, newSplash string
		if old.Splash != "" {
			oldSplash = discordgo.EndpointGuildSplash(old.ID, old.Splash)
		}
		if new.Splash != "" {
			newSplash = discordgo.EndpointGuildSplash(new.ID, new.Splash)
		}
		options = append(options, ItemOption{
			Key:           "splash",
			PreviousValue: oldSplash,
			NewValue:      newSplash,
			Type:          EntityTypeImageURL,
		})
	}
	if old.AfkTimeout != new.AfkTimeout {
		options = append(options, ItemOption{
			Key:           "afk_timeout",
			PreviousValue: strconv.Itoa(old.AfkTimeout),
			NewValue:      strconv.Itoa(new.AfkTimeout),
			Type:          EntityTypeNumber, // TODO: duration type? seconds or what unit?
		})
	}
	if old.VerificationLevel != new.VerificationLevel {
		options = append(options, ItemOption{
			Key:           "verification_level",
			PreviousValue: strconv.Itoa(int(old.VerificationLevel)),
			NewValue:      strconv.Itoa(int(new.VerificationLevel)),
			Type:          EntityTypeGuildVerificationLevel,
		})
	}
	// TODO: not sent with normal state data?
	// if old.EmbedEnabled != new.EmbedEnabled {
	// 	options = append(options, ItemOption{
	// 		Key:           "embed_enabled",
	// 		PreviousValue: strconv.FormatBool(old.EmbedEnabled),
	// 		NewValue:      strconv.FormatBool(new.EmbedEnabled),
	// 		Type:          EntityTypeBool,
	// 	})
	// }
	if old.DefaultMessageNotifications != new.DefaultMessageNotifications {
		options = append(options, ItemOption{
			Key:           "default_message_notifications",
			PreviousValue: strconv.Itoa(old.DefaultMessageNotifications),
			NewValue:      strconv.Itoa(new.DefaultMessageNotifications),
			Type:          EntityTypeNumber,
		})
	}
	if old.ExplicitContentFilter != new.ExplicitContentFilter {
		options = append(options, ItemOption{
			Key:           "explicit_content_filter",
			PreviousValue: strconv.Itoa(int(old.ExplicitContentFilter)),
			NewValue:      strconv.Itoa(int(new.ExplicitContentFilter)),
			Type:          EntityTypeGuildExplicitContentLevel,
		})
	}
	if old.MfaLevel != new.MfaLevel {
		options = append(options, ItemOption{
			Key:           "mfa_level",
			PreviousValue: strconv.Itoa(int(old.MfaLevel)),
			NewValue:      strconv.Itoa(int(new.MfaLevel)),
			Type:          EntityTypeGuildMfaLevel,
		})
	}
	// TODO: not sent with normal state data?
	// if old.WidgetEnabled != new.WidgetEnabled {
	// 	options = append(options, ItemOption{
	// 		Key:           "widget_enabled",
	// 		PreviousValue: strconv.FormatBool(old.WidgetEnabled),
	// 		NewValue:      strconv.FormatBool(new.WidgetEnabled),
	// 		Type:          EntityTypeBool,
	// 	})
	// }
	// TODO: not sent with normal state data?
	// if old.WidgetChannelID != new.WidgetChannelID {
	// 	options = append(options, ItemOption{
	// 		Key:           "widget_channel_id",
	// 		PreviousValue: old.WidgetChannelID,
	// 		NewValue:      new.WidgetChannelID,
	// 		Type:          EntityTypeChannel,
	// 	})
	// }
	if old.SystemChannelID != new.SystemChannelID {
		options = append(options, ItemOption{
			Key:           "system_channel_id",
			PreviousValue: old.SystemChannelID,
			NewValue:      new.SystemChannelID,
			Type:          EntityTypeChannel,
		})
	}
	if old.VanityURLCode != new.VanityURLCode {
		options = append(options, ItemOption{
			Key:           "vanity_url_code",
			PreviousValue: old.VanityURLCode,
			NewValue:      new.VanityURLCode,
			Type:          EntityTypeDiscordInvite,
		})
	}
	if old.Description != new.Description {
		options = append(options, ItemOption{
			Key:           "description",
			PreviousValue: old.Description,
			NewValue:      new.Description,
			Type:          EntityTypeText,
		})
	}
	if old.Banner != new.Banner {
		var oldBanner, newBanner string
		if old.Banner != "" {
			oldBanner = discordgo.EndpointGuildBanner(old.ID, old.Banner)
		}
		if new.Banner != "" {
			newBanner = discordgo.EndpointGuildBanner(new.ID, new.Banner)
		}
		options = append(options, ItemOption{
			Key:           "banner",
			PreviousValue: oldBanner,
			NewValue:      newBanner,
			Type:          EntityTypeImageURL,
		})
	}

	return options
}

func optionsForMember(old, new *discordgo.Member) []ItemOption {
	var options []ItemOption
	if old == nil || new == nil {
		return options
	}

	if old.Nick != new.Nick {
		options = append(options, ItemOption{
			Key:           "nick",
			PreviousValue: old.Nick,
			NewValue:      new.Nick,
			Type:          EntityTypeText,
		})
	}
	// TODO: not sent in diff
	// if old.Deaf != new.Deaf {
	// 	options = append(options, ItemOption{
	// 		Key:           "deaf",
	// 		PreviousValue: strconv.FormatBool(old.Deaf),
	// 		NewValue:      strconv.FormatBool(new.Deaf),
	// 		Type:          EntityTypeBool,
	// 	})
	// }
	// if old.Mute != new.Mute {
	// 	options = append(options, ItemOption{
	// 		Key:           "mute",
	// 		PreviousValue: strconv.FormatBool(old.Mute),
	// 		NewValue:      strconv.FormatBool(new.Mute),
	// 		Type:          EntityTypeBool,
	// 	})
	// }
	if !stringSliceMatches(old.Roles, new.Roles) {
		options = append(options, ItemOption{
			Key:           "roles",
			PreviousValue: strings.Join(old.Roles, ","),
			NewValue:      strings.Join(new.Roles, ","),
			Type:          EntityTypeRole,
		})
	}
	// TODO: not sent in diff
	// if old.User.Username != new.User.Username {
	// 	options = append(options, ItemOption{
	// 		Key:           "username",
	// 		PreviousValue: old.User.Username,
	// 		NewValue:      new.User.Username,
	// 		Type:          EntityTypeText,
	// 	})
	// }
	// if old.User.Discriminator != new.User.Discriminator {
	// 	options = append(options, ItemOption{
	// 		Key:           "discriminator",
	// 		PreviousValue: old.User.Discriminator,
	// 		NewValue:      new.User.Discriminator,
	// 		Type:          EntityTypeText,
	// 	})
	// }

	return options
}
