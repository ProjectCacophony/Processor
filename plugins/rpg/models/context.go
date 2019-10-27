package models

import (
	"github.com/bwmarrin/discordgo"
)

type Context struct {
	GuildID string
	UserID  string
	Message *discordgo.Message
}
