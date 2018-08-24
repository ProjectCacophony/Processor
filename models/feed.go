package models

import (
	"time"

	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"gitlab.com/Cacophony/dhelpers/mongo"
)

const (
	// FeedTable is the MongoDB Collection for FeedEntry entries
	FeedTable mongo.Collection = "feeds"
)

var (
	// FeedRepository contains the database logic for the table
	FeedRepository = mongo.NewRepository(FeedTable)
)

// FeedEntry is an entry for each feed set up
type FeedEntry struct {
	ID            *objectid.ObjectID `bson:"_id,omitempty"`
	GuildID       string
	ChannelID     string
	AddedByUserID string
	FeedURL       string
	FeedTitle     string
	LastCheck     time.Time
	AddedAt       time.Time
	PostedPostIDs []string
}
