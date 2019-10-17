package eventlog

import (
	"errors"
	"fmt"

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
			return errors.New("insufficient permissions")
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
			return errors.New("insufficient permissions")
		}

		if i.TargetType != EntityTypeUser {
			return fmt.Errorf("invalid target type: %s, expected: %s", i.TargetType, EntityTypeUser)
		}

		err = event.Discord().Client.GuildBanCreateWithReason(event.GuildID, i.TargetValue, reason, 0)
		if err != nil {
			return err
		}
	default:
		return errors.New("action not revertable")
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
