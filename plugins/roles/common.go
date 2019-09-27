package roles

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

const (
	ServerRoleNotFound string = "roles.role.role-not-found-on-server"
)

func (p *Plugin) getServerRoleByNameOrID(input string, guildID string) (*discordgo.Role, error) {
	guild, err := p.state.Guild(guildID)
	if err != nil {
		return nil, err
	}

	for _, role := range guild.Roles {
		if input == role.Name || input == role.ID {
			return role, nil
		}
	}

	return nil, errors.New(ServerRoleNotFound)
}
