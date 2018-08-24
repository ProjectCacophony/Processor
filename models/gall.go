package models

import (
	"time"

	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"gitlab.com/Cacophony/dhelpers/mongo"
)

const (
	// GallTable is the MongoDB Collection for GallFeedEntry entries
	GallTable mongo.Collection = "gall"
)

var (
	// GallRepository contains the database logic for the table
	GallRepository = mongo.NewRepository(GallTable)
)

// GallFeedEntry is an entry for each gall feed set up
type GallFeedEntry struct {
	ID            *objectid.ObjectID `bson:"_id,omitempty"`
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
