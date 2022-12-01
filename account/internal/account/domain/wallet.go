package domain

import (
	"context"

	pbWallet "git.paygear.ir/giftino/account/internal/account/proto/wallet"
)

type WalletGrpcPort interface {
	CreateWallet(ctx context.Context, uuid string, balance int64) (resp pbWallet.WalletReply, err error)
	GetWallet(ctx context.Context, request *pbWallet.GetWalletRequest) (resp pbWallet.WalletDetailsReply, err error)
}
