package domain

import (
	"context"
	"git.paygear.ir/giftino/account/internal/account/aggregate"
	"git.paygear.ir/giftino/account/internal/account/entity"
)

type AuthUsecase interface {
	SignUp(ctx context.Context, body entity.SingUpRequest) (string, error)
	SignIn(ctx context.Context, in entity.SingInRequest) (string, error)
	Authentication(ctx context.Context, token, secretKey string) (aggregate.Account, error)
	Authorization(ctx context.Context, rolesID []entity.RoleEntity, path, method string) (bool, error)
}
