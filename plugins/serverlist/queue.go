package serverlist

import (
	"strings"
	"time"

	"gitlab.com/Cacophony/go-kit/permissions"

	lock "github.com/bsm/redis-lock"

	"gitlab.com/Cacophony/go-kit/discord"

	"github.com/bwmarrin/discordgo"

	"github.com/pkg/errors"
	"gitlab.com/Cacophony/go-kit/config"
	"gitlab.com/Cacophony/go-kit/events"
	"go.uber.org/zap"
)

func (p *Plugin) handleQueue(event *events.Event) {
	err := config.GuildSetString(
		p.db, event.GuildID, queueChannelIDKey, event.ChannelID,
	)
	if err != nil {
		event.Except(err)
		return
	}

	_, err = event.Respond("serverlist.queue.content", "channelID", event.ChannelID)
	event.Except(err)

	p.refreshQueue(event.GuildID)
}

func (p *Plugin) handleQueueReaction(event *events.Event) bool {
	if event.MessageReactionAdd.Emoji.Name != emojiApprove {
		return false
	}

	channelID, err := config.GuildGetString(
		p.db, event.GuildID, queueChannelIDKey,
	)
	if err != nil {
		if !strings.Contains(err.Error(), "record not found") {
			event.ExceptSilent(err)
		}
		return false
	}

	if event.ChannelID != channelID {
		return false
	}

	if !event.Has(permissions.BotOwner) {
		return false
	}

	err = p.approveCurrentServer(event.BotUserID, event.GuildID)
	if err != nil &&
		!strings.Contains(err.Error(), "nothing to approve") &&
		!strings.Contains(err.Error(), "nothing to reject") {
		event.ExceptSilent(err)
		return true
	}

	p.refreshQueue(event.GuildID)

	discord.RemoveReact( // nolint: errcheck
		p.redis,
		event.Discord(),
		event.ChannelID,
		event.MessageReactionAdd.MessageID,
		event.UserID,
		false,
		event.MessageReactionAdd.Emoji.Name,
	)

	return true
}

func (p *Plugin) handleQueueReject(event *events.Event) {
	if len(event.Fields()) < 3 {
		event.Respond("serverlist.queue-reject.no-reason") // nolint: errcheck
		return
	}

	err := p.rejectCurrentServer(event.BotUserID, event.GuildID, event.Fields()[2])
	if err != nil {
		event.Except(err)
		return
	}

	_, err = event.Respond("serverlist.queue-reject.success")
	event.Except(err)
}

func (p *Plugin) handleQueueRefresh(event *events.Event) {
	err := p.refreshQueueForGuild(event.GuildID)
	if err != nil {
		event.Except(err)
		return
	}

	_, err = event.Respond("serverlist.queue-refresh.content")
	event.Except(err)
}

func (p *Plugin) refreshQueue(guildIDs ...string) {
	var err error

	for _, guildID := range guildIDs {
		err = p.refreshQueueForGuild(guildID)
		if err != nil {
			p.logger.Error("failure refreshing queue for guild",
				zap.Error(err),
				zap.String("guild_id", guildID),
			)
		}
	}
}

func (p *Plugin) approveCurrentServer(botID, guildID string) error {
	queueMessage, err := p.getCurrentQueueMessage(guildID)
	if err != nil {
		return errors.Wrap(err, "error getting QueueMessage from config")
	}

	queue, err := p.getQueue(botID)
	if err != nil {
		return errors.Wrap(err, "unable to query for queued entries")
	}

	server := queueFind(queueMessage.CurrentServerID, queue)
	if server == nil {
		return errors.New("found nothing to approve")
	}

	return server.QueueApprove(p)
}

func (p *Plugin) rejectCurrentServer(botID, guildID, reason string) error {
	queueMessage, err := p.getCurrentQueueMessage(guildID)
	if err != nil {
		return errors.Wrap(err, "error getting QueueMessage from config")
	}

	queue, err := p.getQueue(botID)
	if err != nil {
		return errors.Wrap(err, "unable to query for queued entries")
	}

	server := queueFind(queueMessage.CurrentServerID, queue)
	if server == nil {
		return errors.New("found nothing to reject")
	}

	return server.QueueReject(p, reason)
}

