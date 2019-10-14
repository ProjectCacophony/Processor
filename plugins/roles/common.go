package roles

import (
	"errors"
	"time"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/events"
)

const (
	ServerRoleNotFound          string = "roles.role.role-not-found-on-server"
	MultipleServerRolesWithName string = "roles.role.multiple-server-roles-with-name"

	DELETE_DELAY time.Duration = 5
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

func (p *Plugin) deleteWithDelay(event *events.Event, messageID string) {
	defer recover()
	time.Sleep(DELETE_DELAY * time.Second)
	discord.Delete(event.Redis(), event.Discord(), event.ChannelID, messageID, false)
}
