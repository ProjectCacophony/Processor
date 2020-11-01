package metrics

import (
	"github.com/jinzhu/gorm"
)

type metric interface {
	Key() string
	Register(*gorm.DB) error
	Get(*gorm.DB) (int, error)
}

var (
	totalCommands         = CounterMetric{key: "total_commands"}
	totalMessagesReceived = CounterMetric{key: "total_messages_received"}
)

var metrics = []metric{
	&totalCommands,
	&totalMessagesReceived,
}
