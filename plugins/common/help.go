package common

import (
	"strings"

	"gitlab.com/Cacophony/go-kit/interfaces"
)

type ParamType int

const (
	Text ParamType = iota
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

	ParamSets []ParamSet

	Hide        bool
	PatreonOnly bool

	PermissionsRequired Permissions
}

type ParamSet struct {
	PatreonOnly         bool
	Description         string
	Params              []PluginParam
	PermissionsRequired Permissions
}

type PluginParam struct {
	Name        string
	Example     string
	Type        ParamType
	Optional    bool
	NotVariable bool // indicates if the parameter is defined by the user
}

type Permissions []interfaces.Permission

func (p Permissions) String() (permissionsText string) {
	for _, permission := range p {
		permissionsText += permission.Name() + ", "
	}
	return strings.TrimSuffix(permissionsText, ", ")
}
