package serverlist

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	humanize "github.com/dustin/go-humanize"
)

func getLogApprovedEmbed(server *Server) *discordgo.MessageEmbed {
	baseEmbed := getBaseServerEmbed(server, true)
	baseEmbed.Title = "serverlist.log.approved.embed.title"

	return baseEmbed
}

// func getLogRejectedEmbed(server *Server) *discordgo.MessageEmbed {
// 	baseEmbed := getBaseServerEmbed(server, false)
// 	baseEmbed.Title = "serverlist.log.rejected.embed.title"
//
// 	return baseEmbed
// }
//
// func getLogRemovedEmbed(server *Server) *discordgo.MessageEmbed {
// 	baseEmbed := getBaseServerEmbed(server, false)
// 	baseEmbed.Title = "serverlist.log.removed.embed.title"
//
// 	return baseEmbed
// }
//
// func getLogHiddenEmbed(server *Server) *discordgo.MessageEmbed {
// 	baseEmbed := getBaseServerEmbed(server, false)
// 	baseEmbed.Title = "serverlist.log.hidden.embed.title"
//
// 	return baseEmbed
// }
//
// func getLogUnhiddenEmbed(server *Server) *discordgo.MessageEmbed {
// 	baseEmbed := getBaseServerEmbed(server, true)
// 	baseEmbed.Title = "serverlist.log.unhidden.embed.title"
//
// 	return baseEmbed
// }

func getBaseServerEmbed(server *Server, invite bool) *discordgo.MessageEmbed {
	var categoryText string
	for _, category := range server.Categories {
		categoryText += "<#" + category.Category.ChannelID + ">, "
	}
	categoryText = strings.TrimRight(categoryText, ", ")

	fields := []*discordgo.MessageEmbedField{
		{
			Name:   "🏷 Name",
			Value:  strings.Join(server.Names, "; "),
			Inline: true,
		},
	}
	if invite {
		fields = append(fields,
			&discordgo.MessageEmbedField{
				Name:   "🚩 Invite",
				Value:  fmt.Sprintf("https://discord.gg/%s", server.InviteCode),
				Inline: true,
			},
		)
	}
	fields = append(fields,
		&discordgo.MessageEmbedField{
			Name:   "📖 Description",
			Value:  server.Description,
			Inline: false,
		},
		&discordgo.MessageEmbedField{
			Name:   "🗃 Category",
			Value:  categoryText,
			Inline: false,
		},
	)

	return &discordgo.MessageEmbed{
		Fields: fields,
	}
}

func getQueueMessageEmbed(server *Server, total int) *discordgo.MessageEmbed {
	if server == nil {
		return &discordgo.MessageEmbed{
			Title:       "⌛ Serverlist Queue",
			Description: "Queue empty!",
		}
	}

	var categoryText string
	for _, category := range server.Categories {
		categoryText += "<#" + category.Category.ChannelID + ">, "
	}
	categoryText = strings.TrimRight(categoryText, ", ")

	return &discordgo.MessageEmbed{
		Title:       "⌛ Serverlist Queue",
		Description: "serverlist.queue.embed.description",
		Timestamp:   server.CreatedAt.Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf(
				"there are %d Servers queued in total • added", total,
			),
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name: "🏷 Name(s)",
				Value: fmt.Sprintf("%s\n#%s",
					strings.Join(server.Names, "; "), server.GuildID,
				),
				Inline: true,
			},
			{
				Name: "👥 Editor(s)",
				Value: fmt.Sprintf("<@%s>",
					strings.Join(server.EditorUserIDs, "> <@"),
				),
				Inline: true,
			},
			{
				Name:   "🚩 Invite",
				Value:  fmt.Sprintf("discord.gg/%s", server.InviteCode),
				Inline: true,
			},
			{
				Name:   "📈 Members",
				Value:  humanize.Comma(int64(server.TotalMembers)),
				Inline: true,
			},
			{
				Name:   "📖 Description",
				Value:  server.Description,
				Inline: false,
			},
			{
				Name:   "🗃 Category",
				Value:  categoryText,
				Inline: false,
			},
		},
	}
}
