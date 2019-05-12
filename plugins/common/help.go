package common

import (
	"gitlab.com/Cacophony/go-kit/permissions"
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

	DiscordPermissionRequired *permissions.Discord
	BotPermissionRequired     *permissions.CacophonyBotPermission
}

type ParamSet struct {
	PatreonOnly               bool
	Description               string
	Params                    []PluginParams
	DiscordPermissionRequired *permissions.Discord
	BotPermissionRequired     *permissions.CacophonyBotPermission
}

type PluginParams struct {
	Name        string
	Example     string
	Type        ParamType
	Optional    bool
	NotVariable bool // indicates if the parameter is defined by the user
}
