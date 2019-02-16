package dev

import (
	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/events"
)

func handleDevTranslate(event *events.Event) {
	_, err := event.RespondComplex(&discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Description: "dev.translate.embed.description",
		},
	})
	event.Except(err)
}
