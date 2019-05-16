package common

import (
	"strings"

	"gitlab.com/Cacophony/go-kit/interfaces"
)

type ParamType int

const (
	Text ParamType = iota
	QuotedText
	Hardcoded
	User
	Channel
	Link
	Attachment
	Duration
	DiscordInvite
)

type PluginHelp struct {
	Name        string
	Description string

	Commands []Command

	Hide        bool
	PatreonOnly bool

	PermissionsRequired Permissions
}

type Command struct {
	Name                string
	PatreonOnly         bool
	Description         string
	Params              []CommandParam
	PermissionsRequired Permissions
}

type CommandParam struct {
	Name     string
	Example  string
	Type     ParamType
	Optional bool
}

type Permissions []interfaces.Permission

func (p Permissions) String() (permissionsText string) {
	for _, permission := range p {
		permissionsText += permission.Name() + ", "
	}
	return strings.TrimSuffix(permissionsText, ", ")
}
