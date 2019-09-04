package weather

import (
	"encoding/json"
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
	var weather *Weather

	// check if they passed a location, otherwise check if they have one saved
	if len(event.Fields()) > 1 {
		weather = p.getWeatherInfo(event, strings.Join(event.Fields()[1:], " "))
	} else {
		weather = p.getUserWeather(event.UserID)
		if weather == nil || weather.UserID == "" {
			event.Respond("weather.nosaved")
			return
		}
	}

	if weather == nil {
		event.Respond("weather.location.not-found")
		return
	}

	var forecast darksky.ForecastResponse

	// check cache
	cacheKey := fmt.Sprintf(placeKey, weather.PlaceID)
	forecastBytes, err := p.redis.Get(cacheKey).Bytes()
	if err != nil || len(forecastBytes) == 0 {
		forecast, err = p.darkSky.Forecast(darksky.ForecastRequest{
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

		forcastData, err := json.Marshal(forecast)
		if err != nil {
			event.Except(err)
			return
		}

		p.redis.Set(cacheKey, forcastData, placeCacheDuration)
	} else {
		err = json.Unmarshal(forecastBytes, &forecast)
		if err != nil {
			event.Except(err)
			return
		}
	}

	forecastLocation, err := time.LoadLocation(forecast.Timezone)
	if err != nil {
		forecastLocation, err = time.LoadLocation("UTC")
		if err != nil {
			event.Except(err)
			return
		}
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
	}

	var threeDays []darksky.DataPoint

	embeds = append(embeds, currentWeatherEmbed)
	if len(forecast.Daily.Data) > 3 {
		embeds = append(embeds, p.makeWeatherEmbed(event, weather))
		embeds = append(embeds, p.makeWeatherEmbed(event, weather))

		todayTime := forecast.Currently.Time
		var pastToday bool
		for _, day := range forecast.Daily.Data {
			todayDate := time.Unix(int64(todayTime), 0).In(forecastLocation).Format("02/01/06")
			dayDate := time.Unix(int64(day.Time), 0).In(forecastLocation).Format("02/01/06")

			// dark sky apy does not return daily data in a reliable order. need to loop
			// through daily info to find the current day, only get days after today
			if !pastToday {
				if todayDate == dayDate {
					pastToday = true
					threeDays = append(threeDays, day)
				}
				continue
			}

			if len(threeDays) < 3 {
				threeDays = append(threeDays, day)
			}

			if len(embeds[1].Fields) < 3 {
				embeds[1].Fields = append(embeds[1].Fields, p.makeFieldFromDay(event, day, weather, forecastLocation))
			} else if len(embeds[2].Fields) < 3 {
				embeds[2].Fields = append(embeds[2].Fields, p.makeFieldFromDay(event, day, weather, forecastLocation))
			}
		}
	}

	currentWeatherEmbed.Fields = append(currentWeatherEmbed.Fields, &discordgo.MessageEmbedField{
		Name: "This week",
		Value: event.Translate(
			"weather.current.week-summary",
			"emoji", p.getWeatherEmoji,
			"summaryIcon", forecast.Daily.Icon,
			"summaryText", forecast.Daily.Summary,
			"threeDays", threeDays,
			"usa", weather.USA(),
			"f", func(i float64) float64 {
				return i*1.8 + 32
			},
		),
		Inline: false,
	})

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

func (p *Plugin) setUserWeather(event *events.Event) {
	if len(event.Fields()) < 3 {
		event.Respond("common.to-few-params")
		return
	}

	inputAddress := strings.Join(event.Fields()[2:], " ")
	newInfo := p.getWeatherInfo(event, inputAddress)
	if newInfo == nil {
		event.Respond("weather.location.not-found")
		return
	}

	var err error
	currentInfo := p.getUserWeather(event.UserID)
	if currentInfo != nil && currentInfo.UserID != "" {
		currentInfo.Longitude = newInfo.Longitude
		currentInfo.Latitude = newInfo.Latitude
		currentInfo.Address = newInfo.Address
		currentInfo.UserEnteredAddress = newInfo.UserEnteredAddress
		currentInfo.PlaceID = newInfo.PlaceID

		err = p.db.Save(currentInfo).Error
	} else {
		err = p.db.Save(newInfo).Error
	}
	if err != nil {
		event.Except(err)
		return
	}

	event.Respond("weather.location.saved")
}

func (p *Plugin) getWeatherInfo(event *events.Event, inputAddress string) *Weather {
	if inputAddress == "" {
		return nil
	}

	escapedAddress := url.QueryEscape(inputAddress)
	link := fmt.Sprintf(geocodeEndpoint, p.config.GoogleMapsKey, escapedAddress)

	resp, err := event.HTTPClient().Get(link)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		event.Except(err)
		return nil
	}

	results, err := gabs.ParseJSON(bytes)
	if err != nil {
		event.Except(err)
		return nil
	}

	if results.Path("status").Data().(string) != "OK" {
		return nil
	}

	allLocationInfo, err := results.Path("results").Children()
	if err != nil || len(allLocationInfo) == 0 {
		return nil
	}

	locationInfo := allLocationInfo[0]
	return &Weather{
		UserID:             event.UserID,
		Longitude:          locationInfo.Path("geometry.location.lng").Data().(float64),
		Latitude:           locationInfo.Path("geometry.location.lat").Data().(float64),
		UserEnteredAddress: inputAddress,
		Address:            locationInfo.Path("formatted_address").Data().(string),
		PlaceID:            locationInfo.Path("place_id").Data().(string),
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

func (p *Plugin) makeFieldFromDay(event *events.Event, day darksky.DataPoint, weather *Weather, loc *time.Location) *discordgo.MessageEmbedField {
	tm := time.Unix(int64(day.Time), 0).In(loc)

	outputFormat := "weather.outputformat.daily"
	temp1 := strconv.FormatFloat(float64(day.TemperatureHigh), 'f', 1, 64)
	temp2 := strconv.FormatFloat(float64(day.TemperatureLow), 'f', 1, 64)
	temp3 := strconv.FormatFloat(float64(day.TemperatureHigh)*1.8+32, 'f', 1, 64)
	temp4 := strconv.FormatFloat(float64(day.TemperatureLow)*1.8+32, 'f', 1, 64)

	if weather.USA() {
		outputFormat = "weather.outputformat.daily.usa"
		temp3 = strconv.FormatFloat(float64(day.TemperatureHigh), 'f', 1, 64)
		temp4 = strconv.FormatFloat(float64(day.TemperatureLow), 'f', 1, 64)
		temp1 = strconv.FormatFloat(float64(day.TemperatureHigh)*1.8+32, 'f', 1, 64)
		temp2 = strconv.FormatFloat(float64(day.TemperatureLow)*1.8+32, 'f', 1, 64)
	}

	return &discordgo.MessageEmbedField{
		Name: tm.Format("**Monday (Jan 2)**"),
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
