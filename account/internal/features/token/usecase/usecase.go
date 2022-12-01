package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/ksuid"

	"git.paygear.ir/giftino/account/internal/account/domain"
	"git.paygear.ir/giftino/account/internal/account/entity"
	"git.paygear.ir/giftino/account/internal/utils"
)

type usecase struct {
	tokenRepo domain.TokenRepository
	jwt       domain.JWTService
	redis     domain.RedisPort
}

func New(tokenRepo domain.TokenRepository, jwtService domain.JWTService, redis domain.RedisPort) domain.TokenUsecase {
	return &usecase{
		tokenRepo: tokenRepo,
		jwt:       jwtService,
		redis:     redis,
	}
}

func (u usecase) NewToken(ctx context.Context, in entity.CreateTokenRequest) (entity.TokenEntity, error) {
	var token entity.TokenEntity
	authToken, err := u.jwt.Generate(in.AccountID, in.RolesID, in.PermissionsID, domain.GiftinoSecretKey, in.ExpireAt)
	if err != nil {
		return entity.TokenEntity{}, err
	}

	token.UUID = uuid.New().String()
	token.AuthToken = authToken
	token.AccountID = in.AccountID
	token.Type = domain.JWTToken
	token.RefreshToken = ksuid.New().String()
	token.ExpireAt = in.ExpireAt
	token.ID, err = u.tokenRepo.CreateToken(ctx, token)
	token.CreatedAt = time.Now()
	if err != nil {
		return entity.TokenEntity{}, err
	}

	return token, nil
}

func (u usecase) RefreshToken(ctx context.Context, in entity.RefreshTokenRequest) (token entity.TokenEntity, err error) {
	req := entity.NewCollectionRequest(domain.TokenTable).
		SetModel(entity.TokenModel{}).
		EQ("auth_token", in.OldAuthToken).
		EQ("refresh_token", in.RefreshToken)

	req.OrderBy = entity.OrderBy{
		Has:       true,
		Field:     domain.KeyCreatedAt,
		Ascending: false,
	}

	req.PaginateData = entity.Paginate{
		Has:    true,
		Limit:  1,
		Offset: 0,
	}

	//var token entity.TokenModel
	err = u.tokenRepo.ReadToken(ctx, req, &token)
	if err != nil {
		return token, err
	}

	claim, err := u.jwt.Extract(in.OldAuthToken, domain.GiftinoSecretKey)
	if err != nil {
		return token, err
	}

	if !claim.ExpiresAt.After(time.Now()) {
		err = utils.NewExpireError(claim.ExpiresAt.Time, "token")
	}

	var newToken entity.UpdateTokenRequest
	newToken.RefreshToken = in.RefreshToken
	newToken.AuthToken, err = u.jwt.Generate(claim.ID, claim.RolesID, claim.PermissionsID, domain.GiftinoSecretKey, time.Now().AddDate(0, 3, 0))
	if err != nil {
		return entity.TokenEntity{}, err
	}

	updateRequest := entity.NewCollectionRequest(domain.TokenTable).SetModel(newToken).EQ("refresh_token", in.RefreshToken)

	err = u.tokenRepo.UpdateToken(ctx, updateRequest)
	if err != nil {
		return entity.TokenEntity{}, err
	}

	return token, nil
}
