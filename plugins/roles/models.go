package roles

import (
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"gitlab.com/Cacophony/go-kit/state"
)

const (
	confirmCategoryDeleteKey = "cacophony:processor:role:questionnaire:delete-category"
)

type Category struct {
	gorm.Model

	Name      string
	GuildID   string
	ChannelID string // the channel this Category will listen to for role assignments
	Roles     []Role `gorm:"foreignkey:CategoryID"`
	Enabled   bool
	Hidden    bool
	Limit     int // 0 = no limit
}

type Role struct {
	gorm.Model

	CategoryID   uint // roles aren't guarenteed to have a category
	GuildID      string
	ServerRoleID string
	PrintName    string
	Enabled      bool
	Aliases      pq.StringArray `gorm:"type:varchar[]"`
	// Reactions []string
}

func (*Category) TableName() string {
	return "role_categories"
}

func (*Role) TableName() string {
	return "roles"
}

func (r *Role) Name(state *state.State) string {
	if r.PrintName != "" {
		return r.PrintName
	}

	serverRole, err := state.Role(r.GuildID, r.ServerRoleID)
	if err != nil {
		return ""
	}

	return serverRole.Name
}

// matches the input to see if its the server role name or the print name, or one of the aliases
func (r *Role) Match(state *state.State, input string) bool {
	if input == "" {
		return false
	}

	if r.PrintName == input {
		return true
	}

	for _, alias := range r.Aliases {
		if alias == input {
			return true
		}
	}

	serverRole, err := state.Role(r.GuildID, r.ServerRoleID)
	if err != nil {
		return false
	}

	if serverRole.Name == input {
		return true
	}

	return false
}
