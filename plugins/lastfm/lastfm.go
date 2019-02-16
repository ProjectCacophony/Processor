package lastfm

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
