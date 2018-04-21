package ping

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/project-d-collab/dhelpers"
	"gitlab.com/project-d-collab/dhelpers/state"
)

func simplePing(channelID string, eventReceivedAt time.Time) {
	_, err := dhelpers.SendMessage(channelID, time.Since(eventReceivedAt).String())
	if err != nil {
		fmt.Println(err.Error())
	}
}

func pingInfo(dequeuedAt time.Time, event dhelpers.EventContainer) {
	message := "pong, Gateway => SqsProcessor: " + dequeuedAt.Sub(event.ReceivedAt).String() + "\n"

	var err error
	var channel *discordgo.Channel
	var guild *discordgo.Guild
	var author, bot *discordgo.User
	channel, err = state.Channel(event.MessageCreate.ChannelID)
	if err == nil {
		message += "channel `" + channel.Name + "`\n"
		guild, err = state.Guild(channel.GuildID)
		if err == nil {
			message += "guild `" + guild.Name + "`\n"
		}
	}
	author, err = state.User(event.MessageCreate.Author.ID)
	if err == nil {
		message += "author `" + author.Username + "#" + author.Discriminator + "`\n"
	}
	bot, err = state.User(event.BotUserID)
	if err == nil {
		message += "bot `" + bot.Username + "#" + bot.Discriminator + "`\n"
	}

	message += "\n"

	guildIDs, err := state.AllGuildIDs()
	if err == nil {
		message += "Guilds (" + strconv.Itoa(len(guildIDs)) + "): "
		var botGuild *discordgo.Guild
		for _, guildID := range guildIDs {
			botGuild, err = state.Guild(guildID)
			if err == nil {
				message += "`" + botGuild.Name + "` "
			}
		}
		message += "\n"
	}

	userIDs, err := state.AllUserIDs()
	if err == nil {
		message += "Users (" + strconv.Itoa(len(userIDs)) + "): "
		var botUser *discordgo.User
		for _, userID := range userIDs {
			botUser, err = state.User(userID)
			if err == nil {
				message += "`" + botUser.Username + "#" + botUser.Discriminator + "` "
			}
		}
		message += "\n"
	}

	_, err = dhelpers.SendMessage(event.MessageCreate.ChannelID, message)
	if err != nil {
		fmt.Println(err.Error())
	}
}
