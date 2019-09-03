package stats

import (
	"sort"
	"time"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleUser(event *events.Event) {
	user, err := event.FindUser()
	if err != nil {
		event.Except(err)
		return
	}

	createdAt, _ := discordgo.SnowflakeTimestamp(user.ID)

	// optional information for members
	var userNr int
	var joinedAt, premiumSince time.Time
	var roles []*discordgo.Role
	member, err := event.State().Member(event.GuildID, user.ID)
	if err == nil {
		joinedAt, _ = member.JoinedAt.Parse()
		premiumSince, _ = member.PremiumSince.Parse()

		roles = make([]*discordgo.Role, len(member.Roles))
		for i, roleID := range member.Roles {
			role, err := event.State().Role(event.GuildID, roleID)
			if err != nil {
				role = &discordgo.Role{
					ID: roleID,
				}
			}

			roles[i] = role
		}
		sort.Slice(roles, func(i, j int) bool {
			return roles[i].Position > roles[j].Position
		})

		// TODO: guild.Members is not available anymore, is there a different way to get the member nr?
		// 	guild, err := event.State().Guild(event.GuildID)
		// 	if err == nil {
		// 		sort.Slice(guild.Members[:], func(i, j int) bool {
		// 			if guild.Members[i].JoinedAt == "" || guild.Members[j].JoinedAt == "" {
		// 				return false
		// 			}
		//
		// 			iMemberTime, err := guild.Members[i].JoinedAt.Parse()
		// 			if err != nil {
		// 				return false
		// 			}
		// 			jMemberTime, err := guild.Members[j].JoinedAt.Parse()
		// 			if err != nil {
		// 				return false
		// 			}
		//
		// 			return iMemberTime.Before(jMemberTime)
		// 		})
		//
		// 		for i, sortedMember := range guild.Members[:] {
		// 			if sortedMember.User.ID == user.ID {
		// 				userNr = i + 1
		// 				break
		// 			}
		// 		}
		// 	}
	}

	var color int
	for _, role := range roles {
		if role.Color > 0 {
			color = role.Color
			break
		}
	}

	_, err = event.RespondComplex(&discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title:       "stats.user.embed.title",
			Description: "stats.user.embed.description",
			Color:       color,
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: "stats.user.embed.thumbnail.url",
			},
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "stats.user.embed.field.created-at.name",
					Value: "stats.user.embed.field.created-at.value",
				},
				{
					Name:  "stats.user.embed.field.joined-at.name",
					Value: "stats.user.embed.field.joined-at.value",
				},
				{
					Name:  "stats.user.embed.field.premium-since.name",
					Value: "stats.user.embed.field.premium-since.value",
				},
				{
					Name:  "stats.user.embed.field.roles.name",
					Value: "stats.user.embed.field.roles.value",
				},
			},
		},
	},
		"user", user,
		"member", member,
		"createdAt", createdAt,
		"joinedAt", joinedAt,
		"premiumSince", premiumSince,
		"roles", roles,
		"userNr", userNr,
	)
	event.Except(err)
}
