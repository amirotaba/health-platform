package usecase

import (
	"context"
	"strconv"

	"git.paygear.ir/giftino/account/internal/account/domain"
)

type usecase struct {
	repo domain.AccountTypeRepository
}

func New(repo domain.AccountTypeRepository) domain.AccountTypeUsecase {
	return &usecase{repo: repo}
}

func (u usecase) AddAccountType(ctx context.Context, accountType *domain.AccountTypeEntity) error {
	return u.repo.CreateAccountType(ctx, accountType)
}

func (u usecase) FetchAccountType(ctx context.Context, ids []string) ([]domain.AccountTypeResponse, error) {
	var accountTypes []domain.AccountTypeResponse
	for _, id := range ids {
		accountTypeId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return nil, err
		}

		accountType, err := u.repo.ReadAccountType(ctx, accountTypeId)
		if err != nil {
			return nil, err
		}

		accountTypes = append(accountTypes, newResponse(accountType))
	}

	return accountTypes, nil
}

func (u usecase) FetchAccountTypes(ctx context.Context) ([]domain.AccountTypeResponse, error) {
	types, err := u.repo.ReadAccountTypes(ctx)
	var resp []domain.AccountTypeResponse
	if err != nil {
		return resp, err
	}

	for _, t := range types {
		resp = append(resp, newResponse(t))
	}

	return resp, nil
}

func (u usecase) PatchAccountType(ctx context.Context, accountType *domain.AccountTypeUpdateRequest) error {
	_, err := u.repo.ReadAccountType(ctx, accountType.ID)
	if err != nil {
		return err
	}

	return u.repo.UpdateAccountType(ctx, accountType)
}

func newResponse(entity domain.AccountTypeEntity) domain.AccountTypeResponse {
	return domain.AccountTypeResponse{
		ID:          entity.ID,
		Name:        entity.Name,
		IsActive:    entity.IsActive,
		Description: entity.Description,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}
