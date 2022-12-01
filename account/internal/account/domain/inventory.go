package domain

import (
	"context"

	"git.paygear.ir/giftino/account/internal/account/entity"
)

type InventoryGrpcPort interface {
	TagExist(ctx context.Context, id int64) (entity.Tag, error)
}
