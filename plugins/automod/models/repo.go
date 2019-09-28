package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

func CreateRule(
	db *gorm.DB,
	rule *Rule,
) error {
	if rule == nil {
		return errors.New("rule cannot be empty")
	}

	return db.Create(rule).Error
}

func UpdateRule(
	db *gorm.DB,
	id uint,
	rule *Rule,
) error {
	if id == 0 {
		return errors.New("id cannot be empty")
	}
	if rule == nil {
		return errors.New("rule cannot be empty")
	}

	return db.Model(Rule{}).Where("id = ?", id).Save(rule).Error
}

func DeleteRule(
	db *gorm.DB,
	rule *Rule,
) error {
	if rule == nil || rule.ID == 0 {
		return errors.New("rule cannot be empty")
	}

	return db.Delete(rule).Error
}
