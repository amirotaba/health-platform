package domain

import (
	"context"

	"git.paygear.ir/giftino/account/internal/account/entity"
)

type OtpUsecase interface {
	OtpVerify(ctx context.Context, in entity.OtpVerifyRequest) (entity.TokenEntity, error)
	NewOtp(ctx context.Context, in entity.OtpRequest) (string, error)
}

type OtpRepository interface {
	CreateOtp(ctx context.Context, in entity.OtpRequest) (int64, error)
	ReadOtp(ctx context.Context, request entity.CollectionRequest, response interface{}) error
}
