package list

import (
	"gitlab.com/Cacophony/Processor/plugins/automod/actions"
	"gitlab.com/Cacophony/Processor/plugins/automod/filters"
	"gitlab.com/Cacophony/Processor/plugins/automod/interfaces"
	"gitlab.com/Cacophony/Processor/plugins/automod/triggers"
)

var (
	TriggerList = []interfaces.TriggerInterface{
		triggers.Message{},
		triggers.BucketUpdated{},
		triggers.Join{},
		triggers.Leave{},
		triggers.Ban{},
		triggers.Unban{},
	}

	FiltersList = []interfaces.FilterInterface{
		filters.RegexMessageContent{},
		filters.True{},
		filters.BucketAmount{},
		filters.RegexUserName{},
		filters.AccountAge{},
		filters.MentionsCount{},
		filters.EmojiCount{},
		filters.ChannelID{},
		filters.AttachmentsCount{},
		filters.RoleID{},
		filters.UserID{},
		filters.InvitesCount{},
		filters.Wait{},
	}

	ActionsList = []interfaces.ActionInterface{
		actions.SendMessage{},
		actions.ApplyRole{},
		actions.IncrBucket{},
		actions.SendMessageTo{},
		actions.DeleteMessage{},
		actions.BanUser{},
		actions.KickUser{},
		actions.ResetBucket{},
		actions.React{},
	}
)
