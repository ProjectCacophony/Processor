package gall

import "gitlab.com/project-d-collab/dhelpers"

var (
	// GallColor is the discord colour used for Gall embeds
	GallColor = dhelpers.GetDiscordColorFromHex("#4064D0")
	// GallIcon is a Gall Logo
	GallIcon = "https://i.imgur.com/SkVJwJP.jpg"

	friendlyBoardURL = func(boardID string) string {
		return "http://gall.dcinside.com/board/lists/?id=" + boardID + "&page=1&exception_mode=recommend"
	}
)
