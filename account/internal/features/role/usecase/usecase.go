package usecase

import (
	"context"
	"strconv"

	log "github.com/sirupsen/logrus"

	"git.paygear.ir/giftino/account/internal/account/domain"
	"git.paygear.ir/giftino/account/internal/account/entity"
)

type usecase struct {
	roleRepo       domain.RoleRepository
	permissionRepo domain.PermissionRepository
	rolePerRepo    domain.RolePermissionRepository
	logger         *log.Logger
}

func New(
	roleRepo domain.RoleRepository,
	permissionRepo domain.PermissionRepository,
	rolePerRepo domain.RolePermissionRepository, logger *log.Logger) domain.RoleUsecase {
	return &usecase{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		rolePerRepo:    rolePerRepo,
		logger:         logger,
	}
}

func (u usecase) AddRole(ctx context.Context, role entity.CreateRoleRequest) (int64, error) {
	req := entity.NewInsertCollectionRequest(domain.RoleTable, role)
	return u.roleRepo.CreateRole(ctx, req)
	//id, err := u.roleRepo.CreateRole(ctx, req)
	//if err != nil {
	//	return 0, err
	//}

	//for _, pole := range role.Roles {
	//	rp := entity.CreateRolesToRoleRequest{
	//		RoleId:       id,
	//		RoleId: int64(pole),
	//	}
	//
	//	createRequest := entity.NewInsertCollectionRequest(domain.RolesRolesTable, rp)
	//	var rpID int64
	//	rpID, err = u.rolePerRepo.CreateRoleRoles(ctx, createRequest)
	//	if err != nil {
	//		return 0, err
	//	}
	//
	//	log.Printf("premission with id %d assigned to role with id %d \n", rpID, id)
	//}

	//return id, nil
}

func (u usecase) FetchRole(ctx context.Context, ids []string) ([]entity.RoleResponse, error) {
	var roles []entity.RoleResponse
	for _, id := range ids {
		roleId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return nil, err
		}

		role, err := u.roleRepo.ReadRole(ctx, roleId)
		if err != nil {
			return nil, err
		}

		//role.Roles, err = u.roleRepo.ReadRoleRoles(ctx, role.ID)
		//if err != nil {
		//	role.Description += "\n" + err.Error()
		//}

		roles = append(roles, role.ToResponse())
	}

	return roles, nil
}

func (u usecase) FetchRoleByID(ctx context.Context, id int64) (entity.RoleResponse, error) {
	role, err := u.roleRepo.ReadRole(ctx, id)
	if err != nil {
		return entity.RoleResponse{}, err
	}

	//role.Roles, err = u.roleRepo.ReadRoleRoles(ctx, role.ID)
	//if err != nil {
	//	role.Description += "\n" + err.Error()
	//}

	return role.ToResponse(), nil
}

func (u usecase) FetchRoles(ctx context.Context) ([]entity.RoleResponse, error) {
	roles, err := u.roleRepo.ReadRoles(ctx)
	var resp []entity.RoleResponse
	if err != nil {
		return resp, err
	}

	for _, role := range roles {
		//role.Roles, err = u.roleRepo.ReadRoleRoles(ctx, role.ID)
		//if err != nil {
		//	role.Description += "\n" + err.Error()
		//}
		//
		resp = append(resp, role.ToResponse())
	}

	return resp, nil
}

func (u usecase) FetchRoleWithPaginate(ctx context.Context, request entity.FilterSearchRequest) ([]entity.RoleResponse, domain.PaginateApp, error) {
	req := entity.CollectionRequest{Filters: request.Filters, Searches: request.Searches}
	req.TableName = domain.RoleTable
	req.Model = entity.RoleModel{}
	roles, err := u.roleRepo.ReadManyRole(ctx, req)
	var resp []entity.RoleResponse
	var paginate domain.PaginateApp
	if err != nil {
		return resp, paginate, err
	}

	paginate.Total = len(roles)
	paginate.Page = request.Page
	paginate.PerPage = request.PerPage
	paginate.TotalPage = paginate.Total / paginate.PerPage
	u.logger.Println(paginate.Total)
	u.logger.Println(paginate.TotalPage)
	if paginate.Total%request.PerPage > 0 {
		paginate.TotalPage += 1
	}

	if paginate.TotalPage >= paginate.Page {
		paginate.HasNext = true
	} else {
		paginate.HasNext = false
	}

	for _, role := range roles {
		resp = append(resp, role.ToResponse())
	}
	return resp, paginate, nil
}

func (u usecase) PatchRole(ctx context.Context, role entity.RoleUpdateRequest) error {
	_, err := u.roleRepo.ReadRole(ctx, role.ID)
	if err != nil {
		return err
	}

	return u.roleRepo.UpdateRole(ctx, role)
}
