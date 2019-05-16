package help

import (
	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/permissions"
)

func containsBotAdminPermission(permissionsInput common.Permissions) bool {
	for _, permission := range permissionsInput {
		if permission == permissions.BotAdmin {
			return true
		}
	}

	return false
}
