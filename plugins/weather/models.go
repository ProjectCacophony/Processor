package weather

import "github.com/jinzhu/gorm"

type Weather struct {
	gorm.Model

	Longitude float64
	Latitude  float64

	UserEnteredAddress string
	Address            string
	UserID             string
	PlaceID            string
}

func (*Weather) TableName() string {
	return "weather"
}
