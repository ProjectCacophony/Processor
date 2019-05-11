package stocks

import (
	"github.com/jinzhu/gorm"
)

func findSymbol(db *gorm.DB, symbol string) (*Symbol, error) {
	var symbolEntry Symbol
	err := db.Where("symbol = ?", symbol).First(&symbolEntry).Error

	return &symbolEntry, err
}
