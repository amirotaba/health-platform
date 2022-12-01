package grpc

import (
	"context"
	"fmt"
	"io"

	"git.paygear.ir/giftino/account/internal/account/domain"
	"git.paygear.ir/giftino/account/internal/account/entity"
)

type server struct {
	service domain.ServiceUsecase
}

func New(because domain.ServiceUsecase) ServiceServiceServer {
	return server{service: because}
}

func (s server) AddOneService(ctx context.Context, request *ServiceRequest) (*ServiceReply, error) {
	insertedID, err := s.service.NewService(ctx, entity.CreateServiceRequest{
		Name:        request.GetName(),
		Code:        request.GetCode(),
		Path:        request.GetPath(),
		Function:    request.GetFunction(),
		Method:      request.GetMethod(),
		IsActive:    request.GetActive(),
		Description: "",
	})

	if err != nil {
		return &ServiceReply{
			Inserted: false,
			Message:  err.Error(),
		}, err
	}

	return &ServiceReply{
		Inserted: true,
		Message:  fmt.Sprintf("service with id %d iserted", insertedID),
	}, nil
}

func (s server) AddManyService(serviceServer ServiceService_AddManyServiceServer) error {
	var reply *StreamServiceReply
	for {
		service, err := serviceServer.Recv()
		if err == io.EOF {
			return serviceServer.SendAndClose(reply)
		}

		if err != nil {
			return err
		}

		ctx := context.TODO()
		id, err := s.service.NewService(ctx, entity.CreateServiceRequest{
			Name:        service.GetName(),
			Code:        service.GetCode(),
			Path:        service.GetPath(),
			Function:    service.GetFunction(),
			Method:      service.GetMethod(),
			IsActive:    service.GetActive(),
			Description: service.GetDescription(),
		})

		if err != nil {
			reply.Resp = append(reply.Resp, &ServiceReply{
				Inserted: false,
				Message:  err.Error(),
			})
		}

		reply.Resp = append(reply.Resp, &ServiceReply{
			Inserted: true,
			Message:  fmt.Sprintf("service with id %d iserted", id),
		})

	}
}

func (s server) mustEmbedUnimplementedServiceServiceServer() {}
