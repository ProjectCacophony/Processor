package models

import (
	"time"

	"github.com/globalsign/mgo/bson"
	"gitlab.com/Cacophony/dhelpers/mdb"
)

const (
	// GallTable is the MongoDB Collection for GallFeedEntry entries
	GallTable mdb.Collection = "gall"
)

// GallFeedEntry is an entry for each gall feed set up
type GallFeedEntry struct {
	ID            bson.ObjectId `bson:"_id,omitempty"`
	GuildID       string
	ChannelID     string
	AddedByUserID string
	BoardID       string
	MinorGallery  bool
	Recommended   bool
	LastCheck     time.Time
	AddedAt       time.Time
	PostedPostIDs []string
}
