package tools

import (
	"math/rand"

	"gitlab.com/Cacophony/go-kit/events"
)

type eightballChoiceType int

const (
	eightballChoiceTypePositive eightballChoiceType = iota
	eightballChoiceTypeNeutral
	eightballChoiceTypeNegative
)

var (
	eightballChoices = map[string]eightballChoiceType{
		"It is certain.":             eightballChoiceTypePositive,
		"It is decidedly so.":        eightballChoiceTypePositive,
		"Without a doubt.":           eightballChoiceTypePositive,
		"Yes — definitely.":          eightballChoiceTypePositive,
		"You may rely on it.":        eightballChoiceTypePositive,
		"As I see it, yes.":          eightballChoiceTypePositive,
		"Most likely.":               eightballChoiceTypePositive,
		"Outlook good.":              eightballChoiceTypePositive,
		"Yes.":                       eightballChoiceTypePositive,
		"Signs point to yes.":        eightballChoiceTypePositive,
		"Reply hazy, try again.":     eightballChoiceTypeNeutral,
		"Ask again later.":           eightballChoiceTypeNeutral,
		"Better not tell you now.":   eightballChoiceTypeNeutral,
		"Cannot predict now.":        eightballChoiceTypeNeutral,
		"Concentrate and ask again.": eightballChoiceTypeNeutral,
		"Don't count on it.":         eightballChoiceTypeNegative,
		"My reply is no.":            eightballChoiceTypeNegative,
		"My sources say no.":         eightballChoiceTypeNegative,
		"Outlook not so good.":       eightballChoiceTypeNegative,
		"Very doubtful.":             eightballChoiceTypeNegative,
	}
	eightballChoicesLength = len(eightballChoices)
)

func (p *Plugin) handle8ball(event *events.Event) {
	pick, choiceType := randomEightballChoice()

	_, err := event.Respond("tools.8ball.result", "pick", pick, "type", choiceType)
	event.Except(err)
}

func randomEightballChoice() (string, eightballChoiceType) {
	i := rand.Intn(eightballChoicesLength)

	for key, value := range eightballChoices {
		if i > 0 {
			i--
			continue
		}

		return key, value
	}

	return "", 0
}
