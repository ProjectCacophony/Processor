package autoroles

import "strings"

func (p *Plugin) getAutoRoles(guildID string) ([]*AutoRole, error) {
	var autoRoles []*AutoRole
	err := p.db.
		Preload("Rule").
		Preload("Rule.Actions").
		Find(&autoRoles, "guild_id = ?", guildID).Error
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		return nil, err
	}
	return autoRoles, nil
}

func (p *Plugin) getAutoRolesByServerRoleID(serverRoleID string, guildID string) (*AutoRole, error) {
	var aRole AutoRole
	err := p.db.
		Preload("Rule").
		Preload("Rule.Actions").
		First(&aRole, "guild_id = ? and server_role_id = ?", guildID, serverRoleID).Error
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		return nil, err
	}
	return &aRole, nil
}
