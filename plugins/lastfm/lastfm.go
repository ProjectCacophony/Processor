package lastfm

// func displayTopArtists(event *events.Event) {
// 	// initialise variables
// 	var newArgs []string
// 	var period lastfm_client.LastFmPeriod
// 	var makeCollage bool
// 	period, newArgs = lastfm_client.LastFmGetPeriodFromArgs(event.Args)
// 	// makeCollage, newArgs = isCollageRequest(newArgs)
//
// 	// get lastFM username to look up
// 	var lastfmUsername string
// 	// if len(event.MessageCreate.Mentions) > 0 {
// 	// 	lastfmUsername = getLastFmUsername(ctx, event.MessageCreate.Mentions[0].ID)
// 	// }
// 	if lastfmUsername == "" && len(newArgs) >= 3 {
// 		lastfmUsername = event.Args[2]
// 	}
// 	// if lastfmUsername == "" {
// 	// 	lastfmUsername = getLastFmUsername(ctx, event.MessageCreate.Author.ID)
// 	// }
// 	// if no username found, post error and stop
// 	if lastfmUsername == "" {
// 		event.Respond("no user passed") // nolint: errcheck
// 		return
// 	}
//
// 	// start typing
// 	event.GoType()
//
// 	// lookup user
// 	userInfo, err := lastfm_client.LastFmGetUserinfo(ctx, lastfmUsername)
// 	if err != nil && strings.Contains(err.Error(), "User not found") {
// 		event.SendMessage(event.MessageCreate.ChannelID, "LastFmUserNotFound") // nolint: errcheck, gas
// 		return
// 	}
// 	dhelpers.CheckErr(err)
//
// 	// get basic embed for user
// 	embed := getLastfmUserBaseEmbed(userInfo)
//
// 	// get top artists
// 	var artists []lastfm_client.LastFmArtistData
// 	artists, err = lastfm_client.LastFmGetTopArtists(ctx, userInfo.Username, 10, period)
// 	dhelpers.CheckErr(err)
//
// 	// if no artists found, post error and stop
// 	if len(artists) < 1 {
// 		_, err = event.SendMessage(event.MessageCreate.ChannelID, dhelpers.Tf("LastFmNoScrobbles", "userData", userInfo))
// 		dhelpers.CheckErr(err)
// 		return
// 	}
//
// 	// set content
// 	embed.Author.Name = dhelpers.Tf("LastFmTopArtistsTitle", "userData", userInfo, "period", period)
//
// 	// create collage if requested
// 	if makeCollage {
// 		// initialise variables
// 		imageUrls := make([]string, 0)
// 		artistNames := make([]string, 0)
// 		for _, artist := range artists {
// 			imageUrls = append(imageUrls, artist.ImageURL)
// 			artistNames = append(artistNames, artist.Name)
// 			if len(imageUrls) >= 9 {
// 				break
// 			}
// 		}
//
// 		// create the collage
// 		collageBytes := collage.FromUrls(
// 			ctx,
// 			imageUrls,
// 			artistNames,
// 			900, 900,
// 			300, 300,
// 			dhelpers.DiscordDarkThemeBackgroundColor,
// 		)
//
// 		// add collage image to embed
// 		embed.Image = &discordgo.MessageEmbedImage{
// 			URL: "attachment://LastFM-Collage.png",
// 		}
// 		// send collage to discord and stop
// 		_, err = event.SendComplex(event.MessageCreate.ChannelID, &discordgo.MessageSend{
// 			Files: []*discordgo.File{
// 				{
// 					Name:   "LastFM-Collage.png",
// 					Reader: bytes.NewReader(collageBytes),
// 				},
// 			},
// 			Embed: &embed,
// 		})
// 		dhelpers.CheckErr(err)
// 		return
// 	}
//
// 	// add artists to embed
// 	for i, artist := range artists {
// 		embed.Description += fmt.Sprintf("`#%2d`", i+1) + " " + dhelpers.Tf("LastFmArtist", "artist", artist) + "\n"
// 	}
//
// 	// add artists image to embed if possible
// 	if artists[0].ImageURL != "" {
// 		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
// 			URL: artists[0].ImageURL,
// 		}
// 	}
//
// 	// send to discord
// 	_, err = event.SendEmbed(event.MessageCreate.ChannelID, &embed)
// 	dhelpers.CheckErr(err)
// }

