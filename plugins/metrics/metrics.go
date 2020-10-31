package metrics

import (
	"github.com/jinzhu/gorm"
)

type metric interface {
	Register(*gorm.DB) error
}

var (
	totalCommands         = CounterMetric{key: "total_commands"}
	totalMessagesReceived = CounterMetric{key: "total_messages_received"}
)

var metrics = []metric{
	&totalCommands,
	&totalMessagesReceived,
}
