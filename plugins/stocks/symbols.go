package stocks

import (
	"strings"

	"github.com/jinzhu/gorm"
)

const priorityRegion = "US"

var priorityOrderExpr = gorm.Expr(
	"case when region = ? then 1 else 2 end ASC", priorityRegion,
)

func findSymbol(db *gorm.DB, symbol string) (*Symbol, error) {
	var symbolEntry Symbol

	// try to match exact symbol
	err := db.
		Where("symbol = ?", symbol).
		Order(priorityOrderExpr).
		First(&symbolEntry).
		Error
	if err != nil &&
		strings.Contains(err.Error(), "record not found") {

		// try to match symbol with undefined country suffix
		err = db.
			Where("symbol LIKE  ? || '-%'", symbol).
			Order(priorityOrderExpr).
			First(&symbolEntry).
			Error
		if err != nil &&
			strings.Contains(err.Error(), "record not found") {

			// try to match company name
			err = db.
				Where("UPPER(name) LIKE  '%' || ? || '%'", symbol).
				Order(priorityOrderExpr).
				First(&symbolEntry).
				Error
		}
	}

	return &symbolEntry, err
}
