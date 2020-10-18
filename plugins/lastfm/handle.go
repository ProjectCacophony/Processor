package lastfm

import (
	"errors"
	"strings"

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
}

func (p *Plugin) Names() []string {
	return []string{"lastfm", "lf", "fm"}
}

func (p *Plugin) Start(params common.StartParameters) error {
	p.db = params.DB
	p.logger = params.Logger

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
		Names:       p.Names(),
		Description: "lastfm.help.description",
		Commands: []common.Command{
			{
				Name:        "lastfm.help.lf.name",
				Description: "lastfm.help.lf.description",
			},
			{
				Name:        "lastfm.help.np.name",
				Description: "lastfm.help.np.description",
				Params: []common.CommandParam{
					{Name: "np", Type: common.Flag},
					{Name: "Last.FM Username", Type: common.Text},
				},
			},
			{
				Name:        "lastfm.help.ta.name",
				Description: "lastfm.help.ta.description",
				Params: []common.CommandParam{
					{Name: "ta", Type: common.Flag},
					{Name: "week|month|quarter|6month|year|overall (default)", Type: common.Text, Optional: true},
					{Name: "Last.FM Username", Type: common.Text, Optional: true},
				},
			},
			{
				Name:        "lastfm.help.tt.name",
				Description: "lastfm.help.tt.description",
				Params: []common.CommandParam{
					{Name: "tt", Type: common.Flag},
					{Name: "week|month|quarter|6month|year|overall (default)", Type: common.Text, Optional: true},
					{Name: "Last.FM Username", Type: common.Text, Optional: true},
				},
			},
			{
				Name:        "lastfm.help.tal.name",
				Description: "lastfm.help.tal.description",
				Params: []common.CommandParam{
					{Name: "tal", Type: common.Flag},
					{Name: "week|month|quarter|6month|year|overall (default)", Type: common.Text, Optional: true},
					{Name: "Last.FM Username", Type: common.Text, Optional: true},
					{Name: "collage", Type: common.Flag, Optional: true},
				},
			},
			{
				Name:        "lastfm.help.recent.name",
				Description: "lastfm.help.recent.description",
				Params: []common.CommandParam{
					{Name: "recent", Type: common.Flag},
					{Name: "Last.FM Username", Type: common.Text, Optional: true},
				},
			},
			{
				Name:        "lastfm.help.about.name",
				Description: "lastfm.help.about.description",
				Params: []common.CommandParam{
					{Name: "about", Type: common.Flag},
					{Name: "Last.FM Username", Type: common.Text, Optional: true},
				},
			},
			{
				Name:        "lastfm.help.set.name",
				Description: "lastfm.help.set.description",
				Params: []common.CommandParam{
					{Name: "set", Type: common.Flag},
					{Name: "Last.FM Username", Type: common.Text},
				},
			},
		},
	}
}

func (p *Plugin) Action(event *events.Event) bool {
	if !event.Command() {
		return false
	}

	if event.Fields()[0] != "lastfm" &&
		event.Fields()[0] != "lf" &&
		event.Fields()[0] != "fm" {
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
