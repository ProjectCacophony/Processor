package models

import (
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"gitlab.com/Cacophony/dhelpers/mongo"
)

const (
	// LastFmTable is the MongoDB Collection for LastFmEntry entries
	LastFmTable mongo.Collection = "lastfm"
)

var (
	// LastFmRepository contains the database logic for the table
	LastFmRepository = mongo.NewRepository(LastFmTable)
)

// LastFmEntry is an entry for each Discord User
type LastFmEntry struct {
	ID             *objectid.ObjectID `bson:"_id,omitempty"`
	UserID         string
	LastFmUsername string
}
