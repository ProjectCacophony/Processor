package lastfm

import (
	"strings"

	"gitlab.com/Cacophony/Processor/plugins/common"

	"github.com/Seklfreak/lastfm-go/lastfm"
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

	p.lastfmClient = lastfm.New(
		"57f55283a6b3d6e65c10192186871cba",
		"46a19473b0482b854e32ada1032e62b6",
	) // TODO: don't store plaintest

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
		// TODO: send message
		return true
	}

	switch strings.ToLower(event.Fields()[1]) {
	case "np", "nowplaying", "now":

		displayNowPlaying(event, p.lastfmClient)
		return true

	case "set":

		setUsername(event)
		return true
	}

	return false
}
