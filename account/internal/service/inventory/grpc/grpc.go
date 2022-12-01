package grpc_proto

import (
	"context"
	"git.paygear.ir/giftino/account/internal/account/entity"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"git.paygear.ir/giftino/account/internal/account/domain"
)

type grpcPort struct {
	address string
	conn    *grpc.ClientConn
}

func New(address string) domain.InventoryGrpcPort {
	return grpcPort{address: address}
}

func (g grpcPort) TagExist(ctx context.Context, id int64) (entity.Tag, error) {
	conn, err := grpc.Dial(g.address, grpc.WithInsecure())
	if err != nil {
		log.Printf("did not connect: %v \n", err)
		return entity.Tag{}, err
	}

	defer conn.Close()
	c := NewTagServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Exist(ctx, &TagRequest{
		ID: id,
	})

	if err != nil {
		log.Printf("could not check: %v \n", err)
		return entity.Tag{}, err
	}

	return entity.Tag{
		Exist:        r.GetExist(),
		Status:       r.GetStatus(),
		Name:         r.GetName(),
		ID:           r.GetID(),
		CategoryName: r.GetCategoryName(),
		CategoryID:   r.GetCategoryID(),
	}, nil
}
