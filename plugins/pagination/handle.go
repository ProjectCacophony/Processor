package pagination

import (
	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/Processor/plugins/pagination/paginator"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/interfaces"
)

type Plugin struct {
	paginator *paginator.Paginator
}

func (p *Plugin) Name() string {
	return "pagination"
}

func (p *Plugin) Start(params common.StartParameters) error {
	var err error
	p.paginator, err = paginator.NewPaginator(
		params.Logger,
		params.Redis,
		params.State,
		params.Tokens,
	)

	return err
}

func (p *Plugin) Stop(params common.StopParameters) error {
	return nil
}

func (p *Plugin) Priority() int {
	return 100000
}

func (p *Plugin) Passthrough() bool {
	return true
}

func (p *Plugin) Localisations() []interfaces.Localisation {
	return nil
}

func (p *Plugin) Action(event *events.Event) bool {
	p.paginator.Handle(event)

	if event.Type != events.MessageCreateType {
		return false
	}

	switch event.MessageCreate.Content {

	case "embed pls":
		p.paginator.FieldsPaginator( // nolint: errcheck
			event.MessageCreate.GuildID,
			event.MessageCreate.ChannelID,
			event.MessageCreate.Author.ID,
			&discordgo.MessageEmbed{
				Description: "test",
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:  "a",
						Value: "a",
					},
					{
						Name:  "b",
						Value: "b",
					},
					{
						Name:  "c",
						Value: "c",
					},
					{
						Name:  "d",
						Value: "d",
					},
					{
						Name:  "e",
						Value: "e",
					},
					{
						Name:  "f",
						Value: "f",
					},
					{
						Name:  "g",
						Value: "g",
					},
					{
						Name:  "h",
						Value: "h",
					},
					{
						Name:  "i",
						Value: "i",
					},
					{
						Name:  "j",
						Value: "j",
					},
					{
						Name:  "k",
						Value: "k",
					},
				},
			},
			3,
		)

	case "embed embed pls":

		p.paginator.EmbedPaginator( // nolint: errcheck
			event.MessageCreate.GuildID,
			event.MessageCreate.ChannelID,
			event.MessageCreate.Author.ID,
			&discordgo.MessageEmbed{
				Description: "first page",
				Footer: &discordgo.MessageEmbedFooter{
					Text: "the footer of the first page",
				},
			},
			&discordgo.MessageEmbed{
				Description: "second page",
				Footer: &discordgo.MessageEmbedFooter{
					Text: "the footer of the second page",
				},
			},
			&discordgo.MessageEmbed{
				Description: "third page",
				Footer: &discordgo.MessageEmbedFooter{
					Text: "the footer of the third page",
				},
			},
			&discordgo.MessageEmbed{
				Description: "fourth page",
				Footer: &discordgo.MessageEmbedFooter{
					Text: "the footer of the fourth page",
				},
			},
			&discordgo.MessageEmbed{
				Description: "fifth page",
				Footer: &discordgo.MessageEmbedFooter{
					Text: "the footer of the fifth page",
				},
			},
		)
	}

	return false
}
