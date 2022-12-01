package grpc

import (
	"context"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"git.paygear.ir/giftino/account/internal/account/domain"
	pbWallet "git.paygear.ir/giftino/account/internal/account/proto/wallet"
)

type grpcPort struct {
	address string
}

func New(address string) domain.WalletGrpcPort {
	return grpcPort{address: address}
}

func (g grpcPort) CreateWallet(ctx context.Context, uuid string, balance int64) (resp pbWallet.WalletReply, err error) {
	conn, err := grpc.Dial(g.address, grpc.WithInsecure())
	if err != nil {
		log.Printf("did not connect: %v", err)
		return //false, err
	}

	defer conn.Close()
	c := pbWallet.NewWalletServiceClient(conn)

	// Contact the server and print out its response.
	r, err := c.CreateWallet(ctx, &pbWallet.WalletRequest{
		UserID:  uuid,
		Balance: balance,
	})
	if err != nil {
		log.Printf("could not check: %v", err)
		return //false, err
	}

	resp = *r
	return
}

func (g grpcPort) GetWallet(ctx context.Context, request *pbWallet.GetWalletRequest) (resp pbWallet.WalletDetailsReply, err error) {
	conn, err := grpc.Dial(g.address, grpc.WithInsecure())
	if err != nil {
		log.Printf("did not connect: %v", err)
		return //false, err
	}

	defer conn.Close()
	c := pbWallet.NewWalletServiceClient(conn)

	// Contact the server and print out its response.
	r, err := c.GetWallet(ctx, request)
	if err != nil {
		log.Printf("could not check: %v", err)
		return //false, err
	}

	resp = *r
	return
}
