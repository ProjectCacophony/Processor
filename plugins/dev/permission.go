package dev

import (
	"strconv"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/permissions"

	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleDevPermission(event *events.Event) {
	if len(event.Fields()) < 3 {
		return
	}

	user, err := event.FindMember()
	if err != nil {
		event.Except(err)
		return
	}

	permissionID, err := strconv.Atoi(event.Fields()[2])
	if err != nil {
		event.Except(err)
		return
	}

	var channel *discordgo.Channel
	for _, field := range event.Fields() {
		channel, err = event.State().ChannelFromMentionTypesEverywhere(field, discordgo.ChannelTypeGuildText)
		if err != nil {
			continue
		}

		break
	}
	if channel == nil {
		event.Respond("common.invalid-params")
		return
	}

	has := permissions.NewDiscordPermission("", permissionID).Match(
		p.state,
		event.DB(),
		user.ID,
		channel.ID,
		event.DM(),
	)

	_, err = event.Respond("dev.permission",
		"has", has, "permissionID", permissionID, "userID", user.ID, "channelID", channel.ID)
	if err != nil {
		event.Except(err)
		return
	}
}
