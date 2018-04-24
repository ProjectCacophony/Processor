package lastfm

import (
	"strings"

	"github.com/sirupsen/logrus"
	"gitlab.com/project-d-collab/dhelpers"
	"gitlab.com/project-d-collab/dhelpers/cache"
)

// Module is a struct for the module
type Module struct{}

// GetDestinations defines valid destinations for the module
func (m *Module) GetDestinations() []string {
	return []string{
		"lastfm",
	}
}

// GetTranslationFiles defines all translation files for the module
func (m *Module) GetTranslationFiles() []string {
	return []string{
		"lastfm.en.toml",
	}
}

// Action is the module entry point when event is triggered
func (m *Module) Action(event dhelpers.EventContainer) {

	switch event.Type {
	case dhelpers.MessageCreateEventType:

		for _, destination := range event.Destinations {

			switch destination.Name {
			case "lastfm":
				if len(event.Args) < 2 { // [p]lastfm|lf

					displayAbout(event)
					return
				}

				switch strings.ToLower(event.Args[1]) {
				case "set", "register", "save": // [p]lastfm|lf set|register|save <last.fm username>

					setUsername(event)
					return
				case "np", "nowplaying", "now": // [p]lastfm|lf np|nowplaying|now [<@user or user id or lastfm username>]

					displayNowPlaying(event)
					return
				case "recent", "recently", "last", "recents": // [p]lastfm|lf recent|recently|last|recents [<@user or user id or lastfm username>]

					displayRecent(event)
					return
				case "topartists", "topartist", "top-artist", "top-artists", "artist", "artists": // [p]lastfm|lf topartists|topartist|top-artist|top-artists|artist|artists [<@user or user id or lastfm username>] [<timerange>] [<period>]

					displayTopArtists(event)
					return
				case "toptracks", "toptrack", "top-track", "top-tracks", "track", "tracks": // [p]lastfm|lf toptracks|toptrack|top-track|top-tracks|track|tracks [<@user or user id or lastfm username>] [<timerange>] [<period>]

					displayTopTracks(event)
					return
				case "topalbums", "topalbum", "top-album", "top-albums", "album", "albums": // [p]lastfm|lf topalbums|topalbum|top-album|top-albums|album|albums [<@user or user id or lastfm username>] [<timerange>] [<period>]

					displayTopAlbums(event)
					return
				case "top":
					if len(event.Args) < 3 {
						return
					}

					switch strings.ToLower(event.Args[2]) {
					case "artist", "artists": // [p]lastfm|lf top artist|artists [<@user or user id or lastfm username>] [<timerange>] [<period>]

						displayTopArtists(event)
						return
					case "track", "tracks": // [p]lastfm|lf top track|tracks [<@user or user id or lastfm username>] [<timerange>] [<period>]

						displayTopTracks(event)
						return
					case "album", "albums": // [p]lastfm|lf top album|albums [<@user or user id or lastfm username>] [<timerange>] [<period>]

						displayTopAlbums(event)
						return
					}

				default: // [p]lastfm|lf <@user or user id or lastfm username>

					displayAbout(event)
					return
				}
			}
		}
	}
}

// Init is called on bot startup
func (m *Module) Init() {
}

// Uninit is called on normal bot shutdown
func (m *Module) Uninit() {
}

func logger() *logrus.Entry {
	return cache.GetLogger().WithField("module", "lastfm")
}
