package weather

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Jeffail/gabs"
	"github.com/bwmarrin/discordgo"
	"github.com/shawntoffel/darksky"
	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) viewWeather(event *events.Event) {

	if len(event.Fields()) == 0 {
		return
	}

	inputAddress := strings.Join(event.Fields()[1:], " ")
	escapedAddress := url.QueryEscape(inputAddress)

	link := fmt.Sprintf(geocodeEndpoint, p.config.GoogleMapsKey, escapedAddress)

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

	if results.Path("status").Data().(string) != "OK" {
		event.Respond("weather.location.not-found")
		return
	}

	allLocationInfo, err := results.Path("results").Children()
	if err != nil || len(allLocationInfo) == 0 {
		event.Except(err)
		return
	}

	locationInfo := allLocationInfo[0]

	weather := &Weather{
		UserID:             event.UserID,
		Longitude:          locationInfo.Path("geometry.location.lng").Data().(float64),
		Latitude:           locationInfo.Path("geometry.location.lat").Data().(float64),
		UserEnteredAddress: inputAddress,
		Address:            locationInfo.Path("formatted_address").Data().(string),
		PlaceID:            locationInfo.Path("place_id").Data().(string),
	}

	forecast, err := p.darkSky.Forecast(darksky.ForecastRequest{
		Latitude:  darksky.Measurement(weather.Latitude),
		Longitude: darksky.Measurement(weather.Longitude),
		Options: darksky.ForecastRequestOptions{
			Exclude: "minutely,hourly,alerts,flags",
			Extend:  "",
			Lang:    "en",
			Units:   "si",
		},
	})
	if err != nil {
		event.Except(err)
		return
	}

	outputFormat := "weather.outputformat"
	temp1 := strconv.FormatFloat(float64(forecast.Currently.Temperature), 'f', 1, 64)
	temp2 := strconv.FormatFloat(float64(forecast.Currently.Temperature)*1.8+32, 'f', 1, 64)

	if strings.Contains(weather.Address, "USA") {
		outputFormat = "weather.outputformat.usa"
		temp1 = strconv.FormatFloat(float64(forecast.Currently.Temperature)*1.8+32, 'f', 1, 64)
		temp2 = strconv.FormatFloat(float64(forecast.Currently.Temperature), 'f', 1, 64)
	}

	var embeds []*discordgo.MessageEmbed
	currentWeatherEmbed := p.makeWeatherEmbed(event, weather)
	currentWeatherEmbed.Fields = []*discordgo.MessageEmbedField{
		{
			Name: "Currently",
			Value: fmt.Sprintf(event.Translate(outputFormat),
				p.getWeatherEmoji(forecast.Currently.Icon),
				forecast.Currently.Summary,
				temp1,
				temp2,
				strconv.FormatFloat(float64(forecast.Currently.WindSpeed), 'f', 1, 64),
				strconv.FormatFloat(float64(forecast.Currently.WindSpeed)*2.23694, 'f', 1, 64),
				strconv.FormatFloat(float64(forecast.Currently.Humidity)*100, 'f', 0, 64),
			),
			Inline: false,
		},
		{
			Name:   "This week",
			Value:  p.getWeatherEmoji(forecast.Daily.Icon) + " " + forecast.Daily.Summary,
			Inline: false,
		},
	}

	embeds = append(embeds, currentWeatherEmbed)
	embeds = append(embeds, p.makeWeatherEmbed(event, weather))
	embeds = append(embeds, p.makeWeatherEmbed(event, weather))

	for i, day := range forecast.Daily.Data {
		if i <= 2 {
			embeds[1].Fields = append(embeds[1].Fields, p.makeFieldFromDay(event, day, weather))
		} else if i >= 3 && i <= 5 {
			embeds[2].Fields = append(embeds[2].Fields, p.makeFieldFromDay(event, day, weather))
		}
	}

	err = event.Paginator().EmbedPaginator(
		event.BotUserID,
		event.ChannelID,
		event.UserID,
		event.DM(),
		embeds...,
	)
	if err != nil {
		event.Except(err)
		return
	}

}

func (p *Plugin) makeWeatherEmbed(event *events.Event, weather *Weather) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: fmt.Sprintf("Weather in %s", weather.Address),
		URL: fmt.Sprintf(darkskyForcastURL,
			strconv.FormatFloat(weather.Latitude, 'f', -1, 64),
			strconv.FormatFloat(weather.Latitude, 'f', -1, 64)),
		Footer: &discordgo.MessageEmbedFooter{
			Text:    event.Translate("weather.darkSky.mention"),
			IconURL: event.Translate("weather.darkSky.logo"),
		},
		Color: discord.HexToColorCode(darkSkyHexColor),
	}
}

func (p *Plugin) makeFieldFromDay(event *events.Event, day darksky.DataPoint, weather *Weather) *discordgo.MessageEmbedField {
	tm := time.Unix(int64(day.Time), 0)

	outputFormat := "weather.outputformat.daily"
	temp1 := strconv.FormatFloat(float64(day.TemperatureHigh), 'f', 1, 64)
	temp2 := strconv.FormatFloat(float64(day.TemperatureLow), 'f', 1, 64)
	temp3 := strconv.FormatFloat(float64(day.TemperatureHigh)*1.8+32, 'f', 1, 64)
	temp4 := strconv.FormatFloat(float64(day.TemperatureLow)*1.8+32, 'f', 1, 64)

	if strings.Contains(weather.Address, "USA") {
		outputFormat = "weather.outputformat.daily.usa"
		temp3 = strconv.FormatFloat(float64(day.TemperatureHigh), 'f', 1, 64)
		temp4 = strconv.FormatFloat(float64(day.TemperatureLow), 'f', 1, 64)
		temp1 = strconv.FormatFloat(float64(day.TemperatureHigh)*1.8+32, 'f', 1, 64)
		temp2 = strconv.FormatFloat(float64(day.TemperatureLow)*1.8+32, 'f', 1, 64)
	}

	return &discordgo.MessageEmbedField{
		Name: tm.Format("__**Monday (Jan 2)**__"),
		Value: fmt.Sprintf(event.Translate(outputFormat),
			p.getWeatherEmoji(day.Icon),
			day.Summary,
			temp1,
			temp2,
			temp3,
			temp4,
			strconv.FormatFloat(float64(day.WindSpeed), 'f', 1, 64),
			strconv.FormatFloat(float64(day.WindSpeed)*2.23694, 'f', 1, 64),
			strconv.FormatFloat(float64(day.Humidity)*100, 'f', 0, 64),
		),
		Inline: false,
	}
}

func (*Plugin) getWeatherEmoji(iconName string) (emoji string) {
	switch iconName {
	case "clear-day":
		return "☀"
	case "clear-night":
		return ""
	case "rain":
		return "🌧"
	case "snow":
		return "☃"
	case "sleet":
		return "🌃"
	case "wind":
		return "🌬"
	case "fog":
		return "🌁"
	case "cloudy":
		return "☁"
	case "partly-cloudy-day":
		return "⛅"
	case "partly-cloudy-night":
		return "☁"
	case "hail":
		return "🌨"
	case "thunderstorm":
		return "⛈"
	case "tornado":
		return "🌪"
	}
	return ""
}
