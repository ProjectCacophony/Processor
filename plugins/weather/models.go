package weather

import "github.com/jinzhu/gorm"

type Weather struct {
	gorm.Model

	Longitude   float32
	Latitude    float32
	AddressText string
	UserId      string
}

func (*Weather) TableName() string {
	return "weather"
}
