package eventlog

import (
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleModEvent(event *events.Event) {
	switch event.Type {
	case events.GuildBanAddType:
		err := CreateItem(event.DB(), event.Publisher(), &Item{
			GuildID:                    event.GuildBanAdd.GuildID,
			ActionType:                 ActionTypeDiscordBan,
			TargetType:                 EntityTypeUser,
			TargetValue:                event.GuildBanAdd.User.ID,
			WaitingForAuditLogBackfill: true,
		})
		event.ExceptSilent(err)
	case events.GuildBanRemoveType:
		err := CreateItem(event.DB(), event.Publisher(), &Item{
			GuildID:                    event.GuildBanRemove.GuildID,
			ActionType:                 ActionTypeDiscordUnban,
			TargetType:                 EntityTypeUser,
			TargetValue:                event.GuildBanRemove.User.ID,
			WaitingForAuditLogBackfill: true,
		})
		event.ExceptSilent(err)
	case events.GuildMemberAddType:
		err := CreateItem(event.DB(), event.Publisher(), &Item{
			GuildID:     event.GuildMemberAdd.GuildID,
			ActionType:  ActionTypeDiscordJoin,
			TargetType:  EntityTypeUser,
			TargetValue: event.GuildMemberAdd.User.ID,
		})
		event.ExceptSilent(err)
	case events.GuildMemberRemoveType:
		err := CreateItem(event.DB(), event.Publisher(), &Item{
			GuildID:                    event.GuildMemberRemove.GuildID,
			ActionType:                 ActionTypeDiscordLeave,
			TargetType:                 EntityTypeUser,
			TargetValue:                event.GuildMemberRemove.User.ID,
			WaitingForAuditLogBackfill: true,
		})
		event.ExceptSilent(err)
	}
}
