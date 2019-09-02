package weather

import (
	"strings"

	"github.com/jinzhu/gorm"
)

type Weather struct {
	gorm.Model

	Longitude float64
	Latitude  float64

	UserEnteredAddress string
	Address            string
	UserID             string
	PlaceID            string
}

func (w *Weather) USA() bool {
	return strings.Contains(w.Address, "USA")
}

func (*Weather) TableName() string {
	return "weather"
}
