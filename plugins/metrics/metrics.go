package metrics

import (
	"github.com/jinzhu/gorm"
)

type metric interface {
	Register(*gorm.DB) error
}

var totalCommands = CounterMetric{key: "total_commands"}

var metrics = []metric{
	&totalCommands,
}
