package serverlist

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/events"
)

type enhancedServer struct {
	*Server
	Categories []string
	Editors    []*discordgo.User
}

func (p *Plugin) handleStatus(event *events.Event) {
	servers, err := serversFindMany(
		p.db, "bot_id = ? AND editor_user_ids @> ARRAY[?]::varchar[]",
		event.BotUserID, event.UserID,
	)
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		event.Except(err)
		return
	}

	enhancedServers := make([]enhancedServer, len(servers))
	for i, server := range servers {
		var categories []string
		for _, category := range server.Categories {
			if len(category.Category.Keywords) == 0 {
				continue
			}

			categories = append(categories, category.Category.Keywords[0])
		}
		var editors []*discordgo.User
		for _, editor := range server.EditorUserIDs {
			user, err := p.state.User(editor)
			if err != nil {
				user = &discordgo.User{
					ID: editor,
				}
			}

			editors = append(editors, user)
		}

		enhancedServers[i] = enhancedServer{
			Server:     server,
			Categories: categories,
			Editors:    editors,
		}
	}

	_, err = event.Respond("serverlist.status.content", "entries", enhancedServers)
	event.Except(err)
}
