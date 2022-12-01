package aggregate

import "git.paygear.ir/giftino/account/internal/account/entity"

type PermissionReals struct {
	entity.PermissionEntity
	Services []entity.Service
}
