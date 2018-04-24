package lastfm

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"github.com/globalsign/mgo/bson"
	"gitlab.com/project-d-collab/SqsProcessor/models"
	"gitlab.com/project-d-collab/dhelpers"
	"gitlab.com/project-d-collab/dhelpers/collage"
	"gitlab.com/project-d-collab/dhelpers/mdb"
)

func displayTopArtists(event dhelpers.EventContainer) {
	var newArgs []string
	var period dhelpers.LastFmPeriod
	var makeCollage bool
	period, newArgs = dhelpers.LastFmGetPeriodFromArgs(event.Args)
	makeCollage, newArgs = isCollageRequest(newArgs)

	var lastfmUsername string
	if len(event.MessageCreate.Mentions) > 0 {
		lastfmUsername = getLastFmUsername(event.MessageCreate.Mentions[0].ID)
	}
	if lastfmUsername == "" && len(newArgs) >= 3 {
		lastfmUsername = event.Args[2]
	}
	if lastfmUsername == "" {
		lastfmUsername = getLastFmUsername(event.MessageCreate.Author.ID)
	}

	if lastfmUsername == "" {
		event.SendMessage(event.MessageCreate.ChannelID, "LastFmNoUserPassed") // nolint: errcheck
		return
	}

	event.GoType(event.MessageCreate.ChannelID)

	// look user
	userInfo, err := dhelpers.LastFmGetUserinfo(lastfmUsername)
	if err != nil && strings.Contains(err.Error(), "User not found") {
		event.SendMessage(event.MessageCreate.ChannelID, "LastFmUserNotFound") // nolint: errcheck, gas
		return
	}
	dhelpers.CheckErr(err)

	// get basic embed for user
	embed := getLastfmBaseEmbed(userInfo)

	// get top artists
	var artists []dhelpers.LastfmArtistData
	artists, err = dhelpers.LastFmGetTopArtists(userInfo.Username, 10, period)
	dhelpers.CheckErr(err)

	if len(artists) < 1 {
		_, err = event.SendMessage(event.MessageCreate.ChannelID, dhelpers.Tf("LastFmNoScrobbles", "userData", userInfo))
		dhelpers.CheckErr(err)
		return
	}

	// set content
	embed.Author.Name = dhelpers.Tf("LastFmTopArtistsTitle", "userData", userInfo, "period", period)

	if makeCollage {
		imageUrls := make([]string, 0)
		artistNames := make([]string, 0)
		for _, artist := range artists {
			imageUrls = append(imageUrls, artist.ImageURL)
			artistNames = append(artistNames, artist.Name)
			if len(imageUrls) >= 9 {
				break
			}
		}

		collageBytes := collage.FromUrls(
			imageUrls,
			artistNames,
			900, 900,
			300, 300,
			dhelpers.DiscordDarkThemeBackgroundColor,
		)

		embed.Image = &discordgo.MessageEmbedImage{
			URL: "attachment://LastFM-Collage.png",
		}
		_, err = event.SendComplex(event.MessageCreate.ChannelID, &discordgo.MessageSend{
			Files: []*discordgo.File{
				{
					Name:   "LastFM-Collage.png",
					Reader: bytes.NewReader(collageBytes),
				},
			},
			Embed: &embed,
		})
		dhelpers.CheckErr(err)
		return
	}

	for i, artist := range artists {
		embed.Description += fmt.Sprintf("`#%2d`", i+1) + " " + dhelpers.Tf("LastFmArtist", "artist", artist, "scrobbles", humanize.Comma(int64(artist.Scrobbles))) + "\n"
	}

	if artists[0].ImageURL != "" {
		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
			URL: artists[0].ImageURL,
		}
	}

	_, err = event.SendEmbed(event.MessageCreate.ChannelID, &embed)
	dhelpers.CheckErr(err)
}

