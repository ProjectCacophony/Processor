package autoroles

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

const (
	ServerRoleNotFound          string = "roles.role.role-not-found-on-server"
	MultipleServerRolesWithName string = "roles.role.multiple-server-roles-with-name"
)

func (p *Plugin) getServerRoleByNameOrID(input string, guildID string) (*discordgo.Role, error) {
	guild, err := p.state.Guild(guildID)
	if err != nil {
		return nil, err
	}

	var roles []*discordgo.Role
	for _, role := range guild.Roles {
		if input == role.Name || input == role.ID {
			roles = append(roles, role)
		}
	}

	if len(roles) == 0 {
		return nil, errors.New(ServerRoleNotFound)
	}

	if len(roles) > 1 {
		return nil, errors.New(MultipleServerRolesWithName)
	}

	return roles[0], nil
}
