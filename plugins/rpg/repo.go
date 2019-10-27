package rpg

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/Cacophony/Processor/plugins/rpg/models"
)

type repo struct {
	db *gorm.DB
}

func (r *repo) StoreInteraction(interaction *models.Interaction) error {
	return r.db.Save(interaction).Error
}
