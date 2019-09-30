package actions

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
	"gitlab.com/Cacophony/go-kit/discord"
)

func automodReason(env *models.Env, action string) string {
	botUserText := "Cacophony"
	botUser, err := env.State.User(env.Event.BotUserID)
	if err == nil {
		botUserText = botUser.String()
	}
	ruleNameText := "Unknown"
	if env.Rule != nil {
		ruleNameText = env.Rule.Name
	}

	return strings.Title(action) + " by " + botUserText + " Automod Rule: " + ruleNameText
}

func ReplaceText(env *models.Env, input string) string {
	// prepare texts
	var userUsername, userFullUsername, userID, userDiscriminator, userMention, userAvatarURL string

	for _, id := range env.UserID {
		user, err := env.State.User(id)
		if err != nil {
			user = &discordgo.User{
				ID:            id,
				Username:      "N/A",
				Discriminator: "N/A",
			}
		}

		userUsername += user.Username + ", "
		userFullUsername += user.String() + ", "
		userID += user.ID + ", "
		userDiscriminator += user.Discriminator + ", "
		userMention += user.Mention() + ", "
		userAvatarURL += user.AvatarURL("2048") + ", "
	}

	userUsername = strings.TrimRight(userUsername, ", ")
	userFullUsername = strings.TrimRight(userFullUsername, ", ")
	userID = strings.TrimRight(userID, ", ")
	userDiscriminator = strings.TrimRight(userDiscriminator, ", ")
	userMention = strings.TrimRight(userMention, ", ")
	userAvatarURL = strings.TrimRight(userAvatarURL, ", ")

	var channelID, channelName string

	for _, id := range env.ChannelID {
		channel, err := env.State.Channel(id)
		if err != nil {
			channel = &discordgo.Channel{
				ID:   id,
				Name: "N/A",
			}
		}

		channelID += channel.ID + ", "
		channelName += channel.Name + ", "
	}

	channelID = strings.TrimRight(channelID, ", ")
	channelName = strings.TrimRight(channelName, ", ")

	var messageID, messageContent, messageLink string

	for _, item := range env.Messages {
		message, err := discord.FindMessage(env.State, nil, item.ChanneID, item.ID)
		if err != nil {
			message = &discordgo.Message{
				ID:        item.ID,
				ChannelID: item.ChanneID,
				Content:   "N/A",
			}
		}

		messageID += message.ID + ", "
		messageContent += discord.MessageCodeFromMessage(message) + ", "
		messageLink += fmt.Sprintf(
			"https://discordapp.com/channels/%s/%s/%s",
			message.GuildID,
			message.ChannelID,
			message.ID,
		) + ", "
	}

	messageID = strings.TrimRight(messageID, ", ")
	messageContent = strings.TrimRight(messageContent, ", ")
	messageLink = strings.TrimRight(messageLink, ", ")

	var guildName, guildID, guildIconURL string

	if env.GuildID != "" {
		guild, err := env.State.Guild(env.GuildID)
		if err != nil {
			guild = &discordgo.Guild{
				ID:   env.GuildID,
				Name: "N/A",
			}
		}

		guildName = guild.Name
		guildID = guild.ID
		if guild.Icon != "" {
			guildIconURL = discordgo.EndpointGuildIcon(env.GuildID, guild.Icon) + "?size=2048"
		}
	}

	// replace texts
	input = strings.Replace(input, "{USER_USERNAME}", userUsername, -1)
	input = strings.Replace(input, "{USER_USERNAME_FULL}", userFullUsername, -1)
	input = strings.Replace(input, "{USER_ID}", userID, -1)
	input = strings.Replace(input, "{USER_DISCRIMINATOR}", userDiscriminator, -1)
	input = strings.Replace(input, "{USER_MENTION}", userMention, -1)
	input = strings.Replace(input, "{USER_AVATARURL}", userAvatarURL, -1) // legacy for Robyul compatibility
	input = strings.Replace(input, "{USER_AVATAR_URL}", userAvatarURL, -1)
	input = strings.Replace(input, "{CHANNEL_ID}", channelID, -1)
	input = strings.Replace(input, "{CHANNEL_NAME}", channelName, -1)
	input = strings.Replace(input, "{GUILD_NAME}", guildName, -1)
	input = strings.Replace(input, "{GUILD_ID}", guildID, -1)
	input = strings.Replace(input, "{GUILD_ICON_URL}", guildIconURL, -1)
	input = strings.Replace(input, "{MESSAGE_ID}", messageID, -1)
	input = strings.Replace(input, "{MESSAGE_CONTENT}", messageContent, -1)
	input = strings.Replace(input, "{MESSAGE_LINK}", messageLink, -1)

	return input
}
