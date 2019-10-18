package eventlog

import (
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleModEvent(event *events.Event) {
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
	case events.GuildMemberRemoveType:
		items = append(items, &Item{
			GuildID:                    event.GuildMemberRemove.GuildID,
			ActionType:                 ActionTypeDiscordLeave,
			TargetType:                 EntityTypeUser,
			TargetValue:                event.GuildMemberRemove.User.ID,
			WaitingForAuditLogBackfill: true,
		})
	case events.ChannelCreateType:
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

	}

	for _, item := range items {
		err := CreateItem(event.DB(), event.Publisher(), item)
		if err != nil {
			event.ExceptSilent(err)
		}
	}
}
