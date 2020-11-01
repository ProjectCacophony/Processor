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

func (m *CounterMetric) Key() string {
	return m.key
}

func (m *CounterMetric) Register(db *gorm.DB) error {
	var count int
	err := db.Model(&Counter{}).Where("key = ?", m.Key()).Count(&count).Error
	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	return db.Create(&Counter{
		Key:   m.Key(),
		Value: 0,
	}).Error
}

func (m *CounterMetric) Inc(db *gorm.DB) error {
	return db.Model(&Counter{}).Where("key = ?", m.Key()).Update("value", gorm.Expr("value + 1")).Error
}

func (m *CounterMetric) Get(db *gorm.DB) (int, error) {
	var metricCounter struct {
		Value int
	}

	err := db.
		Table("metrics_counters").
		Where("key = ?", m.Key()).
		Select("value").
		First(&metricCounter).
		Error
	if err != nil {
		return 0, err
	}

	return metricCounter.Value, nil
}