func displayTopTracks(event dhelpers.EventContainer) {
	var newArgs []string
	var period dhelpers.LastFmPeriod
	var makeCollage bool
	period, newArgs = dhelpers.LastFmGetPeriodFromArgs(event.Args)
	makeCollage, newArgs = isCollageRequest(newArgs)

	var lastfmUsername string
	if len(event.MessageCreate.Mentions) > 0 {
		lastfmUsername = getLastFmUsername(event.MessageCreate.Mentions[0].ID)
	}
	if lastfmUsername == "" && len(newArgs) >= 3 {
		lastfmUsername = event.Args[2]
	}
	if lastfmUsername == "" {
		lastfmUsername = getLastFmUsername(event.MessageCreate.Author.ID)
	}

	if lastfmUsername == "" {
		event.SendMessage(event.MessageCreate.ChannelID, "LastFmNoUserPassed") // nolint: errcheck
		return
	}

	event.GoType(event.MessageCreate.ChannelID)

	// look user
	userInfo, err := dhelpers.LastFmGetUserinfo(lastfmUsername)
	if err != nil && strings.Contains(err.Error(), "User not found") {
		event.SendMessage(event.MessageCreate.ChannelID, "LastFmUserNotFound") // nolint: errcheck
		return
	}
	dhelpers.CheckErr(err)

	// get basic embed for user
	embed := getLastfmBaseEmbed(userInfo)

	// get top artists
	var tracks []dhelpers.LastfmTrackData
	tracks, err = dhelpers.LastFmGetTopTracks(userInfo.Username, 10, period)
	dhelpers.CheckErr(err)

	if makeCollage {
		imageUrls := make([]string, 0)
		trackNames := make([]string, 0)
		for _, track := range tracks {
			imageUrls = append(imageUrls, track.ImageURL)
			trackNames = append(trackNames, track.Name)
			if len(imageUrls) >= 9 {
				break
			}
		}

		collageBytes := collage.FromUrls(
			imageUrls,
			trackNames,
			900, 900,
			300, 300,
			dhelpers.DiscordDarkThemeBackgroundColor,
		)

		embed.Image = &discordgo.MessageEmbedImage{
			URL: "attachment://LastFM-Collage.png",
		}
		_, err = event.SendComplex(event.MessageCreate.ChannelID, &discordgo.MessageSend{
			Files: []*discordgo.File{
				{
					Name:   "LastFM-Collage.png",
					Reader: bytes.NewReader(collageBytes),
				},
			},
			Embed: &embed,
		})
		dhelpers.CheckErr(err)
		return
	}

	if len(tracks) < 1 {
		_, err = event.SendMessage(event.MessageCreate.ChannelID, dhelpers.Tf("LastFmNoScrobbles", "userData", userInfo))
		dhelpers.CheckErr(err)
		return
	}

	// set content
	embed.Author.Name = dhelpers.Tf("LastFmTopTracksTitle", "userData", userInfo, "period", period)

	for i, track := range tracks {
		embed.Description += fmt.Sprintf("`#%2d`", i+1) + " " + dhelpers.Tf("LastFmTrack", "track", track, "scrobbles", humanize.Comma(int64(track.Scrobbles))) + "\n"
	}

	if tracks[0].ImageURL != "" {
		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
			URL: tracks[0].ImageURL,
		}
	}

	_, err = event.SendEmbed(event.MessageCreate.ChannelID, &embed)
	dhelpers.CheckErr(err)
}

func displayTopAlbums(event dhelpers.EventContainer) {
	var newArgs []string
	var period dhelpers.LastFmPeriod
	var makeCollage bool
	period, newArgs = dhelpers.LastFmGetPeriodFromArgs(event.Args)
	makeCollage, newArgs = isCollageRequest(newArgs)

	var lastfmUsername string
	if len(event.MessageCreate.Mentions) > 0 {
		lastfmUsername = getLastFmUsername(event.MessageCreate.Mentions[0].ID)
	}
	if lastfmUsername == "" && len(newArgs) >= 3 {
		lastfmUsername = event.Args[2]
	}
	if lastfmUsername == "" {
		lastfmUsername = getLastFmUsername(event.MessageCreate.Author.ID)
	}

	if lastfmUsername == "" {
		event.SendMessage(event.MessageCreate.ChannelID, "LastFmNoUserPassed") // nolint: errcheck
		return
	}

	event.GoType(event.MessageCreate.ChannelID)

	// look user
	userInfo, err := dhelpers.LastFmGetUserinfo(lastfmUsername)
	if err != nil && strings.Contains(err.Error(), "User not found") {
		event.SendMessage(event.MessageCreate.ChannelID, "LastFmUserNotFound") // nolint: errcheck
		return
	}
	dhelpers.CheckErr(err)

	// get basic embed for user
	embed := getLastfmBaseEmbed(userInfo)

	// get top artists
	var albums []dhelpers.LastfmAlbumData
	albums, err = dhelpers.LastFmGetTopAlbums(userInfo.Username, 10, period)
	dhelpers.CheckErr(err)

	if len(albums) < 1 {
		_, err = event.SendMessage(event.MessageCreate.ChannelID, dhelpers.Tf("LastFmNoScrobbles", "userData", userInfo))
		dhelpers.CheckErr(err)
		return
	}

	// set content
	embed.Author.Name = dhelpers.Tf("LastFmTopAlbumsTitle", "userData", userInfo, "period", period)

	if makeCollage {
		imageUrls := make([]string, 0)
		albumNames := make([]string, 0)
		for _, album := range albums {
			imageUrls = append(imageUrls, album.ImageURL)
			albumNames = append(albumNames, album.Name)
			if len(imageUrls) >= 9 {
				break
			}
		}

		collageBytes := collage.FromUrls(
			imageUrls,
			albumNames,
			900, 900,
			300, 300,
			dhelpers.DiscordDarkThemeBackgroundColor,
		)

		embed.Image = &discordgo.MessageEmbedImage{
			URL: "attachment://LastFM-Collage.png",
		}
		_, err = event.SendComplex(event.MessageCreate.ChannelID, &discordgo.MessageSend{
			Files: []*discordgo.File{
				{
					Name:   "LastFM-Collage.png",
					Reader: bytes.NewReader(collageBytes),
				},
			},
			Embed: &embed,
		})
		dhelpers.CheckErr(err)
		return
	}

	for i, album := range albums {
		embed.Description += fmt.Sprintf("`#%2d`", i+1) + " " + dhelpers.Tf("LastFmAlbum", "album", album, "scrobbles", humanize.Comma(int64(album.Scrobbles))) + "\n"
	}

	if albums[0].ImageURL != "" {
		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
			URL: albums[0].ImageURL,
		}
	}

	_, err = event.SendEmbed(event.MessageCreate.ChannelID, &embed)
	dhelpers.CheckErr(err)
}

