package grpc

import (
	"context"

	"git.paygear.ir/giftino/account/internal/account/domain"
	"git.paygear.ir/giftino/account/internal/utils"
)

type server struct {
	usecase domain.AuthUsecase
}

func New(usecase domain.AuthUsecase) AuthenticationServiceServer {
	return server{usecase: usecase}
}

func (s server) Authentication(ctx context.Context, request *AuthenticationRequest) (*AuthenticationReply, error) {
	entity, err := s.usecase.Authentication(ctx, request.Token, request.Secret)
	if err != nil {
		return nil, err
	}

	ok := false
	if utils.StringInSlice(utils.RoleToString(entity.RolesData), "administrator") {
		ok = true
	} else {
		ok, err = s.usecase.Authorization(ctx, entity.RolesData, request.Path, request.Token)
		if err != nil {
			return nil, err
		}
	}

	return &AuthenticationReply{
		ID:        entity.ID,
		UUID:      entity.UUID,
		FirstName: entity.AccountData.FirstName,
		LastName:  entity.AccountData.LastName,
		Email:     entity.AccountData.Email,
		RoleID:    0, //entity.AccountData.RoleId,
		IsActive:  entity.AccountData.IsActive,
		Access:    ok,
	}, nil
}

func (s server) mustEmbedUnimplementedAuthenticationServiceServer() {}
