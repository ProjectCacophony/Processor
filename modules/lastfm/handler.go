package lastfm

import (
	"strings"

	"context"

	"github.com/sirupsen/logrus"
	"gitlab.com/Cacophony/dhelpers"
	"gitlab.com/Cacophony/dhelpers/cache"
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
func (m *Module) Action(ctx context.Context, event dhelpers.EventContainer) {
	switch event.Type {
	case dhelpers.MessageCreateEventType:

		for _, destination := range event.Destinations {

			switch destination.Name {
			case "lastfm":
				if len(event.Args) < 2 { // [p]lastfm|lf

					displayAbout(ctx, event)
					return
				}

				switch strings.ToLower(event.Args[1]) {
				case "set", "register", "save": // [p]lastfm|lf set|register|save <last.fm username>

					setUsername(ctx, event)
					return
				case "np", "nowplaying", "now": // [p]lastfm|lf np|nowplaying|now [<@user or user id or lastfm username>]

					displayNowPlaying(ctx, event)
					return
				case "recent", "recently", "last", "recents": // [p]lastfm|lf recent|recently|last|recents [<@user or user id or lastfm username>]

					displayRecent(ctx, event)
					return
				case "topartists", "topartist", "top-artist", "top-artists", "artist", "artists", "ta": // [p]lastfm|lf topartists|topartist|top-artist|top-artists|artist|artists|ta [<@user or user id or lastfm username>] [<timerange>] [<collage>]

					displayTopArtists(ctx, event)
					return
				case "toptracks", "toptrack", "top-track", "top-tracks", "track", "tracks", "tt", "ts": // [p]lastfm|lf toptracks|toptrack|top-track|top-tracks|track|tracks|tt|ts [<@user or user id or lastfm username>] [<timerange>] [<collage>] [<server>]

					if serverRequest, _ := isServerRequest(event.Args); serverRequest {
						displayServerTopTracks(ctx, event)
						return
					}

					displayTopTracks(ctx, event)
					return
				case "topalbums", "topalbum", "top-album", "top-albums", "album", "albums", "tal": // [p]lastfm|lf topalbums|topalbum|top-album|top-albums|album|albums|tal [<@user or user id or lastfm username>] [<timerange>] [<collage>]

					displayTopAlbums(ctx, event)
					return
				case "server-top", "server-toptracks": // [p]lastfm|lf server-top|server-toptracks [<timerange>]
					displayServerTopTracks(ctx, event)
					return
				case "top":
					if len(event.Args) < 3 {
						return
					}

					switch strings.ToLower(event.Args[2]) {
					case "artist", "artists": // [p]lastfm|lf top artist|artists [<@user or user id or lastfm username>] [<timerange>] [<period>]
						event.Args = append(event.Args[:1], event.Args[1+1:]...)

						displayTopArtists(ctx, event)
						return
					case "track", "tracks": // [p]lastfm|lf top track|tracks [<@user or user id or lastfm username>] [<timerange>] [<period>] [<server>]
						event.Args = append(event.Args[:1], event.Args[1+1:]...)

						if serverRequest, _ := isServerRequest(event.Args); serverRequest {
							displayServerTopTracks(ctx, event)
							return
						}

						displayTopTracks(ctx, event)
						return
					case "album", "albums": // [p]lastfm|lf top album|albums [<@user or user id or lastfm username>] [<timerange>] [<period>]
						event.Args = append(event.Args[:1], event.Args[1+1:]...)

						displayTopAlbums(ctx, event)
						return
					}

				default: // [p]lastfm|lf <@user or user id or lastfm username>

					displayAbout(ctx, event)
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