//
// func displayTopTracks(event *events.Event) {
// 	// initialise variables
// 	var newArgs []string
// 	var period lastfm_client.LastFmPeriod
// 	var makeCollage bool
// 	period, newArgs = lastfm_client.LastFmGetPeriodFromArgs(event.Args)
// 	makeCollage, newArgs = isCollageRequest(newArgs)
//
// 	// get lastFM username to look up
// 	var lastfmUsername string
// 	if len(event.MessageCreate.Mentions) > 0 {
// 		lastfmUsername = getLastFmUsername(ctx, event.MessageCreate.Mentions[0].ID)
// 	}
// 	if lastfmUsername == "" && len(newArgs) >= 3 {
// 		lastfmUsername = event.Args[2]
// 	}
// 	if lastfmUsername == "" {
// 		lastfmUsername = getLastFmUsername(ctx, event.MessageCreate.Author.ID)
// 	}
// 	// if no username found, post error and stop
// 	if lastfmUsername == "" {
// 		event.SendMessagef(event.MessageCreate.ChannelID, "LastFmNoUserPassed") // nolint: errcheck
// 		return
// 	}
//
// 	// start typing
// 	event.GoType()
//
// 	// lookup user
// 	userInfo, err := lastfm_client.LastFmGetUserinfo(ctx, lastfmUsername)
// 	if err != nil && strings.Contains(err.Error(), "User not found") {
// 		event.SendMessage(event.MessageCreate.ChannelID, "LastFmUserNotFound") // nolint: errcheck
// 		return
// 	}
// 	dhelpers.CheckErr(err)
//
// 	// get basic embed for user
// 	embed := getLastfmUserBaseEmbed(userInfo)
//
// 	// get top tracks
// 	var tracks []lastfm_client.LastFmTrackData
// 	tracks, err = lastfm_client.LastFmGetTopTracks(ctx, userInfo.Username, 10, period)
// 	dhelpers.CheckErr(err)
//
// 	// if no tracks found, post error and stop
// 	if len(tracks) < 1 {
// 		_, err = event.SendMessage(event.MessageCreate.ChannelID, dhelpers.Tf("LastFmNoScrobbles", "userData", userInfo))
// 		dhelpers.CheckErr(err)
// 		return
// 	}
//
// 	// create collage if requested
// 	if makeCollage {
// 		// initialise variables
// 		imageUrls := make([]string, 0)
// 		trackNames := make([]string, 0)
// 		for _, track := range tracks {
// 			imageUrls = append(imageUrls, track.ImageURL)
// 			trackNames = append(trackNames, track.Name)
// 			if len(imageUrls) >= 9 {
// 				break
// 			}
// 		}
//
// 		// create the collage
// 		collageBytes := collage.FromUrls(
// 			ctx,
// 			imageUrls,
// 			trackNames,
// 			900, 900,
// 			300, 300,
// 			dhelpers.DiscordDarkThemeBackgroundColor,
// 		)
//
// 		// add collage image to embed
// 		embed.Image = &discordgo.MessageEmbedImage{
// 			URL: "attachment://LastFM-Collage.png",
// 		}
// 		// send collage to discord and stop
// 		_, err = event.SendComplex(event.MessageCreate.ChannelID, &discordgo.MessageSend{
// 			Files: []*discordgo.File{
// 				{
// 					Name:   "LastFM-Collage.png",
// 					Reader: bytes.NewReader(collageBytes),
// 				},
// 			},
// 			Embed: &embed,
// 		})
// 		dhelpers.CheckErr(err)
// 		return
// 	}
//
// 	// set embed title
// 	embed.Author.Name = dhelpers.Tf("LastFmTopTracksTitle", "userData", userInfo, "period", period)
//
// 	// add tracks to embed
// 	for i, track := range tracks {
// 		embed.Description += fmt.Sprintf("`#%2d`", i+1) + " " + dhelpers.Tf("LastFmTrack", "track", track) + "\n"
// 	}
//
// 	// add track image to embed if possible
// 	if tracks[0].ImageURL != "" {
// 		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
// 			URL: tracks[0].ImageURL,
// 		}
// 	}
//
// 	// send to discord
// 	_, err = event.SendEmbed(event.MessageCreate.ChannelID, &embed)
// 	dhelpers.CheckErr(err)
// }
//
// func displayTopAlbums(event *events.Event) {
// 	// initialise variables
// 	var newArgs []string
// 	var period lastfm_client.LastFmPeriod
// 	var makeCollage bool
// 	period, newArgs = lastfm_client.LastFmGetPeriodFromArgs(event.Args)
// 	makeCollage, newArgs = isCollageRequest(newArgs)
//
// 	// get lastFM username to look up
// 	var lastfmUsername string
// 	if len(event.MessageCreate.Mentions) > 0 {
// 		lastfmUsername = getLastFmUsername(ctx, event.MessageCreate.Mentions[0].ID)
// 	}
// 	if lastfmUsername == "" && len(newArgs) >= 3 {
// 		lastfmUsername = event.Args[2]
// 	}
// 	if lastfmUsername == "" {
// 		lastfmUsername = getLastFmUsername(ctx, event.MessageCreate.Author.ID)
// 	}
// 	// if no username found, post error and stop
// 	if lastfmUsername == "" {
// 		event.SendMessagef(event.MessageCreate.ChannelID, "LastFmNoUserPassed") // nolint: errcheck
// 		return
// 	}
//
// 	// start typing
// 	event.GoType()
//
// 	// lookup user
// 	userInfo, err := lastfm_client.LastFmGetUserinfo(ctx, lastfmUsername)
// 	if err != nil && strings.Contains(err.Error(), "User not found") {
// 		event.SendMessage(event.MessageCreate.ChannelID, "LastFmUserNotFound") // nolint: errcheck
// 		return
// 	}
// 	dhelpers.CheckErr(err)
//
// 	// get basic embed for user
// 	embed := getLastfmUserBaseEmbed(userInfo)
//
// 	// get top albums
// 	var albums []lastfm_client.LastFmAlbumData
// 	albums, err = lastfm_client.LastFmGetTopAlbums(ctx, userInfo.Username, 10, period)
// 	dhelpers.CheckErr(err)
//
// 	// if no albums found, post error and stop
// 	if len(albums) < 1 {
// 		_, err = event.SendMessage(event.MessageCreate.ChannelID, dhelpers.Tf("LastFmNoScrobbles", "userData", userInfo))
// 		dhelpers.CheckErr(err)
// 		return
// 	}
//
// 	// set content
// 	embed.Author.Name = dhelpers.Tf("LastFmTopAlbumsTitle", "userData", userInfo, "period", period)
//
// 	// create collage if requested
// 	if makeCollage {
// 		// initialise variables
// 		imageUrls := make([]string, 0)
// 		albumNames := make([]string, 0)
// 		for _, album := range albums {
// 			imageUrls = append(imageUrls, album.ImageURL)
// 			albumNames = append(albumNames, album.Name)
// 			if len(imageUrls) >= 9 {
// 				break
// 			}
// 		}
//
// 		// create the collage
// 		collageBytes := collage.FromUrls(
// 			ctx,
// 			imageUrls,
// 			albumNames,
// 			900, 900,
// 			300, 300,
// 			dhelpers.DiscordDarkThemeBackgroundColor,
// 		)
//
// 		// add collage image to embed
// 		embed.Image = &discordgo.MessageEmbedImage{
// 			URL: "attachment://LastFM-Collage.png",
// 		}
// 		// send collage to discord and stop
// 		_, err = event.SendComplex(event.MessageCreate.ChannelID, &discordgo.MessageSend{
// 			Files: []*discordgo.File{
// 				{
// 					Name:   "LastFM-Collage.png",
// 					Reader: bytes.NewReader(collageBytes),
// 				},
// 			},
// 			Embed: &embed,
// 		})
// 		dhelpers.CheckErr(err)
// 		return
// 	}
//
// 	// add albums to embed
// 	for i, album := range albums {
// 		embed.Description += fmt.Sprintf("`#%2d`", i+1) + " " + dhelpers.Tf("LastFmAlbum", "album", album) + "\n"
// 	}
//
// 	// add album image to embed if possible
// 	if albums[0].ImageURL != "" {
// 		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
// 			URL: albums[0].ImageURL,
// 		}
// 	}
//
// 	// send to discord
// 	_, err = event.SendEmbed(event.MessageCreate.ChannelID, &embed)
// 	dhelpers.CheckErr(err)
// }
//
// func displayRecent(event *events.Event) {
// 	// get lastFM username to look up
// 	var lastfmUsername string
// 	if len(event.MessageCreate.Mentions) > 0 {
// 		lastfmUsername = getLastFmUsername(ctx, event.MessageCreate.Mentions[0].ID)
// 	}
// 	if lastfmUsername == "" && len(event.Args) >= 3 {
// 		lastfmUsername = event.Args[2]
// 	}
// 	if lastfmUsername == "" {
// 		lastfmUsername = getLastFmUsername(ctx, event.MessageCreate.Author.ID)
// 	}
// 	// if no username found, post error and stop
// 	if lastfmUsername == "" {
// 		event.SendMessagef(event.MessageCreate.ChannelID, "LastFmNoUserPassed") // nolint: errcheck
// 		return
// 	}
//
// 	// start typing
// 	event.GoType()
//
// 	// lookup user
// 	userInfo, err := lastfm_client.LastFmGetUserinfo(ctx, lastfmUsername)
// 	if err != nil && strings.Contains(err.Error(), "User not found") {
// 		event.SendMessage(event.MessageCreate.ChannelID, "LastFmUserNotFound") // nolint: errcheck
// 		return
// 	}
// 	dhelpers.CheckErr(err)
//
// 	// get basic embed for user
// 	embed := getLastfmUserBaseEmbed(userInfo)
//
// 	// get recent tracks
// 	var tracks []lastfm_client.LastFmTrackData
// 	tracks, err = lastfm_client.LastFmGetRecentTracks(ctx, userInfo.Username, 10)
// 	dhelpers.CheckErr(err)
//
// 	// if no tracks found, post error and stop
// 	if len(tracks) < 1 {
// 		_, err = event.SendMessage(event.MessageCreate.ChannelID, dhelpers.Tf("LastFmNoScrobbles", "userData", userInfo))
// 		dhelpers.CheckErr(err)
// 		return
// 	}
//
// 	// set embed title
// 	embed.Author.Name = dhelpers.Tf("LastFmRecentTitle", "userData", userInfo)
//
// 	// add tracks to embed
// 	for _, track := range tracks {
// 		embed.Description += dhelpers.Tf("LastFmTrackLong", "track", track, "hidenp", true)
//
// 		if track.NowPlaying {
// 			embed.Description += " - _" + dhelpers.T("LastFmNowPlaying") + "_"
// 		} else if !track.Time.IsZero() {
// 			embed.Description += " - " + humanize.Time(track.Time)
// 		}
//
// 		embed.Description += "\n"
// 	}
//
// 	// send to discord
// 	_, err = event.SendEmbed(event.MessageCreate.ChannelID, &embed)
// 	dhelpers.CheckErr(err)
// }
//
//
// func displayAbout(event *events.Event) {
// 	// get lastFM username to look up
// 	var lastfmUsername string
// 	if len(event.MessageCreate.Mentions) > 0 {
// 		lastfmUsername = getLastFmUsername(ctx, event.MessageCreate.Mentions[0].ID)
// 	}
// 	if lastfmUsername == "" && len(event.Args) >= 2 {
// 		lastfmUsername = event.Args[1]
// 	}
// 	if lastfmUsername == "" {
// 		lastfmUsername = getLastFmUsername(ctx, event.MessageCreate.Author.ID)
// 	}
// 	// if no username found, post error and stop
// 	if lastfmUsername == "" {
// 		event.SendMessagef(event.MessageCreate.ChannelID, "LastFmNoUserPassed") // nolint: errcheck
// 		return
// 	}
//
// 	// start typing
// 	event.GoType()
//
// 	// lookup user
// 	userInfo, err := lastfm_client.LastFmGetUserinfo(ctx, lastfmUsername)
// 	if err != nil && strings.Contains(err.Error(), "User not found") {
// 		event.SendMessage(event.MessageCreate.ChannelID, "LastFmUserNotFound") // nolint: errcheck
// 		return
// 	}
// 	dhelpers.CheckErr(err)
//
// 	// get basic embed for user
// 	embed := getLastfmUserBaseEmbed(userInfo)
// 	embed.Author.Name = dhelpers.Tf("LastFmAboutTitle", "userData", userInfo)
// 	dhelpers.CheckErr(err)
//
// 	// add fields
// 	// replace scrobbles count in footer with field
// 	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
// 		Name:   "🎶 Scrobbles",
// 		Value:  humanize.Comma(int64(userInfo.Scrobbles)),
// 		Inline: false,
// 	})
// 	if strings.Contains(embed.Footer.Text, "|") {
// 		embed.Footer.Text = strings.TrimSpace(strings.SplitN(embed.Footer.Text, "|", 2)[0])
// 	}
//
// 	// add country to embed if possible
// 	if userInfo.Country != "" {
// 		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
// 			Name:   "🗺 Country",
// 			Value:  userInfo.Country,
// 			Inline: false,
// 		})
// 	}
//
// 	// add account creation date to embed if possible
// 	if !userInfo.AccountCreation.IsZero() {
// 		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
// 			Name:   "🗓 Account creation",
// 			Value:  humanize.Time(userInfo.AccountCreation),
// 			Inline: false,
// 		})
// 	}
//
// 	// replace author icon with bigger thumbnail for about
// 	if userInfo.Icon != "" {
// 		embed.Author.IconURL = ""
// 		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
// 			URL: userInfo.Icon,
// 		}
// 	}
//
// 	// send to discord
// 	_, err = event.SendEmbed(event.MessageCreate.ChannelID, &embed)
// 	dhelpers.CheckErr(err)
// }
//
// func handleSet(event *events.Event) {
// 	// we need at least three args
// 	if len(event.Args) < 3 {
// 		return
// 	}
//
// 	// get last.fm username from args
// 	username := event.Args[2]
//
// 	// upsert username to db
// 	err := models.LastFmRepository.Upsert(
// 		ctx,
// 		bson.NewDocument(bson.EC.String("userid", event.MessageCreate.Author.ID)),
// 		bson.NewDocument(
// 			bson.EC.Interface(
// 				"$set",
// 				models.LastFmEntry{
// 					UserID:         event.MessageCreate.Author.ID,
// 					LastFmUsername: username,
// 				},
// 			),
// 		),
// 	)
// 	dhelpers.CheckErr(err)
//
// 	// send to discord
// 	_, err = event.SendMessage(event.MessageCreate.ChannelID, "LastFmUsernameSaved")
// 	dhelpers.CheckErr(err)
// }
//
// func displayServerTopTracks(event *events.Event) {
// 	// initialise variables
// 	var err error
// 	var period lastfm_client.LastFmPeriod
// 	period, _ = lastfm_client.LastFmGetPeriodFromArgs(event.Args)
//
// 	// start typing
// 	event.GoType()
//
// 	// get source guild
// 	var guild *discordgo.Guild
// 	guild, err = state.Guild(event.MessageCreate.GuildID)
// 	dhelpers.CheckErr(err)
//
// 	// get stats data from redis and unmarshal
// 	statsBytes, err := cache.GetRedisClient().Get(lastfm_client.LastFmGuildTopTracksKey(guild.ID, period)).Bytes()
// 	if err == redis.Nil {
// 		_, err = event.SendMessage(event.MessageCreate.ChannelID, dhelpers.Tf("LastFmGuildNoScrobbles"))
// 		dhelpers.CheckErr(err)
// 		return
// 	}
// 	dhelpers.CheckErr(err)
//
// 	var stats lastfm_client.LastFmGuildTopTracks
// 	err = jsoniter.Unmarshal(statsBytes, &stats)
// 	dhelpers.CheckErr(err)
//
// 	// get basic embed for user
// 	embed := getLastfmGuildBaseEmbed(guild, stats.NumberOfUsers)
//
// 	// if no tracks found, post error and stop
// 	if len(stats.Tracks) < 1 {
// 		_, err = event.SendMessage(event.MessageCreate.ChannelID, dhelpers.Tf("LastFmGuildNo"))
// 		dhelpers.CheckErr(err)
// 		return
// 	}
//
// 	// set embed title, footer, and timestamp
// 	embed.Author.Name = dhelpers.Tf("LastFmGuildTopTracksTitle", "guild", guild, "period", period)
// 	embed.Footer.Text += " | " + dhelpers.T("LastFmCachedAt")
// 	embed.Timestamp = dhelpers.DiscordTime(stats.CachedAt)
//
// 	// add tracks to embed
// 	for i, track := range stats.Tracks {
// 		embed.Description += fmt.Sprintf("`#%2d`", i+1) + " " + dhelpers.Tf(
// 			"LastFmTrack", "track", track) + "\n"
// 		if i >= 9 {
// 			break
// 		}
// 	}
//
// 	// add track image to embed if possible
// 	if stats.Tracks[0].ImageURL != "" {
// 		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
// 			URL: stats.Tracks[0].ImageURL,
// 		}
// 	}
//
// 	// send to discord
// 	_, err = event.SendEmbed(event.MessageCreate.ChannelID, &embed)
// 	dhelpers.CheckErr(err)
// }
