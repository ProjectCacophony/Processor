// nolint: dupl
package lastfm

import (
	"fmt"
	"strings"

	"github.com/Seklfreak/lastfm-go/lastfm"
	"github.com/bwmarrin/discordgo"
	lastfmclient "gitlab.com/Cacophony/Processor/plugins/lastfm/lastfm-client"
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleTopArtists(event *events.Event, lastfmClient *lastfm.Api, offset int) {
	fields := event.Fields()[offset:]

	// var makeCollage bool
	period, fields := lastfmclient.GetPeriodFromArgs(fields)
	// makeCollage, newArgs = isCollageRequest(newArgs)

	// get lastFM username to look up
	username, _ := extractUsername(event, p.db, fields, -1)

	// lookup user
	userInfo, err := lastfmclient.GetUserinfo(lastfmClient, username)
	if err != nil {
		if strings.Contains(err.Error(), "User not found") {
			event.Respond("lastfm.user-not-found", "username", username) // nolint: errcheck
			return
		}
		event.Except(err)
		return
	}

	// get basic embed for user
	embed := getLastfmUserBaseEmbed(userInfo)

	// get top artists
	var artists []lastfmclient.ArtistData
	artists, err = lastfmclient.GetTopArtists(lastfmClient, userInfo.Username, 10, period)
	if err != nil {
		event.Except(err)
		return
	}

	// if no artists found, post error and stop
	if len(artists) < 1 {
		_, err = event.Respond("lastfm.no-scrobbles", "username", username)
		event.Except(err)
		return
	}

	// set content
	embed.Author.Name = "lastfm.artists.embed.title"

	// // create collage if requested
	// if makeCollage {
	// 	// initialise variables
	// 	imageUrls := make([]string, 0)
	// 	artistNames := make([]string, 0)
	// 	for _, artist := range artists {
	// 		imageUrls = append(imageUrls, artist.ImageURL)
	// 		artistNames = append(artistNames, artist.Name)
	// 		if len(imageUrls) >= 9 {
	// 			break
	// 		}
	// 	}
	//
	// 	// create the collage
	// 	collageBytes := collage.FromUrls(
	// 		ctx,
	// 		imageUrls,
	// 		artistNames,
	// 		900, 900,
	// 		300, 300,
	// 		dhelpers.DiscordDarkThemeBackgroundColor,
	// 	)
	//
	// 	// add collage image to embed
	// 	embed.Image = &discordgo.MessageEmbedImage{
	// 		URL: "attachment://LastFM-Collage.png",
	// 	}
	// 	// send collage to discord and stop
	// 	_, err = event.SendComplex(event.MessageCreate.ChannelID, &discordgo.MessageSend{
	// 		Files: []*discordgo.File{
	// 			{
	// 				Name:   "LastFM-Collage.png",
	// 				Reader: bytes.NewReader(collageBytes),
	// 			},
	// 		},
	// 		Embed: &embed,
	// 	})
	// 	dhelpers.CheckErr(err)
	// 	return
	// }

	// add artists to embed
	for i, artist := range artists {
		embed.Description += fmt.Sprintf("`#%2d`", i+1) + " " +
			event.Translate("lastfm.artist.short", "artist", artist) +
			"\n"
	}

	// add artists image to embed if possible
	if artists[0].ImageURL != "" {
		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
			URL: artists[0].ImageURL,
		}
	}

	// send to discord
	_, err = event.SendComplex(event.MessageCreate.ChannelID, &discordgo.MessageSend{
		Embed: &embed,
	}, "userData", userInfo, "period", period)
	event.Except(err)
}
