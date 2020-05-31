package admin

import (
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/events"
)

var (
	interceptionMap     = map[string]*interceptionDetails{}
	interceptionMapLock sync.RWMutex
)

type interceptionDetails struct {
	Expire      time.Time
	ToChannelID string
}

func interceptionMapAdd(botUserID, fromChannelID, toChannelID string) {
	interceptionMapLock.Lock()
	defer interceptionMapLock.Unlock()

	interceptionMap[botUserID+"-"+fromChannelID] = &interceptionDetails{
		Expire:      time.Now().Add(5 * time.Minute),
		ToChannelID: toChannelID,
	}
}

func interceptionMapRead(botUserID, fromChannelID string) string {
	interceptionMapLock.RLock()

	if interceptionMap[botUserID+"-"+fromChannelID] == nil {
		interceptionMapLock.RUnlock()

		return ""
	}

	interceptionDetails := interceptionMap[botUserID+"-"+fromChannelID]
	interceptionMapLock.RUnlock()

	if time.Now().After(interceptionDetails.Expire) {
		interceptionMapLock.Lock()
		defer interceptionMapLock.Unlock()

		interceptionMap[botUserID+"-"+fromChannelID] = nil
		return ""
	}

	return interceptionDetails.ToChannelID
}

func (p *Plugin) handleIntercept(event *events.Event) {
	fromChannel, err := p.state.ChannelFromMention(event.GuildID, event.Fields()[2])
	if err != nil {
		event.Except(err)
		return
	}

	botForChannel, err := event.State().BotForChannel(fromChannel.ID)
	if err != nil {
		event.Except(err)
		return
	}

	fromGuild, err := event.State().Guild(event.MessageCreate.GuildID)
	if err != nil {
		event.ExceptSilent(err)
		return
	}

	interceptionMapAdd(botForChannel, fromChannel.ID, event.ChannelID)

	event.Respond(
		"admin.intercept.start-local",
		"fromChannel", fromChannel,
		"fromGuild", fromGuild,
	)
	event.Send(
		fromChannel.ID,
		"admin.intercept.start-remote",
		"fromChannel", fromChannel,
		"fromGuild", fromGuild,
	)
}

func (p *Plugin) copyMessageCreate(event *events.Event, toChannelID string) {
	fromChannel, err := event.State().Channel(event.MessageCreate.ChannelID)
	if err != nil {
		event.ExceptSilent(err)
		return
	}

	fromGuild, err := event.State().Guild(event.MessageCreate.GuildID)
	if err != nil {
		event.ExceptSilent(err)
		return
	}

	send := &discordgo.MessageSend{}

	send.Content = event.Translate(
		"admin.intercept.copy-from-disclaimer",
		"fromChannel", fromChannel,
		"fromGuild", fromGuild,
		"messageCreate", event.MessageCreate,
	)
	if len(event.MessageCreate.Embeds) > 0 {
		send.Embed = event.MessageCreate.Embeds[0]
	}

	event.SendComplex(toChannelID, send)

	p.logger.Info("need to copy event")
}
