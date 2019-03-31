package serverlist

import (
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"

	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/regexp"

	"gitlab.com/Cacophony/go-kit/events"
)

func extractExistingServerFromArg(redis *redis.Client, db *gorm.DB, session *discord.Session, arg string) *Server {
	server, err := serverFind(db, "invite_code = ? OR guild_id = ?", arg, arg)
	if err == nil {
		return server
	}

	parts := regexp.DiscordInviteRegexp.FindStringSubmatch(arg)
	if len(parts) >= 6 {
		invite, err := discord.Invite(redis, session, parts[5])
		if err != nil {
			return nil
		}

		server, err := serverFind(db, "guild_id = ?", invite.Guild.ID)
		if err == nil {
			return server
		}
	}

	return nil
}

func (p *Plugin) handleRemove(event *events.Event) {
	if len(event.Fields()) < 3 {
		event.Respond("serverlist.remove.too-few-args") // nolint: errcheck
		return
	}

	server := extractExistingServerFromArg(p.redis, p.db, event.Discord(), event.Fields()[2])
	if server == nil {
		event.Respond("serverlist.remove.no-server") // nolint: errcheck
		return
	}

	if !stringSliceContains(event.UserID, server.EditorUserIDs) {
		event.Respond("serverlist.remove.no-editor") // nolint: errcheck
		return
	}

	err := server.Remove(p)
	if err != nil {
		event.Except(err)
		return
	}

	_, err = event.Respond("serverlist.remove.success", "entry", server)
	event.Except(err)
}
