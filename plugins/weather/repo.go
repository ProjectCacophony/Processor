package weather

func (p *Plugin) getUserWeather(userID string) (*Weather, error) {
	var weatherInfo Weather
	err := p.db.
		Model(Weather{}).
		Where(Weather{UserID: userID}).
		First(&weatherInfo).
		Error

	return &weatherInfo, err
}
