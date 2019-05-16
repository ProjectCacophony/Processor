package serverlist

import (
	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleExpired(event *events.Event) {
	server, err := serverFind(p.db, "id = ?", event.ServerlistServerExpire.ID)
	if err != nil {
		event.ExceptSilent(err)
		return
	}

	p.refreshList(server.BotID)

	session, err := discord.NewSession(p.tokens, server.BotID)
	if err != nil {
		event.ExceptSilent(err)
		return
	}

	for _, editorUserID := range server.EditorUserIDs {
		channelID, err := discord.DMChannel(p.redis, session, editorUserID)
		if err != nil {
			event.ExceptSilent(err)
			continue
		}

		discord.SendComplexWithVars(
			session,
			event.Localizations(),
			channelID,
			&discordgo.MessageSend{
				Content: "serverlist.dm.server-expired",
			},
			"server",
			server,
		)
	}
}
