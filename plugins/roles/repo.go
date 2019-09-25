package roles

import "strings"

func (p *Plugin) getAllCategories(guildId string) ([]*Category, error) {
	var categories []*Category
	err := p.db.Find(&categories, "guild_id = ?", guildId).Error
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		return nil, err
	}
	return categories, nil
}

func (p *Plugin) getCategoryByName(name string, guildId string) (*Category, error) {
	var category Category
	err := p.db.First(&category, "name = ? and guild_id = ?", name, guildId).Error
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		return nil, err
	}
	return &category, nil
}
