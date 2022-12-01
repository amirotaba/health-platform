package usecase

import (
	"context"
	"errors"
	"log"
	"time"

	"git.paygear.ir/giftino/account/internal/account/aggregate"
	"git.paygear.ir/giftino/account/internal/account/domain"
	"git.paygear.ir/giftino/account/internal/account/entity"
	"git.paygear.ir/giftino/account/internal/utils"
)

type usecase struct {
	accountRepo        domain.AccountRepository
	accountUsecase     domain.AccountUsecase
	roleRepo           domain.RoleRepository
	serviceRepo        domain.ServiceRepository
	permissionRepo     domain.PermissionRepository
	rolePermission     domain.RolePermissionUsecase
	permissionService  domain.PermissionServiceUsecase
	accountRoleUsecase domain.AccountRoleUsecase
	redis              domain.RedisPort
	jwt                domain.JWTService
}

func New(accountRepo domain.AccountRepository, redis domain.RedisPort, jwt domain.JWTService,
	rolePermission domain.RolePermissionUsecase, permissionService domain.PermissionServiceUsecase,
	serviceRepo domain.ServiceRepository, roleRepo domain.RoleRepository, accountUsecase domain.AccountUsecase, accountRoleUsecase domain.AccountRoleUsecase) domain.AuthUsecase {
	return &usecase{
		accountRepo:        accountRepo,
		serviceRepo:        serviceRepo,
		rolePermission:     rolePermission,
		permissionService:  permissionService,
		accountUsecase:     accountUsecase,
		accountRoleUsecase: accountRoleUsecase,
		roleRepo:           roleRepo,
		redis:              redis,
		jwt:                jwt,
	}
}

func (u usecase) SignUp(ctx context.Context, request entity.SingUpRequest) (string, error) {
	log.Println(request)
	phone, err := utils.CheckPhone(request.Phone)
	if err != nil {
		return "", err
	}

	getAccountReq := entity.NewCollectionRequest(domain.AccountTable).SetModel(entity.AccountModel{}).EQ("phone_number", phone)
	_, err = u.accountRepo.ReadOneAccount(ctx, getAccountReq)
	if err != nil {
		switch err.(type) {
		case utils.NotFoundError:
			_, err = u.accountUsecase.AddAccount(ctx, entity.CreateAccountRequest{
				PhoneNumber: phone,
				Password:    request.Password,
			})
			if err != nil {
				return "", err
			}

			otp := "1234"
			// todo: set token in redis
			return otp, nil
		default:
			return "", err
		}
	}

	// todo: set token in redis
	otp := "1234"
	return otp, nil
}

func (u usecase) SignIn(ctx context.Context, in entity.SingInRequest) (string, error) {
	phone, err := utils.CheckPhone(in.Phone)
	getAccountReq := entity.NewCollectionRequest(domain.AccountTable).SetModel(entity.AccountModel{}).EQ("phone_number", phone)
	//account, err := u.accountRepo.ReadAccountByOwnerPhoneNumber(ctx, phone)
	account, err := u.accountRepo.ReadOneAccount(ctx, getAccountReq)
	switch err.(type) {
	case utils.NotFoundError:
		_, err = u.accountUsecase.AddAccount(ctx, entity.CreateAccountRequest{PhoneNumber: phone, Password: in.Password})
		if err != nil {
			return "", err
		}
	}
	// todo error handling here
	if account.Password != in.Password {
		return "", errors.New("incorrect password account")
	}

	err = u.accountRepo.UpdateAccountLastLogin(ctx, phone, time.Now())
	if account.Password != in.Password {
		return "", errors.New("error when last login updated")
	}

	if !account.IsActive {
		return "", errors.New("user inactive")
	}
	//log.Println(account.ID, account.RoleId)
	return u.jwt.Generate(account.ID, nil, nil, domain.GiftinoSecretKey, time.Now().AddDate(0, 3, 0))
}

