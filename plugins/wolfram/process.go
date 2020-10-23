package wolfram

import (
	"net/url"
	"strings"

	"github.com/Krognol/go-wolfram"
	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) askWolfram(event *events.Event) {

	if len(event.Fields()) == 1 {
		event.Respond("wolfram.no-question")
		return
	}

	question := strings.Join(event.Fields()[1:], " ")
	isImageQuestion := false
	if event.Fields()[1] == "image" || event.Fields()[1] == "img" {
		isImageQuestion = true

		if len(event.Fields()) < 3 {
			event.Respond("wolfram.no-question")
			return
		}

		question = strings.Join(event.Fields()[2:], " ")
	}

	if !isImageQuestion {
		res, err := p.wolframClient.GetShortAnswerQuery(question, wolfram.Metric, 10)
		if err != nil {
			event.Except(err)
			return
		}

		if res == "" {
			event.Respond("wolfram.no-response")
			return
		}

		event.Respond(res)
	} else {
		urlValues := url.Values{}
		urlValues.Add("layout", "labelbar")
		urlValues.Add("timeout", "30")

		image, _, err := p.wolframClient.GetSimpleQuery(question, urlValues)
		if err != nil {
			event.Except(err)
			return
		}
		event.RespondComplex(&discordgo.MessageSend{
			Files: []*discordgo.File{{
				Name:   "wolframalpha-caco.png",
				Reader: image,
			}},
		})
	}
}
