package quickactions

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/discord/emoji"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/permissions"
)

var remindEmoji = map[string]time.Duration{
	"quickaction_remind_1h":  1 * time.Hour,
	"quickaction_remind_8h":  8 * time.Hour,
	"quickaction_remind_24h": 24 * time.Hour,
}

func (p *Plugin) remindMessage(event *events.Event) {
	newEvent, err := events.New(events.CacophonyQuickactionRemind)
	if err != nil {
		event.ExceptSilent(err)
		return
	}
	newEvent.QuickactionRemind = &events.QuickactionRemind{
		GuildID:   event.MessageReactionAdd.GuildID,
		ChannelID: event.MessageReactionAdd.ChannelID,
		MessageID: event.MessageReactionAdd.MessageID,
		Emoji:     &event.MessageReactionAdd.Emoji,
		ToUserID:  event.MessageReactionAdd.UserID,
	}
	newEvent.BotUserID = event.BotUserID

	err = p.publisher.PublishAt(
		event.Context(),
		newEvent,
		time.Now().Add(remindEmoji[event.MessageReactionAdd.Emoji.Name]),
	)
	if err != nil {
		event.ExceptSilent(err)
		return
	}

	ackEmoji := emoji.GetWithout("ok")

	err = discord.React(
		p.redis,
		event.Discord(),
		event.MessageReactionAdd.ChannelID,
		event.MessageReactionAdd.MessageID,
		false,
		ackEmoji,
	)
	if err != nil {
		event.ExceptSilent(err)
		return
	}

	go func() {
		time.Sleep(3 * time.Second)

		// remove reaction if possible
		if permissions.DiscordManageMessages.Match(
			event.State(),
			event.BotUserID,
			event.MessageReactionAdd.ChannelID,
			false,
		) {
			err = event.Discord().Client.MessageReactionRemove(
				event.MessageReactionAdd.ChannelID,
				event.MessageReactionAdd.MessageID,
				event.MessageReactionAdd.Emoji.APIName(),
				event.MessageReactionAdd.UserID,
			)
			if err != nil {
				event.ExceptSilent(err)
				return
			}
		}

		err = event.Discord().Client.MessageReactionRemove(
			event.MessageReactionAdd.ChannelID,
			event.MessageReactionAdd.MessageID,
			ackEmoji,
			event.BotUserID,
		)
		if err != nil {
			event.ExceptSilent(err)
			return
		}
	}()
}

func (p *Plugin) handleQuickactionRemind(event *events.Event) {
	message, err := getMessage(
		event.State(),
		event.Discord(),
		event.QuickactionRemind.ChannelID,
		event.QuickactionRemind.MessageID,
	)
	if err != nil {
		event.ExceptSilent(err)
		return
	}

	message.GuildID = event.QuickactionRemind.GuildID // :(

	// copy message into DMs
	_, err = event.SendComplexDM(
		event.QuickactionRemind.ToUserID,
		&discordgo.MessageSend{
			Content: "quickactions.message.content",
			Embed:   convertMessageToEmbed(message),
		},
		"message",
		message,
		"emoji",
		event.QuickactionRemind.Emoji.MessageFormat(),
	)
	if err != nil {
		event.ExceptSilent(err)
		return
	}
}
