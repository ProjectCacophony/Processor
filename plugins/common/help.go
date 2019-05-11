package common

import (
	"gitlab.com/Cacophony/go-kit/permissions"
)

type PluginHelp struct {
	Name        string
	Description string

	ParamSets []ParamSet

	Hide        bool
	PatreonOnly bool

	DiscordPermissionRequired permissions.Discord
	BotPermissionRequired     permissions.CacophonyBotPermission
}

type ParamSet struct {
	Params []PluginParams

	DiscordPermissionRequired permissions.Discord
	BotPermissionRequired     permissions.CacophonyBotPermission
}

type PluginParams struct {
	Name    string
	Type    string
	Example string
}