func (u usecase) Authentication(ctx context.Context, token, secretKey string) (resp aggregate.Account, err error) {
	if !u.jwt.Validate(token, secretKey) {
		err = errors.New("invalid token")
		return
	}

	claim, err := u.jwt.Extract(token, secretKey)
	if err != nil {
		return
	}

	getAccountRequest := entity.NewCollectionRequest(domain.AccountTable).SetModel(entity.AccountModel{}).EQ("id", claim.ID)
	accountEntity, err := u.accountRepo.ReadOneAccount(ctx, getAccountRequest)
	if err != nil {
		return
	}

	if !accountEntity.IsActive {
		return aggregate.Account{}, err
	}

	accountRolesRequest := entity.NewCollectionRequest(domain.AccountRoleTable).
		SetModel(entity.AccountRolesResponse{}).
		EQ("account_id", accountEntity.ID)
	accountRoles, _, err := u.accountRoleUsecase.FetchAccountRolesN(ctx, accountRolesRequest)
	if err != nil {
		return aggregate.Account{}, err
	}

	//log.Println("accountRoles :  ", accountRoles)
	//var accountRolesID []interface{}
	//for _, accountRole := range accountRoles {
	//	log.Println(accountRole.RoleID)
	//	accountRolesID = append(accountRolesID, accountRole.RoleID)
	//}

	var accountRolesID []interface{}
	if len(accountRoles) != 0 {
		log.Println("accountRoles :  ", accountRoles)
		for _, accountRole := range accountRoles {
			log.Println(accountRole.RoleID)
			accountRolesID = append(accountRolesID, accountRole.RoleID)
		}
		log.Println(accountRolesID)
	} else {
		accountRolesID = append(accountRolesID, accountEntity.RoleId)
	}

	//log.Println(accountRolesID)
	//if len(accountRoles) == 0 {
	//	accountRolesID = append(accountRolesID, accountEntity.RoleId)
	//}

	rolesRequest := entity.NewCollectionRequest(domain.RoleTable).SetModel(entity.RoleModel{})
	if len(accountRolesID) != 0 {
		//rolesRequest.Filters = append(rolesRequest.Filters, entity.Filter{
		//	Field: "id",
		//	Value: accountRolesID,
		//})
		log.Println("7777777777777777777777777777777777777777777777777777777777777777777")
		rolesRequest = rolesRequest.IN("id", accountRolesID...)
		roles, err := u.roleRepo.ReadManyRole(ctx, rolesRequest)
		if err != nil {
			return aggregate.Account{}, err
		}
		log.Println("*****************************", roles)
		resp.RolesData = roles
	}

	// todo : find best way
	var accountRolesNames []string
	for _, role := range resp.RolesData {
		accountRolesNames = append(accountRolesNames, role.Name)
	}

	if utils.StringInSlice(accountRolesNames, "administrator") {
		resp.IsAdmin = true
	}

	resp.AccountData = accountEntity
	return
}

func (u usecase) Authorization(ctx context.Context, roles []entity.RoleEntity, path, method string) (bool, error) {
	//permissions, err := u.roleRepo.ReadRolePermissions(ctx, roleID)
	//if err != nil {
	//	return false, err
	//}
	//
	//req := domain.NewCollectionRequest(domain.PermissionServiceTable, entity.ServiceResponse{})
	//u.permissionRepo.ReadPermissionServices(ctx, req)
	log.Println("-------------------------------------------------------------------")
	var req entity.FilterSearchRequest
	var rolesID []interface{}
	log.Println(roles)
	for _, role := range roles {
		rolesID = append(rolesID, role.ID)
	}

	req.Filters = append(req.Filters, entity.Filter{
		Field:  "role_id",
		Type:   entity.IN,
		Values: rolesID,
	})

	permissions, err := u.rolePermission.FetchRolePermissions(ctx, req)
	if err != nil {
		return false, err
	}

	for _, permission := range permissions {
		var req2 entity.FilterSearchRequest
		req2.Filters = append(req2.Filters, entity.Filter{
			Field: "permission_id",
			Type:  entity.EQ,
			Value: permission.PermissionID,
		})

		//req2.TableName = domain.PermissionServiceTable
		services, _, err := u.permissionService.FetchPermissionServices(ctx, req2)
		if err != nil {
			return false, err
		}

		for _, service := range services {
			var req3 entity.CollectionRequest
			req3.Filters = append(req3.Filters, entity.Filter{
				Field: "id",
				Type:  entity.EQ,
				Value: service.ServiceID,
			})

			req3.TableName = domain.ServicesTable
			var rms entity.ServiceResponse
			rms, err = u.serviceRepo.ReadOneService(ctx, req3)
			if err != nil {
				return false, err
			}

			//for _, manyService := range manyServices {
			if rms.Path == path && rms.Method == method && rms.IsActive {
				return true, nil
			}
			//}
		}
	}

	return false, nil
}
