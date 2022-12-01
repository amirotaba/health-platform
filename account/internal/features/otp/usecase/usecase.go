package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"git.paygear.ir/giftino/account/internal/account/domain"
	"git.paygear.ir/giftino/account/internal/account/entity"
	"git.paygear.ir/giftino/account/internal/utils"
)

type usecase struct {
	otpRepo            domain.OtpRepository
	tokenUsecase       domain.TokenUsecase
	accountRepo        domain.AccountRepository
	accountType        domain.AccountTypeRepository
	accountRoleUsecase domain.AccountRoleUsecase
	roleRepo           domain.RoleRepository
	rolePermRepo       domain.RolePermissionUsecase
	redis              domain.RedisPort
}

func New(otpRepo domain.OtpRepository, tokenUsecase domain.TokenUsecase, accountRepo domain.AccountRepository,
	roleRepo domain.RoleRepository, permRepo domain.RolePermissionUsecase,
	accountRoleUsecase domain.AccountRoleUsecase,
	accountTypeRepo domain.AccountTypeRepository, redis domain.RedisPort) domain.OtpUsecase {
	return &usecase{
		otpRepo:            otpRepo,
		tokenUsecase:       tokenUsecase,
		accountRepo:        accountRepo,
		accountRoleUsecase: accountRoleUsecase,
		redis:              redis,
		roleRepo:           roleRepo,
		rolePermRepo:       permRepo,
		accountType:        accountTypeRepo,
	}
}

func (u usecase) NewOtp(ctx context.Context, in entity.OtpRequest) (string, error) {
	phone, err := utils.CheckPhone(in.OwnerPhoneNumber)
	if err != nil {
		return "", err
	}

	var accountID int64
	getAccountReq := entity.NewCollectionRequest(domain.AccountTable).SetModel(entity.AccountModel{}).EQ("phone_number", phone)
	account, err := u.accountRepo.ReadOneAccount(ctx, getAccountReq)
	if err != nil {
		if err.Error() == "not found" {
			//todo : create new account with phone
			var req entity.CollectionRequest
			req.Filters = append(req.Filters, entity.Filter{
				Field: "name",
				Value: []interface{}{"guest"},
			})

			at, err := u.accountType.ReadOneAccountType(ctx, req)
			if err != nil {
				return "", err
			}

			var req2 entity.CollectionRequest
			req2.Filters = append(req2.Filters, entity.Filter{
				Field: "name",
				Value: []interface{}{"guest"},
			})

			role, err := u.roleRepo.ReadOneRole(ctx, req2)
			if err != nil {
				return "", err
			}

			a := entity.CreateAccountRequest{
				UUID:        uuid.New().String(),
				PhoneNumber: phone,
				TypeID:      at.ID,
				RoleID:      role.ID,
				IsActive:    true,
				ExpireAt:    time.Now().AddDate(0, 1, 0),
			}

			log.Println(a)
			createRequest := entity.NewInsertCollectionRequest(domain.AccountTable, a)
			accountID, err = u.accountRepo.CreateAccountNew(ctx, createRequest)
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	} else {
		accountID = account.ID
	}

	log.Println(accountID)
	otp := utils.Otp()
	_, err = u.otpRepo.CreateOtp(ctx, entity.OtpRequest{
		AccountID:        accountID,
		OwnerPhoneNumber: phone,
		OTP:              otp,
		ExpireAt:         time.Now().Add(time.Second * 120),
	})

	go utils.SMS(entity.SendSmsRequest{
		Mobile:  in.OwnerPhoneNumber,
		Message: otp,
	})

	go u.redis.Set(fmt.Sprintf("otp::%s", in.OwnerPhoneNumber), otp, time.Second*120)
	//go u.redis.Set(fmt.Sprintf("%s:%d", phone, accountID), otp, time.Second*120)
	if err != nil {
		return "", err
	}

	return otp, nil
}

func (u usecase) OtpVerify(ctx context.Context, in entity.OtpVerifyRequest) (token entity.TokenEntity, err error) {
	phone, err := utils.CheckPhone(in.PhoneNumber)
	if err != nil {
		return
	}

	// todo : get from redis is best practice
	var otp entity.OtpModel
	req := entity.NewCollectionRequest(domain.AccountOtpTable).
		SetModel(entity.OtpModel{}).
		EQ("phone_number", phone).
		EQ("otp", in.Otp)

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

	otp.OTP.String = u.redis.Get(fmt.Sprintf("otp::%s", in.PhoneNumber))
	if otp.OTP.String == "" {
		err = u.otpRepo.ReadOtp(ctx, req, &otp)
		if err != nil {
			return
		}

		if !otp.ExpireAt.Time.After(time.Now()) {
			err = utils.NewExpireError(otp.ExpireAt.Time, "otp")
			return
		}
	}

	if otp.OTP.String != in.Otp {
		err = utils.NewWrongOtpError(otp.OTP.String)
		return
	}

	//account, err := u.accountRepo.ReadAccountByOwnerPhoneNumber(ctx, phone)
	//if err != nil {
	//	return
	//}

	getAccountReq := entity.NewCollectionRequest(domain.AccountTable).
		SetModel(entity.AccountModel{}).
		EQ("phone_number", phone)
	account, err := u.accountRepo.ReadOneAccount(ctx, getAccountReq)
	if err != nil {
		return
	}
	// todo get all account role id

	if !account.IsActive {
		return entity.TokenEntity{}, errors.New("user inactive")
	}

	log.Println(account)
	log.Println(account.RoleId)
	accountRolesRequest := entity.NewCollectionRequest(domain.AccountRoleTable).SetModel(entity.AccountRolesResponse{}).EQ("account_id", account.ID)
	// optimize select query for get all account data like perm and role
	accountRoles, _, err := u.accountRoleUsecase.FetchAccountRolesN(ctx, accountRolesRequest)
	if err != nil {
		return entity.TokenEntity{}, err
	}

	var accountRolesID []int64
	if len(accountRoles) != 0 {
		log.Println("accountRoles :  ", accountRoles)
		for _, accountRole := range accountRoles {
			log.Println(accountRole.RoleID)
			accountRolesID = append(accountRolesID, accountRole.RoleID)
		}
		log.Println(accountRolesID)
	} else {
		accountRolesID = append(accountRolesID, account.RoleId)
	}

	var accountPermissionsID []int64
	for _, roleID := range accountRolesID {
		var rolePermsRequest entity.FilterSearchRequest
		rolePermsRequest.Filters = append(rolePermsRequest.Filters, entity.Filter{
			Field: "role_id",
			Type:  entity.EQ,
			Value: roleID,
		})

		permissions, err := u.rolePermRepo.FetchRolePermissions(ctx, rolePermsRequest)
		if err != nil {
			return entity.TokenEntity{}, err
		}

		log.Println(permissions)
		for _, permission := range permissions {
			accountPermissionsID = append(accountPermissionsID, permission.PermissionID)
		}
	}

	log.Println(accountPermissionsID)
	token, err = u.tokenUsecase.NewToken(ctx,
		entity.CreateTokenRequest{
			AccountID:     account.ID,
			RolesID:       accountRolesID,
			PermissionsID: accountPermissionsID,
			ExpireAt:      time.Now().AddDate(0, 3, 0)})
	if err != nil {
		return
	}

	return
}
