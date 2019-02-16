package lastfm

import (
	"errors"
	"strings"

	"github.com/Seklfreak/lastfm-go/lastfm"
	"github.com/kelseyhightower/envconfig"
	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/interfaces"
	"gitlab.com/Cacophony/go-kit/localisation"
)

type Plugin struct {
	lastfmClient *lastfm.Api
}

func (p *Plugin) Name() string {
	return "lastfm"
}

func (p *Plugin) Start(params common.StartParameters) error {
	params.DB.AutoMigrate(User{})

	var config Config
	err := envconfig.Process("", &config)
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

func (p *Plugin) Localisations() []interfaces.Localisation {
	local, err := localisation.NewFileSource("assets/translations/lastfm.en.toml", "en")
	if err != nil {
		panic(err) // TODO: handle error
	}

	return []interfaces.Localisation{local}
}

func (p *Plugin) Action(event *events.Event) bool {
	if !event.Command() {
		return false
	}

	if event.Fields()[0] != "lastfm" &&
		event.Fields()[0] != "lf" {
		return false
	}

	if len(event.Fields()) < 2 {
		event.Respond("lastfm.no-subcommand") // nolint: errcheck
		return true
	}

	event.Typing()

	switch strings.ToLower(event.Fields()[1]) {
	case "np", "nowplaying", "now":

		handleNowPlaying(event, p.lastfmClient)
		return true

	case "set":

		handleSet(event)
		return true
	}

	event.Respond("lastfm.no-subcommand") // nolint: errcheck
	return true
}
