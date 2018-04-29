package gall

import (
	"regexp"

	"github.com/globalsign/mgo/bson"
	"gitlab.com/project-d-collab/SqsProcessor/models"
	"gitlab.com/project-d-collab/dhelpers"
	"gitlab.com/project-d-collab/dhelpers/mdb"
)

var (
	// GallColor is the discord colour used for Gall embeds
	GallColor = dhelpers.GetDiscordColorFromHex("#4064D0")
	// GallIcon is a Gall Logo
	GallIcon = "https://i.imgur.com/SkVJwJP.jpg"

	friendlyBoardURL = func(boardID string) string {
		return "http://gall.dcinside.com/board/lists/?id=" + boardID + "&page=1&exception_mode=recommend"
	}
)

// alreadySetUp returns true if the board is already set up in the channel
func alreadySetUp(boardID, channelID string) (already bool) {
	count, _ := mdb.Count(
		models.GallTable, bson.M{
			"boardid":   bson.M{"$regex": bson.RegEx{Pattern: "^" + regexp.QuoteMeta(boardID) + "$", Options: "i"}},
			"channelid": channelID,
		})
	return count > 0
}
