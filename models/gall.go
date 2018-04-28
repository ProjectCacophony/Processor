package models

import (
	"time"

	"github.com/globalsign/mgo/bson"
	"gitlab.com/project-d-collab/dhelpers/mdb"
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
	LastCheck     time.Time
}
