package models

import (
	"time"

	"github.com/globalsign/mgo/bson"
	"gitlab.com/Cacophony/dhelpers/mdb"
)

const (
	// FeedTable is the MongoDB Collection for FeedEntry entries
	FeedTable mdb.Collection = "feeds"
)

// FeedEntry is an entry for each feed set up
type FeedEntry struct {
	ID            bson.ObjectId `bson:"_id,omitempty"`
	GuildID       string
	ChannelID     string
	AddedByUserID string
	FeedURL       string
	FeedTitle     string
	LastCheck     time.Time
	AddedAt       time.Time
	PostedPostIDs []string
}
