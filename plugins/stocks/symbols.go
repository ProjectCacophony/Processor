package stocks

import (
	"strings"

	"github.com/jinzhu/gorm"
)

func findSymbol(db *gorm.DB, symbol string) (*Symbol, error) {
	var symbolEntry Symbol
	err := db.Where("symbol = ?", symbol).First(&symbolEntry).Error
	if err != nil &&
		strings.Contains(err.Error(), "record not found") {

		err = db.Where("symbol LIKE  ? || '-%'", symbol).First(&symbolEntry).Error
		if err != nil &&
			strings.Contains(err.Error(), "record not found") {

			err = db.Where("UPPER(name) LIKE  '%' || ? || '%'", symbol).First(&symbolEntry).Error
		}
	}

	return &symbolEntry, err
}
