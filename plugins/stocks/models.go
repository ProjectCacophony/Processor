package stocks

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

// Symbol model maintained by Worker
type Symbol struct {
	gorm.Model
	Symbol   string
	Exchange string
	Name     string
	Date     time.Time
	Type     string
	IEXID    string
	Region   string
	Currency string
}

func (*Symbol) TableName() string {
	return "stocks_symbols"
}

func (s *Symbol) FormatCurrency(myValue float64) string {
	currencyFormat := map[string]string{
		"USD": "$ %.2f",
		"EUR": "%.2f €",
	}

	if currencyFormat[s.Currency] != "" {
		return fmt.Sprintf(currencyFormat[s.Currency], myValue)
	}

	return fmt.Sprintf("%.2f %s", myValue, s.Currency)
}
