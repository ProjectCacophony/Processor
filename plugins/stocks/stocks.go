package stocks

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/events"
)

const (
	colorRed   = "#FF3B2F"
	colorGreen = "#35C759"
)

func (p *Plugin) handleStocks(event *events.Event) {
	if len(event.Fields()) <= 1 {
		event.Respond("stocks.no-symbol")
		return
	}

	symbol := event.Fields()[1]

	symbolData, err := p.lookupSymbol(event.Context(), symbol)
	if err != nil {
		if strings.Contains(err.Error(), "symbol not found") {
			event.Respond("stocks.symbol-not-found")
			return
		}
		event.Except(err)
		return
	}

	quote, err := p.iexClient.StocksQuote(event.Context(), symbol)
	if err != nil {
		event.Except(errors.Wrap(err, "failure looking up stocks quote"))
		return
	}

	logo, err := p.iexClient.StocksLogo(event.Context(), symbol)
	if err != nil {
		event.Except(errors.Wrap(err, "failure looking up stocks logo"))
		return
	}

	color := discord.HexToColorCode(colorGreen)
	if quote.Change < 0 {
		color = discord.HexToColorCode(colorRed)
	}

	var thumbnail *discordgo.MessageEmbedThumbnail
	if logo.URL != "" {
		thumbnail = &discordgo.MessageEmbedThumbnail{
			URL: logo.URL,
		}
	}

	_, err = event.RespondComplex(&discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title:       "stocks.embed.title",
			Description: "stock.embed.description",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "stock.embed.field.open.name",
					Value:  "stock.embed.field.open.value",
					Inline: false,
				},
				{
					Name:   "stock.embed.field.high.name",
					Value:  "stock.embed.field.high.value",
					Inline: true,
				},
				{
					Name:   "stock.embed.field.low.name",
					Value:  "stock.embed.field.low.value",
					Inline: true,
				},
				{
					Name:   "stock.embed.field.52whigh.name",
					Value:  "stock.embed.field.52whigh.value",
					Inline: true,
				},
				{
					Name:   "stock.embed.field.52wlow.name",
					Value:  "stock.embed.field.52wlow.value",
					Inline: true,
				},
			},
			Color:     color,
			Thumbnail: thumbnail,
		},
	},
		"symbol", symbolData,
		"quote", quote,
	)
	event.Except(err)
}
