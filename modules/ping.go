package modules

import (
	"fmt"

	"time"

	"gitlab.com/project-d-collab/SqsProcessor/cache"
	"gitlab.com/project-d-collab/dhelpers"
	"gitlab.com/project-d-collab/dhelpers/state"
)

func Action(receivedAt time.Time, event dhelpers.EventContainer) {
	message := "pong, Gateway => SqsProcessor: " + receivedAt.Sub(event.ReceivedAt).String() + "\n"

	channel, err := state.Channel(event.MessageCreate.ChannelID)
	if err == nil {
		message += "channel `" + channel.Name + "`\n"
		guild, err := state.Guild(channel.GuildID)
		if err == nil {
			message += "guild `" + guild.Name + "`\n"
		}
	}
	author, err := state.User(event.MessageCreate.Author.ID)
	if err == nil {
		message += "author `" + author.Username + "#" + author.Discriminator + "`\n"
	}
	bot, err := state.User(event.BotUserID)
	if err == nil {
		message += "bot `" + bot.Username + "#" + bot.Discriminator + "`\n"
	}

	_, err = cache.GetDiscord().ChannelMessageSend(event.MessageCreate.ChannelID, message)
	if err != nil {
		fmt.Println(err.Error())
	}
}
