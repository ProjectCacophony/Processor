package eventlog

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleModEvent(event *events.Event) {
	if !isEnabled(event) {
		return
	}

	var items []*Item

	switch event.Type {
	case events.GuildBanAddType:
		items = append(items, &Item{
			GuildID:                    event.GuildBanAdd.GuildID,
			ActionType:                 ActionTypeDiscordBan,
			TargetType:                 EntityTypeUser,
			TargetValue:                event.GuildBanAdd.User.ID,
			WaitingForAuditLogBackfill: true,
		})
	case events.GuildBanRemoveType:
		items = append(items, &Item{
			GuildID:                    event.GuildBanRemove.GuildID,
			ActionType:                 ActionTypeDiscordUnban,
			TargetType:                 EntityTypeUser,
			TargetValue:                event.GuildBanRemove.User.ID,
			WaitingForAuditLogBackfill: true,
		})
	case events.GuildMemberAddType:
		items = append(items, &Item{
			GuildID:     event.GuildMemberAdd.GuildID,
			ActionType:  ActionTypeDiscordJoin,
			TargetType:  EntityTypeUser,
			TargetValue: event.GuildMemberAdd.User.ID,
		})
	// case events.CacophonyGuildMemberAddExtra:
	// 	var options []ItemOption
	// 	if event.GuildMemberAddExtra.UsedInvite != nil {
	// 		if event.GuildMemberAddExtra.UsedInvite.Code != "" {
	// 			options = append(options, ItemOption{
	// 				Key:      "used_invite",
	// 				NewValue: event.GuildMemberAddExtra.UsedInvite.Code,
	// 				Type:     EntityTypeDiscordInvite,
	// 			})
	// 		}
	// 		if event.GuildMemberAddExtra.UsedInvite.Inviter != nil &&
	// 			event.GuildMemberAddExtra.UsedInvite.Inviter.ID != "" {
	// 			options = append(options, ItemOption{
	// 				Key:      "used_invite_author",
	// 				NewValue: event.GuildMemberAddExtra.UsedInvite.Inviter.ID,
	// 				Type:     EntityTypeUser,
	// 			})
	// 		}
	// 	}
	// 	items = append(items, &Item{
	// 		GuildID:     event.GuildMemberAddExtra.GuildID,
	// 		ActionType:  ActionTypeDiscordJoin,
	// 		TargetType:  EntityTypeUser,
	// 		TargetValue: event.GuildMemberAddExtra.User.ID,
	// 		Options:     options,
	// 	})
	case events.GuildMemberRemoveType:
		items = append(items, &Item{
			GuildID:                    event.GuildMemberRemove.GuildID,
			ActionType:                 ActionTypeDiscordLeave,
			TargetType:                 EntityTypeUser,
			TargetValue:                event.GuildMemberRemove.User.ID,
			WaitingForAuditLogBackfill: true,
		})
	case events.ChannelCreateType:
		if event.ChannelCreate.GuildID == "" {
			return
		}

		items = append(items, &Item{
			GuildID:                    event.ChannelCreate.GuildID,
			ActionType:                 ActionTypeChannelCreate,
			TargetType:                 EntityTypeChannel,
			TargetValue:                event.ChannelCreate.ID,
			WaitingForAuditLogBackfill: true,
			Options:                    optionsForChannel(nil, event.ChannelCreate.Channel),
		})
	case events.GuildRoleCreateType:
		items = append(items, &Item{
			GuildID:                    event.GuildRoleCreate.GuildID,
			ActionType:                 ActionTypeRoleCreate,
			TargetType:                 EntityTypeRole,
			TargetValue:                event.GuildRoleCreate.Role.ID,
			WaitingForAuditLogBackfill: true,
			Options:                    optionsForRole(nil, event.GuildRoleCreate.Role),
		})
	case events.CacophonyDiffGuild:
		options := optionsForGuild(event.DiffGuild.Old, event.DiffGuild.New)
		if len(options) <= 0 {
			return
		}

		items = append(items, &Item{
			GuildID:                    event.GuildID,
			ActionType:                 ActionTypeGuildUpdate,
			TargetType:                 EntityTypeGuild,
			TargetValue:                event.GuildID,
			WaitingForAuditLogBackfill: true,
			Options:                    options,
		})
	case events.CacophonyDiffMember:
		options := optionsForMember(event.DiffMember.Old, event.DiffMember.New)
		if len(options) <= 0 {
			return
		}

		items = append(items, &Item{
			GuildID:                    event.GuildID,
			ActionType:                 ActionTypeMemberUpdate,
			TargetType:                 EntityTypeUser,
			TargetValue:                event.DiffMember.Old.User.ID,
			WaitingForAuditLogBackfill: true,
			Options:                    options,
		})
	case events.CacophonyDiffChannel:
		options := optionsForChannel(event.DiffChannel.Old, event.DiffChannel.New)
		if event.DiffChannel.New == nil {
			options = optionsForChannel(nil, event.DiffChannel.Old)
		}
		if len(options) <= 0 {
			return
		}

		item := &Item{
			GuildID:                    event.GuildID,
			ActionType:                 ActionTypeChannelUpdate,
			TargetType:                 EntityTypeChannel,
			TargetValue:                event.DiffChannel.Old.ID,
			WaitingForAuditLogBackfill: true,
			Options:                    options,
		}
		if event.DiffChannel.New == nil {
			item.ActionType = ActionTypeChannelDelete
		}
		items = append(items, item)
	case events.CacophonyDiffRole:
		options := optionsForRole(event.DiffRole.Old, event.DiffRole.New)
		if event.DiffRole.New == nil {
			options = optionsForRole(nil, event.DiffRole.Old)
		}
		if len(options) <= 0 {
			return
		}

		item := &Item{
			GuildID:                    event.GuildID,
			ActionType:                 ActionTypeRoleUpdate,
			TargetType:                 EntityTypeRole,
			TargetValue:                event.DiffRole.Old.ID,
			WaitingForAuditLogBackfill: true,
			Options:                    options,
		}
		if event.DiffRole.New == nil {
			item.ActionType = ActionTypeRoleDelete
		}
		items = append(items, item)
	case events.CacophonyDiffEmoji:
		new, updated, deleted := compareEmojiDiff(event.DiffEmoji)
		for _, emoji := range new {
			items = append(items, &Item{
				GuildID:                    event.GuildID,
				ActionType:                 ActionTypeEmojiCreate,
				TargetType:                 EntityTypeEmoji,
				TargetValue:                emoji.ID,
				WaitingForAuditLogBackfill: true,
				Options:                    optionsForEmoji(nil, emoji),
			})
		}
		for _, emoji := range updated {
			items = append(items, &Item{
				GuildID:                    event.GuildID,
				ActionType:                 ActionTypeEmojiUpdate,
				TargetType:                 EntityTypeEmoji,
				TargetValue:                emoji[1].ID,
				WaitingForAuditLogBackfill: true,
				Options:                    optionsForEmoji(emoji[0], emoji[1]),
			})
		}
		for _, emoji := range deleted {
			items = append(items, &Item{
				GuildID:                    event.GuildID,
				ActionType:                 ActionTypeEmojiDelete,
				TargetType:                 EntityTypeEmoji,
				TargetValue:                emoji.ID,
				WaitingForAuditLogBackfill: true,
				Options:                    optionsForEmoji(nil, emoji),
			})
		}
	case events.CacophonyDiffWebhooks:
		new, updated, deleted := compareWebhooksDiff(event.DiffWebhooks)
		for _, webhook := range new {
			items = append(items, &Item{
				GuildID:                    event.GuildID,
				ActionType:                 ActionTypeWebhookCreate,
				TargetType:                 EntityTypeWebhook,
				TargetValue:                webhook.ID,
				WaitingForAuditLogBackfill: true,
				Options:                    optionsForWebhook(nil, webhook),
			})
		}
		for _, webhook := range updated {
			items = append(items, &Item{
				GuildID:                    event.GuildID,
				ActionType:                 ActionTypeWebhookUpdate,
				TargetType:                 EntityTypeWebhook,
				TargetValue:                webhook[1].ID,
				WaitingForAuditLogBackfill: true,
				Options:                    optionsForWebhook(webhook[0], webhook[1]),
			})
		}
		for _, webhook := range deleted {
			items = append(items, &Item{
				GuildID:                    event.GuildID,
				ActionType:                 ActionTypeWebhookDelete,
				TargetType:                 EntityTypeWebhook,
				TargetValue:                webhook.ID,
				WaitingForAuditLogBackfill: true,
				Options:                    optionsForWebhook(nil, webhook),
			})
		}
	case events.CacophonyDiffInvites:
		var inviterID string

		new, updated, deleted := compareInvitesDiff(event.DiffInvites)
		for _, invite := range new {
			inviterID = ""
			if invite.Inviter != nil {
				inviterID = invite.Inviter.ID
			}
			at, _ := invite.CreatedAt.Parse()

			items = append(items, &Item{
				Model: gorm.Model{
					CreatedAt: at,
				},
				GuildID:     event.GuildID,
				ActionType:  ActionTypeInviteCreate,
				TargetType:  EntityTypeDiscordInvite,
				TargetValue: invite.Code,
				AuthorID:    inviterID,
				Options:     optionsForInvite(nil, invite),
			})
		}
		for _, invite := range updated {
			items = append(items, &Item{
				GuildID:                    event.GuildID,
				ActionType:                 ActionTypeInviteUpdate,
				TargetType:                 EntityTypeDiscordInvite,
				TargetValue:                invite[1].Code,
				WaitingForAuditLogBackfill: true,
				Options:                    optionsForInvite(invite[0], invite[1]),
			})
		}
		for _, invite := range deleted {
			items = append(items, &Item{
				GuildID:                    event.GuildID,
				ActionType:                 ActionTypeInviteDelete,
				TargetType:                 EntityTypeDiscordInvite,
				TargetValue:                invite.Code,
				WaitingForAuditLogBackfill: true,
				Options:                    optionsForInvite(nil, invite),
			})
		}
	}

	for _, item := range items {
		err := CreateItem(event.DB(), event.Publisher(), item)
		if err != nil {
			event.ExceptSilent(err, "action_type", string(item.ActionType))
		}
	}
}
