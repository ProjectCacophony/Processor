package filters

import (
	"errors"
	"strings"
)

type AmountComparisonType string

const (
	AmountComparisonLT AmountComparisonType = "lt"
	AmountComparisonGT AmountComparisonType = "gt"
	AmountComparisonEQ AmountComparisonType = "eq"
)

func extractAmountComparisonType(field string) (AmountComparisonType, error) {
	switch strings.ToLower(field) {
	case "lt", "<":
		return AmountComparisonLT, nil
	case "gt", ">":
		return AmountComparisonGT, nil
	case "eq", "=":
		return AmountComparisonEQ, nil
	}

	return "", errors.New("invalid Comparison type, try \"lt\", \"eq\", or \"gt\"")
}
