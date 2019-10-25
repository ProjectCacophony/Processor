package stats

import (
	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/events"
)

var (
	avatarSizes = [...]string{
		"1024",
		"512",
		"256",
		"128",
		"64",
		"32",
		"16",
	}
)

func (p *Plugin) handleAvatar(event *events.Event) {

	var user *discordgo.User
	var err error
	size := avatarSizes[0]

	if len(event.Fields()) > 1 {
		var validSize bool
		for _, v := range avatarSizes {
			if event.Fields()[1] == v {
				size = v
				validSize = true
				break
			}
		}

		// if the second field isn't a valid size, check to see if that is a user mention
		if !validSize {
			user, _ = event.State().UserFromMention(event.Fields()[1])
		}

		if !validSize && user == nil {
			event.Respond("stats.avatar.invalid-size")
			return
		}
	}

	if user == nil {
		user, err = event.FindUser()
		if err != nil {
			event.Except(err)
			return
		}

		if user == nil {
			event.Respond("stats.avatar.no-user")
			return
		}
	}

	_, err = event.RespondComplex(&discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{
				Name: user.String(),
				URL:  user.AvatarURL(size),
			},
			Image: &discordgo.MessageEmbedImage{
				URL: user.AvatarURL(size),
			},
		},
	})
	if err != nil {
		event.Except(err)
		return
	}
}
