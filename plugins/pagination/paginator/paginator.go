package paginator

import (
	"bytes"
	"fmt"
	"io"
	"regexp"

	"go.uber.org/zap"

	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/state"

	"github.com/go-redis/redis"

	"github.com/bwmarrin/discordgo"
)

const (
	LeftArrowEmoji  = "⬅"
	RightArrowEmoji = "➡"
	CloseEmoji      = "🇽"
	NumbersEmoji    = "🔢"

	FieldMessageType = iota
	ImageMessageType
)

type Paginator struct {
	logger *zap.Logger
	redis  *redis.Client
	state  *state.State
	tokens map[string]string

	messageRegexp *regexp.Regexp
}

func NewPaginator(
	logger *zap.Logger,
	redis *redis.Client,
	state *state.State,
	tokens map[string]string,
) (*Paginator, error) {
	p := &Paginator{
		logger: logger,
		redis:  redis,
		state:  state,
		tokens: tokens,
	}

	var err error
	p.messageRegexp, err = regexp.Compile("^[0-9]+$") // nolint: gocritic
	return p, err
}

// nolint: gochecknoglobals
var (
	validReactions = map[string]bool{
		LeftArrowEmoji:  true,
		RightArrowEmoji: true,
		CloseEmoji:      true,
		NumbersEmoji:    true,
	}
)

func (p *Paginator) getSession(guildID string) (*discordgo.Session, error) {
	botID, err := p.state.BotForGuild(guildID)
	if err != nil {
		return nil, err
	}

	return discord.NewSession(p.tokens, botID)
}

func (p *Paginator) sendComplex(
	guildID, channelID string, send *discordgo.MessageSend,
) ([]*discordgo.Message, error) {
	session, err := p.getSession(guildID)
	if err != nil {
		return nil, err
	}

	return discord.SendComplexWithVars(
		session,
		nil,
		channelID,
		send,
	)
}

func (p *Paginator) editComplex(
	guildID string, edit *discordgo.MessageEdit) (*discordgo.Message, error) {
	session, err := p.getSession(guildID)
	if err != nil {
		return nil, err
	}

	return session.ChannelMessageEditComplex(edit)
}

// setupAndSendFirstMessage
func (p *Paginator) setupAndSendFirstMessage(message *PagedEmbedMessage) {
	var sentMessage []*discordgo.Message
	var err error

	// copy the embedded message so changes can be made to it
	tempEmbed := &discordgo.MessageEmbed{}
	*tempEmbed = *message.FullEmbed

	// set footer which will hold information about the page it is on
	tempEmbed.Footer = p.getEmbedFooter(message)

	if message.MsgType == ImageMessageType {

		// if fields were sent with image embed, handle those
		if len(message.FullEmbed.Fields) > 0 {

			// get start and end fields based on current page and fields per page
			startField := (message.CurrentPage - 1) * message.FieldsPerPage
			endField := startField + message.FieldsPerPage
			if endField > len(message.FullEmbed.Fields) {
				endField = len(message.FullEmbed.Fields)
			}

			tempEmbed.Fields = tempEmbed.Fields[startField:endField]
		}

		var buf bytes.Buffer
		newReader := io.TeeReader(message.Files[message.CurrentPage-1].Reader, &buf)
		message.Files[message.CurrentPage-1].Reader = &buf

		tempEmbed.Image.URL = fmt.Sprintf("attachment://%s", message.Files[message.CurrentPage-1].Name)
		sentMessage, err = p.sendComplex(
			message.GuildID, message.ChannelID, &discordgo.MessageSend{
				Embed: tempEmbed,
				Files: []*discordgo.File{{
					Name:   message.Files[message.CurrentPage-1].Name,
					Reader: newReader,
				}},
			})
		if p.hasError(message, err) {
			return
		}

	} else {
		// reduce fields to the fields per page
		tempEmbed.Fields = tempEmbed.Fields[:message.FieldsPerPage]

		sentMessage, err = p.sendComplex(
			message.GuildID, message.ChannelID, &discordgo.MessageSend{
				Embed: tempEmbed,
			})
		if p.hasError(message, err) {
			return
		}
	}

	message.MessageID = sentMessage[0].ID
	p.addReactionsToMessage(message)
}

// getEmbedFooter is a simlple helper function to return the footer for the embed message
func (p *Paginator) getEmbedFooter(message *PagedEmbedMessage) *discordgo.MessageEmbedFooter {
	var footerText string

	// check if embed had a footer, if so attach to page count
	if message.FullEmbed.Footer != nil && message.FullEmbed.Footer.Text != "" {
		footerText = fmt.Sprintf(
			"Page: %d / %d | %s",
			message.CurrentPage, message.TotalNumOfPages, message.FullEmbed.Footer.Text,
		)
	} else {
		footerText = fmt.Sprintf(
			"Page: %d / %d",
			message.CurrentPage, message.TotalNumOfPages,
		)
	}

	return &discordgo.MessageEmbedFooter{Text: footerText}
}

func (p *Paginator) addReactionsToMessage(message *PagedEmbedMessage) {
	session, err := p.getSession(message.GuildID)
	if err != nil {
		return
	}

	err = session.MessageReactionAdd(message.ChannelID, message.MessageID, LeftArrowEmoji)
	if err != nil {
		return
	}
	err = session.MessageReactionAdd(message.ChannelID, message.MessageID, RightArrowEmoji)
	if err != nil {
		return
	}

	if message.TotalNumOfPages > 2 {
		err = session.MessageReactionAdd(message.ChannelID, message.MessageID, NumbersEmoji)
		if err != nil {
			return
		}
	}

	session.MessageReactionAdd(message.ChannelID, message.MessageID, CloseEmoji) // nolint: errcheck
}

func (p *Paginator) hasError(message *PagedEmbedMessage, err error) bool {
	if err == nil {
		return false
	}

	// delete from current embeds
	deletePagedMessage(p.redis, message.MessageID) // nolint: errcheck

	// check if error is a permissions error
	if err, ok := err.(*discordgo.RESTError); ok && err.Message.Code == discordgo.ErrCodeMissingPermissions {
		if message.MsgType == ImageMessageType {
			p.sendComplex( // nolint: errcheck
				message.GuildID, message.ChannelID, &discordgo.MessageSend{
					Content: "bot.errors.no-embed-or-file", // TODO
				})
		} else {
			p.sendComplex( // nolint: errcheck
				message.GuildID, message.ChannelID, &discordgo.MessageSend{
					Content: "bot.errors.no-embed", // TODO
				})
		}
	} else {
		p.logger.Error("unexpected error",
			zap.Error(err),
		)
		// TODO: capture in raven
	}

	return true
}
