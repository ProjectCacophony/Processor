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
	err := p.db.
		Preload("Roles").
		First(&category, "name = ? and guild_id = ?", name, guildID).
		Error
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		return nil, err
	}
	return &category, nil
}

func (p *Plugin) getCategoryByChannel(channelID string) ([]*Category, error) {
	var categories []*Category
	err := p.db.
		Preload("Roles").
		Find(&categories, "channel_id = ?", channelID).
		Error
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		return nil, err
	}
	return categories, nil
}

func (p *Plugin) getAllRoles(guildID string) ([]*Role, error) {
	var roles []*Role
	err := p.db.
		Find(&roles, "guild_id = ?", guildID).
		Error
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		return nil, err
	}
	return roles, nil
}

func (p *Plugin) getUncategorizedRoles(guildID string) ([]*Role, error) {
	var roles []*Role
	err := p.db.
		Find(&roles, "guild_id = ? and (category_id is null or category_id = 0)", guildID).
		Error
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		return nil, err
	}
	return roles, nil
}

func (p *Plugin) getRoleByServerRoleID(serverRoleID string, guildID string) (*Role, error) {
	var role Role
	err := p.db.First(&role, "server_role_id = ? and guild_id = ?", serverRoleID, guildID).Error
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		return nil, err
	}
	return &role, nil
}
