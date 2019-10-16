package eventlog

import (
	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/config"
	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/permissions"
)

func (p *Plugin) handleEventlogUpdate(event *events.Event) {
	item, err := GetItem(event.DB(), event.EventlogUpdate.ItemID)
	if err != nil {
		event.ExceptSilent(err)
		return
	}

	channelID, err := config.GuildGetString(p.db, event.EventlogUpdate.GuildID, eventlogChannelKey)
	if err != nil {
		event.ExceptSilent(err)
		return
	}

	botID, err := event.State().BotForChannel(channelID, permissions.DiscordSendMessages, permissions.DiscordEmbedLinks, permissions.DiscordAddReactions)
	if err != nil {
		event.ExceptSilent(err)
		return
	}
	event.BotUserID = botID

	embed := item.Embed(event.State())

	if item.LogMessage.MessageID != "" && item.LogMessage.ChannelID != "" {
		_, err = discord.EditComplexWithVars(
			event.Redis(),
			event.Discord(),
			event.Localizations(),
			&discordgo.MessageEdit{
				Embed:   embed,
				ID:      item.LogMessage.MessageID,
				Channel: item.LogMessage.ChannelID,
			},
			false,
		)
		event.ExceptSilent(err)
		return
	}

	messages, err := event.SendComplex(
		channelID,
		&discordgo.MessageSend{Embed: embed},
	)
	if err != nil {
		event.ExceptSilent(err)
		return
	}

	err = saveItemMessage(event.DB(), item.ID, messages[0].ID, messages[0].ChannelID)
	if err != nil {
		event.ExceptSilent(err)
		return
	}

	discord.React(event.Redis(), event.Discord(), messages[0].ChannelID, messages[0].ID, false, reactionEditReason)

	if item.ActionType.Revertable() {
		discord.React(event.Redis(), event.Discord(), messages[0].ChannelID, messages[0].ID, false, reactionRevert)
	}
}
