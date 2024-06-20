package utils

import (
	"eskimoe-server/config"
	"eskimoe-server/database"
)

func VerifyOwnerOrPermission(member database.Member, permission string) bool {
	authorized := false

	for _, role := range member.Roles {
		for _, permission := range role.Permissions {
			if permission == "manage_rooms" {
				authorized = true
				break
			}
		}
	}

	if !authorized {
		if config.Owner != member.UniqueID {
			authorized = false
		} else {
			authorized = true
		}
	}

	return authorized
}
