package eventlog

import (
	"sort"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/state"
)

var (
	discordColourGreen  = discord.HexToColorCode("#73d016")
	discordColourOrange = discord.HexToColorCode("#ffb80a")
	discordColourRed    = discord.HexToColorCode("#b22222")
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

	WaitingForAuditLogBackfill bool
	Reverted                   bool

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

	// sort items by name
	sort.Slice(i.Options, func(j, k int) bool {
		return i.Options[j].Key < i.Options[k].Key
	})

	var optionAuthor *discordgo.User
	var embedOptionName, embedOptionValue string
	for _, option := range i.Options {
		optionAuthor = nil
		if option.AuthorID != "" {
			optionAuthor, _ = state.User(option.AuthorID)
		}

		embedOptionName = ""
		embedOptionValue = ""

		embedOptionName += titleify(option.Key)

		if optionAuthor != nil && !optionAuthor.Bot && optionAuthor.ID != option.NewValue {
			embedOptionName += " By " + optionAuthor.String() + " #" + optionAuthor.ID
		}

		if option.PreviousValue != "" {
			embedOptionValue += option.Type.String(state, i.GuildID, option.PreviousValue) + " ➡ "
		}
		if option.NewValue != "" {
			embedOptionValue += option.Type.String(state, i.GuildID, option.NewValue)
		} else {
			embedOptionValue += "_/_"
		}

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  embedOptionName,
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
		embed.Description += "On " + i.TargetType.String(state, i.GuildID, i.TargetValue)

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

	if i.ActionType.Destructive() {
		embed.Color = discordColourRed
	}

	if i.WaitingForAuditLogBackfill {
		embed.Color = discordColourOrange
	}

	return embed
}

func (i *Item) Summary(state *state.State, highlightID string) string {
	var summary string
	summary += "**" + i.ActionType.String() + ":**"

	// TODO: add options to summary?

	if i.AuthorID != "" {
		author, err := state.User(i.AuthorID)
		if err != nil {
			author = &discordgo.User{
				ID:       author.ID,
				Username: "N/A",
			}
		}

		summary += " By "
		if author.ID == highlightID {
			summary += "**"
		}
		summary += discord.EscapeDiscordStrict(author.String()) + " #" + author.ID
		if author.ID == highlightID {
			summary += "**"
		}
	}

	if i.TargetValue != "" {
		summary += " On "
		if i.TargetValue == highlightID {
			summary += "**"
		}
		summary += i.TargetType.StringWithoutMention(state, i.GuildID, i.TargetValue)
		if i.TargetValue == highlightID {
			summary += "**"
		}
	}

	for _, option := range i.Options {
		switch option.Key {
		case "message code":
			if option.NewValue != "" && (i.ActionType == ActionTypeModDM || i.ActionType == ActionTypeModNote) {
				summary += "\nMessage: `" + discord.EscapeDiscordStrict(option.NewValue) + "`"
			}
		case "reason":
			if option.NewValue != "" {
				optionAuthor, err := state.User(option.AuthorID)
				if err != nil {
					optionAuthor = &discordgo.User{
						Username: "N/A",
					}
				}

				summary += "\nReason: `" + discord.EscapeDiscordStrict(option.NewValue) + "` by "
				if optionAuthor.ID == highlightID {
					summary += "**"
				}
				summary += discord.EscapeDiscordStrict(optionAuthor.String()) + " #" + optionAuthor.ID
				if optionAuthor.ID == highlightID {
					summary += "**"
				}
			}
		}
	}

	// TODO: message link?
	summary += "\n_#" + i.UUID.String() + "_"

	return summary
}

type ItemOption struct {
	gorm.Model
	ItemID uint `gorm:"NOT NULL;unique_index:idx_itemid_authorid_key"`

	AuthorID      string `gorm:"unique_index:idx_itemid_authorid_key"`
	Key           string `gorm:"NOT NULL;unique_index:idx_itemid_authorid_key"`
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
