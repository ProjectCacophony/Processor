package stats

import (
	"errors"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/regexp"
)

func findInvite(event *events.Event) (*discordgo.Invite, error) {
	for _, field := range event.Fields()[2:] {
		inviteCode := regexp.DiscordInviteRegexp.FindStringSubmatch(field)
		if len(inviteCode) < 7 {
			if regexp.DiscordInviteCodeRegexp.MatchString(field) {

				invite, err := event.Discord().Client.InviteWithCounts(field)
				if err == nil {
					return invite, nil
				}
			}

			continue
		}

		invite, err := event.Discord().Client.InviteWithCounts(inviteCode[5])
		if err == nil {
			return invite, nil
		}
	}

	return nil, errors.New("invite not found")
}

func (p *Plugin) handleInvite(event *events.Event) {
	invite, err := findInvite(event)
	if err != nil {
		if strings.Contains(err.Error(), "invite not found") {
			event.Respond("stats.invite.not-found")
			return
		}
		event.Except(err)
		return
	}

	var detailed bool
	var parentChannel *discordgo.Channel

	// get detailed invite information if possible
	_, err = event.State().Guild(invite.Guild.ID)
	if err == nil {
		botID, err := event.State().BotForGuild(invite.Guild.ID, discordgo.PermissionManageServer)
		if err == nil {
			client, err := discord.NewSession(p.tokens, botID)
			if err != nil {
				event.Except(err)
				return
			}

			// TODO: cache guild invites
			invites, err := client.Client.GuildInvites(invite.Guild.ID)
			if err != nil {
				event.Except(err)
				return
			}

			for _, item := range invites {
				if item.Code == invite.Code {
					invite = item
					detailed = true
					break
				}
			}
		}

		if invite.Channel != nil {
			channel, err := event.State().Channel(invite.Channel.ID)
			if err == nil {
				invite.Channel = channel
			}

			if invite.Channel.ParentID != "" {
				parentChannel, _ = event.State().Channel(invite.Channel.ParentID)
			}
		}
	}

	var createdAt time.Time
	if invite.CreatedAt != "" {
		createdAt, _ = invite.CreatedAt.Parse()
	}

	var iconURL, splashURL, bannerURL string
	if invite.Guild != nil {
		if invite.Guild.Icon != "" {
			iconURL = invite.Guild.IconURL() + "?size=2048"
		}
		if invite.Guild.Splash != "" {
			splashURL = discordgo.EndpointGuildSplash(invite.Guild.ID, invite.Guild.Splash) + "?size=2048"
		}
		if invite.Guild.Banner != "" {
			bannerURL = discordgo.EndpointGuildBanner(invite.Guild.ID, invite.Guild.Banner) + "?size=2048"
		}
	}

	maxAge := (time.Duration(invite.MaxAge) * time.Second).Round(time.Second)

	_, err = event.RespondComplex(
		&discordgo.MessageSend{
			Embed: &discordgo.MessageEmbed{
				Title:       "stats.invite.embed.title",
				URL:         "stats.invite.embed.url",
				Description: "stats.invite.embed.description",
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: "stats.invite.embed.thumbnail.url",
				},
				Image: &discordgo.MessageEmbedImage{
					URL: "stats.invite.embed.image.url",
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text: "stats.invite.embed.footer.text",
				},
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:  "stats.invite.embed.field.created-at.name",
						Value: "stats.invite.embed.field.created-at.value",
					},
				},
			},
		},
		"invite", invite,
		"createdAt", createdAt,
		"iconURL", iconURL,
		"splashURL", splashURL,
		"bannerURL", bannerURL,
		"maxAge", maxAge,
		"detailed", detailed,
		"parentChannel", parentChannel,
	)
	event.Except(err)
}
