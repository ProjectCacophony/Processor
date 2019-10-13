package eventlog

import (
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/state"
)

var (
	discordColourGreen  = discord.HexToColorCode("#73d016")
	discordColourOrange = discord.HexToColorCode("#ffb80a")
	discordColourRed    = discord.HexToColorCode("#b22222") // nolint: deadcode, unused, varcheck // TODO: use it for destructive actions
)

type actionType string

type entityType string

type Item struct {
	gorm.Model
	UUID uuid.UUID `gorm:"UNIQUE_INDEX;NOT NULL;Type:uuid"`

	GuildID string `gorm:"NOT NULL"`

	ActionType actionType `gorm:"NOT NULL"`

	AuthorID string // Author UserID

	TargetType  entityType
	TargetValue string

	Reason string

	WaitingForAuditLogBackfill bool

	Options []ItemOption

	LogMessage ItemLogMessage `gorm:"embedded;embedded_prefix:log_message_"`
}

func (*Item) TableName() string {
	return "eventlog_items"
}

func (i *Item) Embed(state *state.State) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Title:       i.ActionType.String(),
		Description: "",
		Timestamp:   discord.Timestamp(i.CreatedAt),
		Color:       discordColourGreen,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Cacophony Eventlog #" + i.UUID.String(),
		},
		Fields:    make([]*discordgo.MessageEmbedField, 0, 1),
		Thumbnail: &discordgo.MessageEmbedThumbnail{},
	}

	if i.Reason != "" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  "Reason",
			Value: i.Reason,
		})
	}

	for _, option := range i.Options {
		var embedOptionValue string
		if option.PreviousValue != "" {
			embedOptionValue = option.Type.String(option.PreviousValue) + " ➡ "
		}
		if option.NewValue != "" {
			embedOptionValue += option.Type.String(option.NewValue)
		} else {
			embedOptionValue += "N/A"
		}

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  option.Key,
			Value: embedOptionValue,
		})
	}

	if i.AuthorID != "" {
		author, err := state.User(i.AuthorID)
		if err != nil {
			author = &discordgo.User{
				ID:       author.ID,
				Username: "N/A",
			}
		}

		embed.Author = &discordgo.MessageEmbedAuthor{
			Name:    "By " + author.String() + " #" + author.ID,
			IconURL: author.AvatarURL("128"),
		}
	}

	if i.TargetValue != "" {
		embed.Description += "On " + i.TargetType.String(i.TargetValue)

		switch i.TargetType {
		case EntityTypeGuild:
			guild, err := state.Guild(i.TargetValue)
			if err == nil && guild.IconURL() != "" {
				embed.Thumbnail.URL = guild.IconURL() + "?size=256"
			}
		case EntityTypeUser:
			user, err := state.User(i.TargetValue)
			if err == nil {
				embed.Thumbnail.URL = user.AvatarURL("256")
			}
		}
	}

	if i.WaitingForAuditLogBackfill {
		embed.Color = discordColourOrange
	}

	return embed
}

type ItemOption struct {
	gorm.Model
	ItemID uint `gorm:"NOT NULL"`

	Key           string `gorm:"NOT NULL"`
	PreviousValue string
	NewValue      string
	Type          entityType
}

func (*ItemOption) TableName() string {
	return "eventlog_item_options"
}

type ItemLogMessage struct {
	ChannelID string
	MessageID string
}
