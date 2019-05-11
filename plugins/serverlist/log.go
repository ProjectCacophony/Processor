package serverlist

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/config"
	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleLog(event *events.Event) {
	err := config.GuildSetString(
		p.db, event.GuildID, logChannelIDKey, event.ChannelID,
	)
	if err != nil {
		event.Except(err)
		return
	}

	_, err = event.Respond("serverlist.log.success", "channelID", event.ChannelID)
	event.Except(err)
}

func (p *Plugin) sendLogMessageForServer(
	session *discord.Session,
	server *Server,
	send *discordgo.MessageSend,
) error {
	logChannelIDsPosted := make(map[string]interface{})

	for _, category := range server.Categories {
		logChannelID, err := config.GuildGetString(
			p.db, category.Category.GuildID, logChannelIDKey,
		)
		if err != nil {
			if strings.Contains(err.Error(), "record not found") {
				continue
			}
			return err
		}

		if _, ok := logChannelIDsPosted[logChannelID]; ok {
			continue
		}

		logChannelIDsPosted[logChannelID] = nil

		discord.SendComplexWithVars(
			session,
			p.Localizations(),
			logChannelID,
			send,
			"server",
			server,
		)
	}

	return nil
}
