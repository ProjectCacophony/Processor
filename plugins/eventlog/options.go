package eventlog

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func optionsForChannel(old, new *discordgo.Channel) []ItemOption {
	var options []ItemOption

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

	option = ItemOption{
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

	option = ItemOption{
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

	if (old != nil && old.NSFW) || (new != nil && new.NSFW) {
		option = ItemOption{
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

	if (old != nil && old.Bitrate > 0) || (new != nil && new.Bitrate > 0) {
		option = ItemOption{
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

	var oldPermissionOverwrites, newPermissionOverwrites []byte
	if old != nil {
		oldPermissionOverwrites, _ = json.Marshal(old.PermissionOverwrites)
	}
	if new != nil {
		newPermissionOverwrites, _ = json.Marshal(new.PermissionOverwrites)
	}
	if len(oldPermissionOverwrites) > 0 || len(newPermissionOverwrites) > 0 {
		option = ItemOption{
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

	if (old != nil && old.UserLimit > 0) || (new != nil && new.UserLimit > 0) {
		option = ItemOption{
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

	option = ItemOption{
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

	if (old != nil && old.RateLimitPerUser > 0) || (new != nil && new.RateLimitPerUser > 0) {
		option = ItemOption{
			Key:  "rate_limit_per_user",
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

	return filterEqualOptions(options)
}

func optionsForRole(old, new *discordgo.Role) []ItemOption {
	var options []ItemOption

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

	if (old != nil && old.Managed) || (new != nil && new.Managed) {
		option = ItemOption{
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

	if (old != nil && old.Mentionable) || (new != nil && new.Mentionable) {
		option = ItemOption{
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

	if (old != nil && old.Hoist) || (new != nil && new.Hoist) {
		option = ItemOption{
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

	if (old != nil && old.Color > 0) || (new != nil && new.Color > 0) {
		option = ItemOption{
			Key:  "color",
			Type: EntityTypeColor,
		}
		if old != nil {
			option.PreviousValue = strconv.Itoa(old.Color)
		}
		if new != nil {
			option.NewValue = strconv.Itoa(new.Color)
		}
		options = append(options, option)
	}

	option = ItemOption{
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

	return filterEqualOptions(options)
}

func optionsForGuild(old, new *discordgo.Guild) []ItemOption {
	var options []ItemOption
	if old == nil || new == nil {
		return options
	}

	options = append(options, ItemOption{
		Key:           "name",
		PreviousValue: old.Name,
		NewValue:      new.Name,
		Type:          EntityTypeText,
	})

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

	options = append(options, ItemOption{
		Key:           "region",
		PreviousValue: old.Region,
		NewValue:      new.Region,
		Type:          EntityTypeText,
	})

	options = append(options, ItemOption{
		Key:           "afk_channel",
		PreviousValue: old.AfkChannelID,
		NewValue:      new.AfkChannelID,
		Type:          EntityTypeChannel,
	})

	// TODO: not sent with normal state data?
	// 	options = append(options, ItemOption{
	// 		Key:           "embed_channel",
	// 		PreviousValue: old.EmbedChannelID,
	// 		NewValue:      new.EmbedChannelID,
	// 		Type:          EntityTypeChannel,
	// 	})

	options = append(options, ItemOption{
		Key:           "owner",
		PreviousValue: old.OwnerID,
		NewValue:      new.OwnerID,
		Type:          EntityTypeUser,
	})

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

	options = append(options, ItemOption{
		Key:           "afk_timeout",
		PreviousValue: strconv.Itoa(old.AfkTimeout),
		NewValue:      strconv.Itoa(new.AfkTimeout),
		Type:          EntityTypeNumber, // TODO: duration type? seconds or what unit?
	})

	options = append(options, ItemOption{
		Key:           "verification_level",
		PreviousValue: strconv.Itoa(int(old.VerificationLevel)),
		NewValue:      strconv.Itoa(int(new.VerificationLevel)),
		Type:          EntityTypeGuildVerificationLevel,
	})
	// TODO: not sent with normal state data?
	// 	options = append(options, ItemOption{
	// 		Key:           "embed_enabled",
	// 		PreviousValue: strconv.FormatBool(old.EmbedEnabled),
	// 		NewValue:      strconv.FormatBool(new.EmbedEnabled),
	// 		Type:          EntityTypeBool,
	// 	})

	options = append(options, ItemOption{
		Key:           "default_message_notifications",
		PreviousValue: strconv.Itoa(old.DefaultMessageNotifications),
		NewValue:      strconv.Itoa(new.DefaultMessageNotifications),
		Type:          EntityTypeNumber,
	})

	options = append(options, ItemOption{
		Key:           "explicit_content_filter",
		PreviousValue: strconv.Itoa(int(old.ExplicitContentFilter)),
		NewValue:      strconv.Itoa(int(new.ExplicitContentFilter)),
		Type:          EntityTypeGuildExplicitContentLevel,
	})

	options = append(options, ItemOption{
		Key:           "mfa_level",
		PreviousValue: strconv.Itoa(int(old.MfaLevel)),
		NewValue:      strconv.Itoa(int(new.MfaLevel)),
		Type:          EntityTypeGuildMfaLevel,
	})

	// TODO: not sent with normal state data?
	// 	options = append(options, ItemOption{
	// 		Key:           "widget_enabled",
	// 		PreviousValue: strconv.FormatBool(old.WidgetEnabled),
	// 		NewValue:      strconv.FormatBool(new.WidgetEnabled),
	// 		Type:          EntityTypeBool,
	// 	})

	// TODO: not sent with normal state data?
	// 	options = append(options, ItemOption{
	// 		Key:           "widget_channel_id",
	// 		PreviousValue: old.WidgetChannelID,
	// 		NewValue:      new.WidgetChannelID,
	// 		Type:          EntityTypeChannel,
	// 	})

	options = append(options, ItemOption{
		Key:           "system_channel_id",
		PreviousValue: old.SystemChannelID,
		NewValue:      new.SystemChannelID,
		Type:          EntityTypeChannel,
	})

	options = append(options, ItemOption{
		Key:           "vanity_url_code",
		PreviousValue: old.VanityURLCode,
		NewValue:      new.VanityURLCode,
		Type:          EntityTypeDiscordInvite,
	})

	options = append(options, ItemOption{
		Key:           "description",
		PreviousValue: old.Description,
		NewValue:      new.Description,
		Type:          EntityTypeText,
	})

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

	return filterEqualOptions(options)
}

func optionsForMember(old, new *discordgo.Member) []ItemOption {
	var options []ItemOption
	if old == nil || new == nil {
		return options
	}

	options = append(options, ItemOption{
		Key:           "nick",
		PreviousValue: old.Nick,
		NewValue:      new.Nick,
		Type:          EntityTypeText,
	})

	// TODO: not sent in diff
	// 	options = append(options, ItemOption{
	// 		Key:           "deaf",
	// 		PreviousValue: strconv.FormatBool(old.Deaf),
	// 		NewValue:      strconv.FormatBool(new.Deaf),
	// 		Type:          EntityTypeBool,
	// 	})
	// 	options = append(options, ItemOption{
	// 		Key:           "mute",
	// 		PreviousValue: strconv.FormatBool(old.Mute),
	// 		NewValue:      strconv.FormatBool(new.Mute),
	// 		Type:          EntityTypeBool,
	// 	})

	if !stringSliceMatches(old.Roles, new.Roles) {
		options = append(options, ItemOption{
			Key:           "roles",
			PreviousValue: strings.Join(old.Roles, ","),
			NewValue:      strings.Join(new.Roles, ","),
			Type:          EntityTypeRole,
		})
	}

	// TODO: not sent in diff
	// 	options = append(options, ItemOption{
	// 		Key:           "username",
	// 		PreviousValue: old.User.Username,
	// 		NewValue:      new.User.Username,
	// 		Type:          EntityTypeText,
	// 	})
	// 	options = append(options, ItemOption{
	// 		Key:           "discriminator",
	// 		PreviousValue: old.User.Discriminator,
	// 		NewValue:      new.User.Discriminator,
	// 		Type:          EntityTypeText,
	// 	})

	return filterEqualOptions(options)
}

func optionsForEmoji(old, new *discordgo.Emoji) []ItemOption {
	var options []ItemOption

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

	option = ItemOption{
		Key:  "roles",
		Type: EntityTypeRole,
	}
	if old != nil {
		option.PreviousValue = strings.Join(old.Roles, ",")
	}
	if new != nil {
		option.NewValue = strings.Join(new.Roles, ",")
	}
	options = append(options, option)

	if (old != nil && old.Managed) || (new != nil && new.Managed) {
		option = ItemOption{
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

	return filterEqualOptions(options)
}

func optionsForWebhook(old, new *discordgo.Webhook) []ItemOption {
	var options []ItemOption

	option := ItemOption{
		Key:  "channel",
		Type: EntityTypeChannel,
	}
	if old != nil {
		option.PreviousValue = old.ChannelID
	}
	if new != nil {
		option.NewValue = new.ChannelID
	}
	options = append(options, option)

	option = ItemOption{
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

	option = ItemOption{
		Key:  "avatar",
		Type: EntityTypeImageURL,
	}
	if old != nil && old.Avatar != "" {
		option.PreviousValue = discordgo.EndpointUserAvatar(old.ID, old.Avatar)
	}
	if new != nil && new.Avatar != "" {
		option.NewValue = discordgo.EndpointUserAvatar(new.ID, new.Avatar)
	}
	options = append(options, option)

	return filterEqualOptions(options)
}

func optionsForInvite(old, new *discordgo.Invite) []ItemOption {
	var options []ItemOption

	var option ItemOption

	if (old != nil && old.MaxAge > 0) || (new != nil && new.MaxAge > 0) {
		option = ItemOption{
			Key:  "max_age",
			Type: EntityTypeSeconds,
		}
		if old != nil {
			option.PreviousValue = strconv.Itoa(old.MaxAge)
		}
		if new != nil {
			option.NewValue = strconv.Itoa(new.MaxAge)
		}
		options = append(options, option)
	}

	if (old != nil && old.MaxUses > 0) || (new != nil && new.MaxUses > 0) {
		option = ItemOption{
			Key:  "max_uses",
			Type: EntityTypeNumber,
		}
		if old != nil {
			option.PreviousValue = strconv.Itoa(old.MaxUses)
		}
		if new != nil {
			option.NewValue = strconv.Itoa(new.MaxUses)
		}
		options = append(options, option)
	}

	if (old != nil && old.Revoked) || (new != nil && new.Revoked) {
		option = ItemOption{
			Key:  "revoked",
			Type: EntityTypeBool,
		}
		if old != nil {
			option.PreviousValue = strconv.FormatBool(old.Revoked)
		}
		if new != nil {
			option.NewValue = strconv.FormatBool(new.Revoked)
		}
		options = append(options, option)
	}

	if (old != nil && old.Temporary) || (new != nil && new.Temporary) {
		option = ItemOption{
			Key:  "temporary",
			Type: EntityTypeBool,
		}
		if old != nil {
			option.PreviousValue = strconv.FormatBool(old.Temporary)
		}
		if new != nil {
			option.NewValue = strconv.FormatBool(new.Temporary)
		}
		options = append(options, option)
	}

	if (old != nil && old.Unique) || (new != nil && new.Unique) {
		option = ItemOption{
			Key:  "unique",
			Type: EntityTypeBool,
		}
		if old != nil {
			option.PreviousValue = strconv.FormatBool(old.Unique)
		}
		if new != nil {
			option.NewValue = strconv.FormatBool(new.Unique)
		}
		options = append(options, option)
	}

	return filterEqualOptions(options)
}

func filterEqualOptions(input []ItemOption) []ItemOption {
	result := make([]ItemOption, 0, len(input))
	for _, item := range input {
		if item.PreviousValue == item.NewValue {
			continue
		}

		result = append(result, item)
	}

	return result
}
