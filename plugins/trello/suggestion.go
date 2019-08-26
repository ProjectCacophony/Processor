package trello

import (
	"strings"
	"time"

	trello "github.com/VojtechVitek/go-trello"
	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleSuggestion(event *events.Event) {
	if len(event.Fields()) < 2 {
		event.Respond("common.invalid-params")
		return
	}

	issueTitle := event.Fields()[1]

	// if there are more than 3 fields, combine all of them to make the title to avoid require quotes
	if len(event.Fields()) > 3 {
		issueTitle = strings.Join(event.Fields()[1:], " ")
	}

	if len(issueTitle) > 50 {
		event.Respond("trello.title.to-long")
		return
	}

	issueDescription := ""
	if len(event.Fields()) == 3 {
		issueDescription = event.Fields()[2]
	}

	list, err := p.trello.List(backlogBoardID)
	if err != nil {
		event.Except(err)
		return
	}

	issueTitle = strings.TrimSpace(issueTitle)
	issueDescription = strings.TrimSpace(issueDescription)

	_, err = list.AddCard(trello.Card{
		Name: issueTitle,
		Desc: issueDescription,
	})
	if err != nil {
		event.Except(err)
		return
	}

	logRequest(event, issueTitle, issueDescription)
	event.Respond("trello.suggestion.received")
}

func logRequest(event *events.Event, title, description string) {
	if trelloLogChannelID == "" {
		return
	}

	user, err := event.State().User(event.UserID)
	if err != nil {
		event.Except(err)
		return
	}

	// trello blue from their logo image
	color := discord.HexToColorCode("0079c1")

	event.SendComplex(trelloLogChannelID, &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title:       title,
			Description: description,
			Color:       color,
			Timestamp:   time.Now().Format(time.RFC3339),
			Author: &discordgo.MessageEmbedAuthor{
				Name:    user.Username + "#" + user.Discriminator + " (#" + user.ID + ")",
				IconURL: user.AvatarURL(""),
			},
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Guild ",
					Value:  event.GuildID,
					Inline: true,
				},
				{
					Name:   "Channel",
					Value:  event.ChannelID,
					Inline: true,
				},
			},
		},
	})
}
