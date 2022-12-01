package domain

import (
	"context"

	"git.paygear.ir/giftino/account/internal/account/entity"
)

type TokenUsecase interface {
	NewToken(ctx context.Context, in entity.CreateTokenRequest) (entity.TokenEntity, error)
	RefreshToken(ctx context.Context, in entity.RefreshTokenRequest) (entity.TokenEntity, error)
}

type TokenRepository interface {
	CreateToken(ctx context.Context, in entity.TokenEntity) (int64, error)
	ReadToken(ctx context.Context, request entity.CollectionRequest, response interface{}) error
	UpdateToken(ctx context.Context, request entity.CollectionRequest) error
}
