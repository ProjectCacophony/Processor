package paginator

import "github.com/bwmarrin/discordgo"

type PagedEmbedMessage struct {
	Files           []*discordgo.File
	FullEmbed       *discordgo.MessageEmbed
	TotalNumOfPages int
	CurrentPage     int
	FieldsPerPage   int
	Color           int
	MessageID       string
	GuildID         string
	ChannelID       string
	UserID          string // user who triggered the message
	MsgType         int
}

type numbersListener struct {
	MessageID           string // message ID of the message asking the user which page to choose
	PagedEmbedMessageID string // message ID of the paged embed message
}
