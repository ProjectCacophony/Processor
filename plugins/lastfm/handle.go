package lastfm

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/jinzhu/gorm"

	"github.com/Seklfreak/lastfm-go/lastfm"
	"github.com/kelseyhightower/envconfig"
	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
)

type Plugin struct {
	logger       *zap.Logger
	db           *gorm.DB
	lastfmClient *lastfm.Api
	httpClient   *http.Client
}

func (p *Plugin) Name() string {
	return "lastfm"
}

func (p *Plugin) Start(params common.StartParameters) error {
	err := params.DB.AutoMigrate(User{}).Error
	if err != nil {
		return err
	}

	var config Config
	err = envconfig.Process("", &config)
	if err != nil {
		return nil
	}

	if config.Key == "" || config.Secret == "" {
		return errors.New("last.fm plugin configuration missing")
	}

	p.lastfmClient = lastfm.New(
		config.Key,
		config.Secret,
	)

	p.db = params.DB
	p.logger = params.Logger

	// TODO: global http client?
	p.httpClient = &http.Client{
		Timeout: time.Second * 10,
	}

	return nil
}

func (p *Plugin) Stop(params common.StopParameters) error {
	return nil
}

func (p *Plugin) Priority() int {
	return 0
}

func (p *Plugin) Passthrough() bool {
	return false
}

func (p *Plugin) Help() *common.PluginHelp {
	return &common.PluginHelp{
		Name:        p.Name(),
		Description: "help.lastfm.description",
	}
}

func (p *Plugin) Action(event *events.Event) bool {
	if !event.Command() {
		return false
	}

	if event.Fields()[0] != "lastfm" &&
		event.Fields()[0] != "lf" {
		return false
	}

	event.Typing()

	if len(event.Fields()) < 2 {
		p.handleNowPlaying(event, p.lastfmClient, 1)
		return true
	}

	switch strings.ToLower(event.Fields()[1]) {
	case "np", "nowplaying", "now":

		p.handleNowPlaying(event, p.lastfmClient, 2)
		return true

	case "topartists", "topartist", "top-artist", "top-artists", "artist", "artists", "ta":

		p.handleTopArtists(event, p.lastfmClient, 2)
		return true

	case "toptracks", "toptrack", "top-track", "top-tracks", "track", "tracks", "tt", "ts":

		p.handleTopTracks(event, p.lastfmClient, 2)
		return true

	case "topalbums", "topalbum", "top-album", "top-albums", "album", "albums", "tal":

		p.handleTopAlbums(event, p.lastfmClient, 2)
		return true

	case "top":
		if len(event.Fields()) < 3 {
			event.Respond("lastfm.no-subcommand")
			return true
		}

		switch strings.ToLower(event.Fields()[2]) {
		case "artist", "artists":

			p.handleTopArtists(event, p.lastfmClient, 3)
			return true

		case "track", "tracks":

			p.handleTopTracks(event, p.lastfmClient, 3)
			return true

		case "album", "albums":

			p.handleTopAlbums(event, p.lastfmClient, 3)
			return true

		}

	case "recent", "recently", "last", "recents":

		p.handleRecent(event, p.lastfmClient)
		return true

	case "set", "register", "save":

		p.handleSet(event)
		return true

	case "about", "user", "u":

		p.handleAbout(event, p.lastfmClient)
		return true
	}

	p.handleNowPlaying(event, p.lastfmClient, 1)
	return true
}
