package help

import (
	"gitlab.com/Cacophony/go-kit/permissions"
)

type PluginHelp struct {
	PluginName  string
	Description string

	ParamSets []ParamSet

	Hide        bool
	PatreonOnly bool

	DiscordPermissionRequired permissions.Discord
	BotPermissionRequired     permissions.CacophonyBotOwner
}

type ParamSet struct {
	Params []PluginParams

	DiscordPermissionRequired permissions.Discord
	BotPermissionRequired     permissions.CacophonyBotOwner
}

type PluginParams struct {
	ParamName    string
	ParamType    string
	ParamExample string
}
