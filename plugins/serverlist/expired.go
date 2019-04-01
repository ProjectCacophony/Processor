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

	p.refreshList(server.BotID) // nolint: errcheck

	session, err := discord.NewSession(p.tokens, server.BotID)
	if err != nil {
		event.ExceptSilent(err)
		return
	}

	for _, editorUserID := range server.EditorUserIDs {
		discord.SendComplexWithVars( // nolint: errcheck
			p.redis,
			session,
			p.Localisations(),
			editorUserID,
			&discordgo.MessageSend{
				Content: "serverlist.dm.server-expired",
			},
			true,
			"server",
			server,
		)
	}
}
