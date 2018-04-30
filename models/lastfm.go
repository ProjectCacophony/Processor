package models

import (
	"github.com/globalsign/mgo/bson"
	"gitlab.com/Cacophony/dhelpers/mdb"
)

const (
	// LastFmTable is the MongoDB Collection for LastFmEntry entries
	LastFmTable mdb.Collection = "lastfm"
)

// LastFmEntry is an entry for each Discord User
type LastFmEntry struct {
	ID             bson.ObjectId `bson:"_id,omitempty"`
	UserID         string
	LastFmUsername string
}
