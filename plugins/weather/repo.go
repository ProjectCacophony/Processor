package weather

import "github.com/jinzhu/gorm"

func (w *Weather) saveToDB(db *gorm.DB) error {
	if db.NewRecord(w) {
		return db.Create(w).Error
	}
	return db.Update(w).Error
}

func (p *Plugin) getUserWeather(userID string) (*Weather, error) {
	var weatherInfo Weather
	err := p.db.
		Model(Weather{}).
		Where(Weather{UserID: userID}).
		First(&weatherInfo).
		Error

	return &weatherInfo, err
}
