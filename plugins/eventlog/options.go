package eventlog

import (
	"encoding/json"
	"strconv"
	"strings"

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

func optionsForRole(role *discordgo.Role) []ItemOption {
	var options []ItemOption
	if role == nil {
		return options
	}

	if role.Name != "" {
		options = append(options, ItemOption{
			Key:      "name",
			NewValue: role.Name,
			Type:     EntityTypeText,
		})
	}
	if role.Managed {
		options = append(options, ItemOption{
			Key:      "managed",
			NewValue: strconv.FormatBool(role.Managed),
			Type:     EntityTypeBool,
		})
	}
	if role.Mentionable {
		options = append(options, ItemOption{
			Key:      "mentionable",
			NewValue: strconv.FormatBool(role.Mentionable),
			Type:     EntityTypeBool,
		})
	}
	if role.Hoist {
		options = append(options, ItemOption{
			Key:      "hoist",
			NewValue: strconv.FormatBool(role.Hoist),
			Type:     EntityTypeBool,
		})
	}
	if role.Color > 0 {
		options = append(options, ItemOption{
			Key:      "color",
			NewValue: strconv.Itoa(role.Color),
			Type:     EntityTypeColor,
		})
	}
	if role.Permissions > 0 {
		options = append(options, ItemOption{
			Key:      "permission",
			NewValue: strconv.Itoa(role.Permissions),
			Type:     EntityTypePermission,
		})
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
