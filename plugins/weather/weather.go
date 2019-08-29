package weather

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"

	"github.com/Jeffail/gabs"
	"github.com/davecgh/go-spew/spew"
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) viewWeather(event *events.Event) {

	if len(event.Fields()) == 0 {
		return
	}

	inputAddress := url.QueryEscape(strings.Join(event.Fields()[1:], " "))

	link := fmt.Sprintf(geocodeEndpoint, p.config.GoogleMapsKey, inputAddress)

	fmt.Println("----------------")
	fmt.Println(link)
	fmt.Println("----------------")

	resp, err := event.HTTPClient().Get(link)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		event.Except(err)
		return
	}

	results, err := gabs.ParseJSON(bytes)
	if err != nil {
		event.Except(err)
		return
	}

	locationInfo, err := results.Path("results").Children()
	if err != nil {
		event.Except(err)
		return
	}

	spew.Dump(locationInfo)
}
