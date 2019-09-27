package roles

import "strings"

func (p *Plugin) getAllCategories(guildID string) ([]*Category, error) {
	var categories []*Category
	err := p.db.
		Preload("Roles").
		Find(&categories, "guild_id = ?", guildID).
		Error
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		return nil, err
	}
	return categories, nil
}

func (p *Plugin) getCategoryByName(name string, guildID string) (*Category, error) {
	var category Category
	err := p.db.First(&category, "name = ? and guild_id = ?", name, guildID).Error
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		return nil, err
	}
	return &category, nil
}

func (p *Plugin) getRoleByServerRoleID(serverRoleID string, guildID string) (*Role, error) {
	var role Role
	err := p.db.First(&role, "server_role_id = ? and guild_id = ?", serverRoleID, guildID).Error
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		return nil, err
	}
	return &role, nil
}
