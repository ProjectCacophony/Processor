package gall

import "gitlab.com/project-d-collab/dhelpers"

var (
	color = dhelpers.GetDiscordColorFromHex("#4064D0")

	friendlyBoardURL = func(boardID string) string {
		return "http://gall.dcinside.com/board/lists/?id=" + boardID + "&page=1&exception_mode=recommend"
	}
)