func displayRecent(event dhelpers.EventContainer) {
	var lastfmUsername string
	if len(event.MessageCreate.Mentions) > 0 {
		lastfmUsername = getLastFmUsername(event.MessageCreate.Mentions[0].ID)
	}
	if lastfmUsername == "" && len(event.Args) >= 3 {
		lastfmUsername = event.Args[2]
	}
	if lastfmUsername == "" {
		lastfmUsername = getLastFmUsername(event.MessageCreate.Author.ID)
	}

	if lastfmUsername == "" {
		event.SendMessage(event.MessageCreate.ChannelID, "LastFmNoUserPassed") // nolint: errcheck
		return
	}

	event.GoType(event.MessageCreate.ChannelID)

	// look user
	userInfo, err := dhelpers.LastFmGetUserinfo(lastfmUsername)
	if err != nil && strings.Contains(err.Error(), "User not found") {
		event.SendMessage(event.MessageCreate.ChannelID, "LastFmUserNotFound") // nolint: errcheck
		return
	}
	dhelpers.CheckErr(err)

	// get basic embed for user
	embed := getLastfmBaseEmbed(userInfo)

	// get recent tracks
	var tracks []dhelpers.LastfmTrackData
	tracks, err = dhelpers.LastFmGetRecentTracks(userInfo.Username, 10)
	dhelpers.CheckErr(err)

	if len(tracks) < 1 {
		_, err = event.SendMessage(event.MessageCreate.ChannelID, dhelpers.Tf("LastFmNoScrobbles", "userData", userInfo))
		dhelpers.CheckErr(err)
		return
	}

	// set content
	embed.Author.Name = dhelpers.Tf("LastFmRecentTitle", "userData", userInfo)

	for _, track := range tracks {
		embed.Description += dhelpers.Tf("LastFmTrackLong", "track", track, "hidenp", true)

		if track.NowPlaying {
			embed.Description += " - _" + dhelpers.T("LastFmNowPlaying") + "_"
		} else if !track.Time.IsZero() {
			embed.Description += " - " + humanize.Time(track.Time)
		}

		embed.Description += "\n"
	}

	_, err = event.SendEmbed(event.MessageCreate.ChannelID, &embed)
	dhelpers.CheckErr(err)
}

