package actions

import (
	"strings"

	"gitlab.com/Cacophony/Processor/plugins/automod/models"
)

func automodReason(env *models.Env, action string) string {
	botUserText := "Cacophony"
	botUser, err := env.State.User(env.Event.BotUserID)
	if err == nil {
		botUserText = botUser.String()
	}
	ruleNameText := "Unknown"
	if env.Rule != nil {
		ruleNameText = env.Rule.Name
	}

	return strings.Title(action) + " by " + botUserText + " Automod Rule: " + ruleNameText
}
