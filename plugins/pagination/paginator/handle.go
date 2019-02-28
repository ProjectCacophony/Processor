package paginator

import (
	"strconv"

	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Paginator) Handle(event *events.Event) {
	switch event.Type {
	case events.MessageCreateType:
		p.handleMessageCreate(event)
	case events.MessageReactionAddType:
		p.handleMessageReactionAdd(event)
	}

}

func (p *Paginator) handleMessageReactionAdd(event *events.Event) {
	if !validReactions[event.MessageReactionAdd.Emoji.Name] {
		return
	}

	pagedMessage, err := getPagedMessage(
		p.redis, event.MessageReactionAdd.MessageID,
	)
	if err != nil {
		return
	}

	if event.MessageReactionAdd.UserID != pagedMessage.UserID {
		return
	}

	err = p.handleReaction(
		pagedMessage,
		event.MessageReactionAdd,
	)
	if err != nil {
		event.ExceptSilent(err)
	}
}

func (p *Paginator) handleMessageCreate(event *events.Event) {
	if !p.messageRegexp.MatchString(event.MessageCreate.Content) {
		return
	}

	page, err := strconv.Atoi(event.MessageCreate.Content)
	if err != nil {
		return
	}

	if !isNumbersListening(
		p.redis, event.MessageCreate.ChannelID, event.MessageCreate.Author.ID,
	) {
		return
	}

	listener, err := getNumbersListeningMessageDelete(
		p.redis, event.MessageCreate.ChannelID, event.MessageCreate.Author.ID,
	)
	if err != nil {
		return
	}

	message, err := getPagedMessage(p.redis, listener.PagedEmbedMessageID)
	if err != nil {
		return
	}

	err = p.setPage(message, page)
	if err != nil {
		event.ExceptSilent(err)
		return
	}

	// clean up
	err = event.Discord().ChannelMessageDelete(event.ChannelID, listener.MessageID)
	if err != nil {
		return
	}
	event.Discord().ChannelMessageDelete(event.ChannelID, event.MessageCreate.ID) // nolint: errcheck
}
