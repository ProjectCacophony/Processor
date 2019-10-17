package eventlog

import (
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleModEvent(event *events.Event) {
	var item *Item

	switch event.Type {
	case events.GuildBanAddType:
		item = &Item{
			GuildID:                    event.GuildBanAdd.GuildID,
			ActionType:                 ActionTypeDiscordBan,
			TargetType:                 EntityTypeUser,
			TargetValue:                event.GuildBanAdd.User.ID,
			WaitingForAuditLogBackfill: true,
		}
	case events.GuildBanRemoveType:
		item = &Item{
			GuildID:                    event.GuildBanRemove.GuildID,
			ActionType:                 ActionTypeDiscordUnban,
			TargetType:                 EntityTypeUser,
			TargetValue:                event.GuildBanRemove.User.ID,
			WaitingForAuditLogBackfill: true,
		}
	case events.GuildMemberAddType:
		item = &Item{
			GuildID:     event.GuildMemberAdd.GuildID,
			ActionType:  ActionTypeDiscordJoin,
			TargetType:  EntityTypeUser,
			TargetValue: event.GuildMemberAdd.User.ID,
		}
	case events.GuildMemberRemoveType:
		item = &Item{
			GuildID:                    event.GuildMemberRemove.GuildID,
			ActionType:                 ActionTypeDiscordLeave,
			TargetType:                 EntityTypeUser,
			TargetValue:                event.GuildMemberRemove.User.ID,
			WaitingForAuditLogBackfill: true,
		}
	case events.ChannelCreateType:
		item = &Item{
			GuildID:                    event.ChannelCreate.GuildID,
			ActionType:                 ActionTypeChannelCreate,
			TargetType:                 EntityTypeChannel,
			TargetValue:                event.ChannelCreate.ID,
			WaitingForAuditLogBackfill: true,
			Options:                    optionsForChannel(event.ChannelCreate.Channel),
		}
	case events.GuildRoleCreateType:
		item = &Item{
			GuildID:                    event.GuildRoleCreate.GuildID,
			ActionType:                 ActionTypeRoleCreate,
			TargetType:                 EntityTypeRole,
			TargetValue:                event.GuildRoleCreate.Role.ID,
			WaitingForAuditLogBackfill: true,
			Options:                    optionsForRole(event.GuildRoleCreate.Role),
		}
	case events.CacophonyDiffGuild:
		options := optionsForGuild(event.DiffGuild.Old, event.DiffGuild.New)
		if len(options) <= 0 {
			return
		}

		item = &Item{
			GuildID:                    event.GuildID,
			ActionType:                 ActionTypeGuildUpdate,
			TargetType:                 EntityTypeGuild,
			TargetValue:                event.GuildID,
			WaitingForAuditLogBackfill: true,
			Options:                    options,
		}
	case events.CacophonyDiffMember:
		options := optionsForMember(event.DiffMember.Old, event.DiffMember.New)
		if len(options) <= 0 {
			return
		}

		item = &Item{
			GuildID:                    event.GuildID,
			ActionType:                 ActionTypeMemberUpdate,
			TargetType:                 EntityTypeUser,
			TargetValue:                event.DiffMember.Old.User.ID,
			WaitingForAuditLogBackfill: true,
			Options:                    options,
		}
	}

	if item != nil {
		err := CreateItem(event.DB(), event.Publisher(), item)
		if err != nil {
			event.ExceptSilent(err)
		}
	}
}
