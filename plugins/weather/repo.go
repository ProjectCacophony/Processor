package weather

func (p *Plugin) getUserWeather(userID string) *Weather {
	var weatherInfo Weather
	p.db.
		Model(Weather{}).
		Where(Weather{UserID: userID}).
		First(&weatherInfo)

	return &weatherInfo
}
