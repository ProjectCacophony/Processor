package dev

import (
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
	humanize "github.com/dustin/go-humanize"
	"gitlab.com/Cacophony/go-kit/events"
)

const (
	unavailablePlaceholder = "N/A"
)

func handleDevState(event *events.Event) {
	targetUserID := event.MessageCreate.Author.ID
	if len(event.Fields()) >= 3 {
		targetUserID = event.Fields()[2]
	}

	allGuilds, err := event.State().AllGuildIDs()
	if err != nil {
		event.Except(err)
		return
	}
	allChannels, err := event.State().AllChannelIDs()
	if err != nil {
		event.Except(err)
		return
	}
	allUsers, err := event.State().AllUserIDs()
	if err != nil {
		event.Except(err)
		return
	}

	user, _ := event.State().User(targetUserID)
	usernameText := unavailablePlaceholder
	if user != nil {
		usernameText = user.String()
	}

	member, _ := event.State().Member(
		event.MessageCreate.GuildID,
		targetUserID,
	)
	memberRolesText := unavailablePlaceholder
	memberJoinedAtText := unavailablePlaceholder
	if member != nil {
		memberRolesText = strconv.Itoa(len(member.Roles))
		memberJoinedAt, _ := member.JoinedAt.Parse()
		if err != nil {
			event.Except(err)
			return
		}
		memberJoinedAtText = humanize.Time(memberJoinedAt)
	}
	memberIs, err := event.State().IsMember(
		event.MessageCreate.GuildID,
		targetUserID,
	)
	if err != nil {
		event.Except(err)
		return
	}
	botID, err := event.State().BotForGuild(event.MessageCreate.GuildID)
	if err != nil {
		event.Except(err)
		return
	}

	_, err = event.RespondComplex(&discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title: "State :spy:",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Guilds",
					Value:  fmt.Sprintf("**%d**", len(allGuilds)),
					Inline: true,
				},
				{
					Name:   "Channels",
					Value:  fmt.Sprintf("**%d**", len(allChannels)),
					Inline: true,
				},
				{
					Name:   "Users",
					Value:  fmt.Sprintf("**%d**", len(allUsers)),
					Inline: true,
				},
				{
					Name: "User",
					Value: fmt.Sprintf("**%s**",
						usernameText,
					),
					Inline: true,
				},
				{
					Name: "Member on this Guild",
					Value: fmt.Sprintf("Joined **%s**\nRoles **%s**\nMember **%v**",
						memberJoinedAtText,
						memberRolesText,
						memberIs,
					),
					Inline: true,
				},
				{
					Name:   "Bot",
					Value:  fmt.Sprintf("<@%s>", botID),
					Inline: true,
				},
			},
		},
	})
	event.Except(err)
}

func handleDevStateGuilds(event *events.Event) {
	allGuilds, err := event.State().AllGuildIDs()
	if err != nil {
		event.Except(err)
		return
	}

	var resp string
	for _, guildID := range allGuilds {
		guild, err := event.State().Guild(guildID)
		if err != nil {
			event.Except(err)
			return
		}

		resp += fmt.Sprintf("**%s** (`#%s`)\n", guild.Name, guild.ID)
	}
	resp += fmt.Sprintf("in total **%d** guilds.\n", len(allGuilds))

	_, err = event.Respond(resp)
	event.Except(err)
}