func (p *Plugin) refreshQueueForGuild(guildID string) error {
	guildLock := p.getGuildLock(guildID)

	locked, err := guildLock.Lock()
	if err != nil {
		return errors.Wrap(err, "error acquiring lock")
	}
	if !locked {
		return errors.Wrap(err, "unable to acquire lock")
	}
	defer guildLock.Unlock() // nolint: errcheck

	channelID, err := config.GuildGetString(
		p.db, guildID, queueChannelIDKey,
	)
	if err != nil {
		return errors.Wrap(err, "failure getting queueChannelIDKey from guild config")
	}

	if channelID == "" {
		return nil
	}

	botID, err := p.state.BotForGuild(guildID)
	if err != nil {
		return errors.Wrap(err, "failure getting Bot ID for Guild")
	}

	session, err := discord.NewSession(p.tokens, botID)
	if err != nil {
		return errors.Wrap(err, "failure creating Discord Session for Bot")
	}

	queue, err := p.getQueue(botID)
	if err != nil {
		return errors.Wrap(err, "unable to query for queued entries")
	}

	currentQueueMessage, err := p.getCurrentQueueMessage(guildID)
	if err != nil {
		return errors.Wrap(err, "error getting QueueMessage from config")
	}

	if currentQueueMessage == nil {
		if len(queue) == 0 {
			return nil
		}

		currentQueueMessage = &QueueMessage{}
	}

	previousServerID := currentQueueMessage.CurrentServerID

	var queueItem *Server
	if len(queue) > 0 {
		queueItem = queue[0]
	}

	// create new queue message if none exists
	if currentQueueMessage == nil {
		return p.newQueueMessage(
			session,
			guildID,
			channelID,
			queueItem,
			queue,
		)
	}

	activeItem := queueFind(currentQueueMessage.CurrentServerID, queue)
	if activeItem != nil {
		queueItem = activeItem
	}

	// update queue message with new server
	err = p.updateQueueMessage(
		session,
		guildID,
		channelID,
		queueItem,
		queue,
		currentQueueMessage,
	)
	if err != nil {
		return err
	}

	// ping if queue was empty before, but is not anymore
	if previousServerID == 0 && len(queue) > 0 {
		pingMessages, err := discord.SendComplexWithVars(
			p.redis,
			session,
			p.Localisations(),
			channelID,
			&discordgo.MessageSend{
				Content: ":nayoung:",
			},
			false,
		)
		if err == nil {
			discord.Delete( // nolint: errcheck
				p.redis,
				session,
				pingMessages[0].ChannelID,
				pingMessages[0].ID,
				false,
			)
		}
	}

	return nil
}

func (p *Plugin) newQueueMessage(
	session *discord.Session,
	guildID string,
	channelID string,
	queueItem *Server,
	queue []*Server,
) error {
	embed := p.getQueueMessageEmbed(queueItem, len(queue))

	messages, err := discord.SendComplexWithVars(
		p.redis,
		session,
		p.Localisations(),
		channelID,
		&discordgo.MessageSend{
			Embed: embed,
		},
		false,
		"server",
		queueItem,
	)
	if err != nil {
		return errors.Wrap(err, "error sending new QueueMessage")
	}

	var currentServerID uint
	if queueItem != nil {
		currentServerID = queueItem.ID
	}

	currentQueueMessage := &QueueMessage{
		CurrentServerID: currentServerID,
		MessageID:       messages[0].ID,
		Embed:           embed,
	}

	err = config.GuildSetInterface(p.db, guildID, queueMessageKey, currentQueueMessage)
	if err != nil {
		return errors.Wrap(err, "error saving new QueueMessage")
	}

	discord.React( // nolint: errcheck
		p.redis,
		session,
		channelID,
		messages[0].ID,
		false,
		emojiApprove,
	)

	return nil
}

func (p *Plugin) updateQueueMessage(
	session *discord.Session,
	guildID string,
	channelID string,
	queueItem *Server,
	queue []*Server,
	currentMessage *QueueMessage,
) error {
	embed := p.getQueueMessageEmbed(queueItem, len(queue))

	_, err := discord.EditComplexWithVars(
		p.redis,
		session,
		p.Localisations(),
		&discordgo.MessageEdit{
			Embed:   embed,
			ID:      currentMessage.MessageID,
			Channel: channelID,
		},
		false,
		"server",
		queueItem,
	)
	if err != nil {
		if errD, ok := err.(*discordgo.RESTError); ok &&
			errD.Message != nil &&
			errD.Message.Code == discordgo.ErrCodeUnknownMessage {
			return p.newQueueMessage(session, guildID, channelID, queueItem, queue)
		}

		return errors.Wrap(err, "error editing existing QueueMessage")
	}

	currentMessage.CurrentServerID = 0
	if queueItem != nil {
		currentMessage.CurrentServerID = queueItem.ID
	}

	currentMessage.Embed = embed

	err = config.GuildSetInterface(p.db, guildID, queueMessageKey, currentMessage)
	if err != nil {
		return errors.Wrap(err, "error saving updated existing QueueMessage")
	}

	return nil
}

func queueFind(n uint, list []*Server) *Server {
	for _, item := range list {
		if item.ID != n {
			continue
		}

		return item
	}

	return nil
}

func (p *Plugin) getGuildLock(guildID string) *lock.Locker {
	return lock.New(
		p.redis,
		refreshQueueLock(guildID),
		&lock.Options{
			LockTimeout: 5 * time.Minute,
			RetryCount:  9, // try for 1 1/2 minutes
			RetryDelay:  10 * time.Second,
		},
	)
}

func (p *Plugin) getCurrentQueueMessage(guildID string) (*QueueMessage, error) {
	var queueMessage *QueueMessage
	err := config.GuildGetInterface(p.db, guildID, queueMessageKey, &queueMessage)
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		return nil, err
	}

	return queueMessage, nil
}

func (p *Plugin) getQueue(botID string) ([]*Server, error) {
	return serversFindMany(
		p.db,
		"state = ? AND bot_id = ?", StateQueued, botID,
	)
}
