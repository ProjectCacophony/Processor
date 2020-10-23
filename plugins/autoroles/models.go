package autoroles

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
)

type AutoRole struct {
	gorm.Model

	GuildID      string
	ServerRoleID string
	RuleID       uint
	Rule         models.Rule
}

func (*AutoRole) TableName() string {
	return "auto_role"
}
