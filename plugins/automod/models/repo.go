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

func GetRule(
	db *gorm.DB,
	id uint,
) (*Rule, error) {
	if id == 0 {
		return nil, errors.New("id cannot be empty")
	}

	var rule Rule
	err := db.
		Preload("Filters").
		Preload("Actions").
		Where("id = ?", id).
		First(&rule).Error
	return &rule, err
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

	oldRule, err := GetRule(db, id)
	if err != nil {
		return err
	}

	rule.ID = oldRule.ID
	for i := range oldRule.Filters {
		if len(rule.Filters) >= i+1 {
			rule.Filters[i].ID = oldRule.Filters[i].ID
			continue
		}

		err = db.Delete(RuleFilter{}, "id = ?", oldRule.Filters[i].ID).Error
		if err != nil {
			return err
		}
	}
	for i := range oldRule.Actions {
		if len(rule.Actions) >= i+1 {
			rule.Actions[i].ID = oldRule.Actions[i].ID
			continue
		}

		err = db.Delete(RuleAction{}, "id = ?", oldRule.Actions[i].ID).Error
		if err != nil {
			return err
		}
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
