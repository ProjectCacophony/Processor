package dev

import (
	"strconv"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/permissions"
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

	permissionID, err := strconv.ParseInt(event.Fields()[2], 10, 64)
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
		event.SuperUser,
	)

	_, err = event.Respond("dev.permission",
		"has", has, "permissionID", permissionID, "userID", user.ID, "channelID", channel.ID)
	if err != nil {
		event.Except(err)
		return
	}
}
