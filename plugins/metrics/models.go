package metrics

import (
	"github.com/jinzhu/gorm"
)

type Counter struct {
	Key   string `gorm:"unique_index"`
	Value int
}

func (*Counter) TableName() string {
	return "metrics_counters"
}

type CounterMetric struct {
	key string
}

func (m *CounterMetric) Register(db *gorm.DB) error {
	var count int
	err := db.Model(&Counter{}).Where("key = ?", m.key).Count(&count).Error
	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	return db.Create(&Counter{
		Key:   m.key,
		Value: 0,
	}).Error
}

func (m *CounterMetric) Inc(db *gorm.DB) error {
	return db.Model(&Counter{}).Where("key = ?", m.key).Update("value", gorm.Expr("value + 1")).Error
}