func displayNowPlaying(event dhelpers.EventContainer) {
	var lastfmUsername string
	if len(event.MessageCreate.Mentions) > 0 {
		lastfmUsername = getLastFmUsername(event.MessageCreate.Mentions[0].ID)
	}
	if lastfmUsername == "" && len(event.Args) >= 3 {
		lastfmUsername = event.Args[2]
	}
	if lastfmUsername == "" {
		lastfmUsername = getLastFmUsername(event.MessageCreate.Author.ID)
	}

	if lastfmUsername == "" {
		event.SendMessage(event.MessageCreate.ChannelID, "LastFmNoUserPassed") // nolint: errcheck
		return
	}

	event.GoType(event.MessageCreate.ChannelID)

	// look user
	userInfo, err := dhelpers.LastFmGetUserinfo(lastfmUsername)
	if err != nil && strings.Contains(err.Error(), "User not found") {
		event.SendMessage(event.MessageCreate.ChannelID, "LastFmUserNotFound") // nolint: errcheck
		return
	}
	dhelpers.CheckErr(err)

	// get basic embed for user
	embed := getLastfmBaseEmbed(userInfo)

	// get recent tracks
	var tracks []dhelpers.LastfmTrackData
	tracks, err = dhelpers.LastFmGetRecentTracks(userInfo.Username, 2)
	dhelpers.CheckErr(err)

	if len(tracks) < 1 {
		_, err = event.SendMessage(event.MessageCreate.ChannelID, dhelpers.Tf("LastFmNoScrobbles", "userData", userInfo))
		dhelpers.CheckErr(err)
		return
	}

	// set content
	embed.Author.Name = dhelpers.Tf("LastFmNowPlayingTitle", "userData", userInfo, "tracks", tracks)

	embed.Description = dhelpers.Tf("LastFmTrackLong", "track", tracks[0])

	if tracks[0].Album != "" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   dhelpers.T("LastFmAlbumTitle"),
			Value:  tracks[0].Album,
			Inline: false,
		})
	}

	if tracks[0].ImageURL != "" {
		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
			URL: tracks[0].ImageURL,
		}
	}

	if len(tracks) >= 2 {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   dhelpers.T("LastFmListenedBeforeTitle"),
			Value:  dhelpers.Tf("LastFmTrackLong", "track", tracks[1]),
			Inline: false,
		})
	}

	_, err = event.SendEmbed(event.MessageCreate.ChannelID, &embed)
	dhelpers.CheckErr(err)
}

func displayAbout(event dhelpers.EventContainer) {
	var lastfmUsername string
	if len(event.MessageCreate.Mentions) > 0 {
		lastfmUsername = getLastFmUsername(event.MessageCreate.Mentions[0].ID)
	}
	if lastfmUsername == "" && len(event.Args) >= 2 {
		lastfmUsername = event.Args[1]
	}
	if lastfmUsername == "" {
		lastfmUsername = getLastFmUsername(event.MessageCreate.Author.ID)
	}

	if lastfmUsername == "" {
		event.SendMessage(event.MessageCreate.ChannelID, "LastFmNoUserPassed") // nolint: errcheck
		return
	}

	event.GoType(event.MessageCreate.ChannelID)

	// look user
	userInfo, err := dhelpers.LastFmGetUserinfo(lastfmUsername)
	if err != nil && strings.Contains(err.Error(), "User not found") {
		event.SendMessage(event.MessageCreate.ChannelID, "LastFmUserNotFound") // nolint: errcheck
		return
	}
	dhelpers.CheckErr(err)

	// get basic embed for user
	embed := getLastfmBaseEmbed(userInfo)
	embed.Author.Name = dhelpers.Tf("LastFmAboutTitle", "userData", userInfo)
	dhelpers.CheckErr(err)

	// add fields
	// replace scrobbles count in footer with field
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:   "🎶 Scrobbles",
		Value:  humanize.Comma(int64(userInfo.Scrobbles)),
		Inline: false,
	})
	if strings.Contains(embed.Footer.Text, "|") {
		embed.Footer.Text = strings.TrimSpace(strings.SplitN(embed.Footer.Text, "|", 2)[0])
	}

	if userInfo.Country != "" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "🗺 Country",
			Value:  userInfo.Country,
			Inline: false,
		})
	}

	if !userInfo.AccountCreation.IsZero() {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "🗓 Account creation",
			Value:  humanize.Time(userInfo.AccountCreation),
			Inline: false,
		})
	}

	// replace author icon with bigger thumbnail for about
	if userInfo.Icon != "" {
		embed.Author.IconURL = ""
		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
			URL: userInfo.Icon,
		}
	}

	_, err = event.SendEmbed(event.MessageCreate.ChannelID, &embed)
	dhelpers.CheckErr(err)
}

func setUsername(event dhelpers.EventContainer) {
	if len(event.Args) < 3 {
		return
	}

	username := event.Args[2]

	err := mdb.UpsertQuery(
		models.LastFmTable,
		bson.M{"userid": event.MessageCreate.Author.ID},
		models.LastFmEntry{
			UserID:         event.MessageCreate.Author.ID,
			LastFmUsername: username,
		},
	)
	dhelpers.CheckErr(err)

	_, err = event.SendMessage(event.MessageCreate.ChannelID, "LastFmUsernameSaved")
	dhelpers.CheckErr(err)
}
