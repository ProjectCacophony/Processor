package common

import (
	"strings"

	"gitlab.com/Cacophony/go-kit/interfaces"
	"gitlab.com/Cacophony/go-kit/permissions"
)

type ParamType int

const (
	Text ParamType = iota
	QuotedText
	Flag
	User
	Channel
	Link
	Attachment
	Duration
	DiscordInvite
)

type PluginHelp struct {
	Names               []string
	Description         string
	Commands            []Command
	Reactions           []Reaction
	Hide                bool
	PermissionsRequired Permissions
}

type Reaction struct {
	Name                string
	EmojiName           string
	Description         string
	PermissionsRequired Permissions
}

type Command struct {
	Name                string
	SkipRootCommand     bool
	SkipPrefix          bool
	Params              []CommandParam
	Description         string
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
		if permission == nil {
			permissionsText += "N/A, "
			continue
		}

		// skip Patron permissions, as we display it in a special way
		if permission == permissions.Patron {
			continue
		}

		permissionsText += permission.Name() + ", "
	}

	return strings.TrimSuffix(permissionsText, ", ")
}
