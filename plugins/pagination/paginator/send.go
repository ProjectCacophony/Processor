package paginator

import (
	"errors"
	"fmt"
	"math"

	"github.com/bwmarrin/discordgo"
)

// CreatePagedMessage creates the paged messages
func (p *Paginator) SendPagedMessage(
	guildID, channelID, userID string, embed *discordgo.MessageEmbed, fieldsPerPage int,
) error {

	// if there aren't multiple fields to be paged through,
	// or if the amount of fields is less than the requested fields per page
	// just send a normal embed
	if len(embed.Fields) < 2 || len(embed.Fields) <= fieldsPerPage {
		_, err := p.sendComplex(guildID, channelID, &discordgo.MessageSend{
			Embed: embed,
		})
		return err
	}

	// fields per page can not be less than 1
	if fieldsPerPage < 1 {
		return errors.New("FieldsPerPage can not be less than 1")
	}

	// create paged message
	pagedMessage := &PagedEmbedMessage{
		FullEmbed:       embed,
		ChannelID:       channelID,
		GuildID:         guildID,
		CurrentPage:     1,
		FieldsPerPage:   fieldsPerPage,
		TotalNumOfPages: int(math.Ceil(float64(len(embed.Fields)) / float64(fieldsPerPage))),
		UserID:          userID,
		MsgType:         FieldMessageType,
	}

	p.setupAndSendFirstMessage(pagedMessage)

	err := setPagedMessage(p.redis, pagedMessage.MessageID, pagedMessage)
	return err
}

// SendPagedImageMessage creates the paged image messages
func (p *Paginator) SendPagedImageMessage(
	guildID, channelID, userID string, msgSend *discordgo.MessageSend, fieldsPerPage int,
) error {
	if msgSend.Embed == nil {
		return errors.New("parameter msgSend must contain an embed")
	}

	// make sure the image url is set to the name of the first file incease it wasn't set
	msgSend.Embed.Image.URL = fmt.Sprintf("attachment://%s", msgSend.Files[0].Name)

	// check if there are multiple Files, not just send it normally
	if len(msgSend.Files) < 2 {
		_, err := p.sendComplex(guildID, channelID, msgSend)
		return err
	}

	// create paged message
	pagedMessage := &PagedEmbedMessage{
		FullEmbed:       msgSend.Embed,
		ChannelID:       channelID,
		GuildID:         guildID,
		CurrentPage:     1,
		FieldsPerPage:   fieldsPerPage,
		TotalNumOfPages: len(msgSend.Files),
		Files:           msgSend.Files,
		UserID:          userID,
		MsgType:         ImageMessageType,
	}

	p.setupAndSendFirstMessage(pagedMessage)

	err := setPagedMessage(p.redis, pagedMessage.MessageID, pagedMessage)
	return err
}
