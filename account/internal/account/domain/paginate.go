package domain

import "git.paygear.ir/giftino/account/internal/account/entity"

type PaginationRequest struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

type ServicePaginateResponseApp struct {
	Result     []entity.ServiceResponse `json:"result"`
	Pagination PaginateApp              `json:"Pagination"`
}

type AccessPaginateResponseApp struct {
	Result     []entity.RolePermissionsResponse `json:"result"`
	Pagination PaginateApp                      `json:"Pagination"`
}

type AccountRolePaginateResponseApp struct {
	Result     []entity.AccountRolesResponse `json:"result"`
	Pagination PaginateApp                   `json:"Pagination"`
}

type ChannelAccountPaginateResponseApp struct {
	Result     []entity.ChannelAccountsResponse `json:"result"`
	Pagination PaginateApp                      `json:"Pagination"`
}

type ChannelPaginateResponseApp struct {
	Result     []entity.ChannelResponse `json:"result"`
	Pagination PaginateApp              `json:"Pagination"`
}

type AccountPaginateResponseApp struct {
	Result     []entity.AccountResponse `json:"result"`
	Pagination PaginateApp              `json:"Pagination"`
}

type PermissionPaginateResponseApp struct {
	Result     []entity.PermissionResponse `json:"result"`
	Pagination PaginateApp                 `json:"Pagination"`
}

type RolePaginateResponseApp struct {
	Result     []entity.RoleResponse `json:"result"`
	Pagination PaginateApp           `json:"Pagination"`
}

type PermissionServicesPaginateResponseApp struct {
	Result     []entity.PermissionServicesResponse `json:"result"`
	Pagination PaginateApp                         `json:"Pagination"`
}

type PaginateApp struct {
	HasNext   bool `json:"has_next"`
	Page      int  `json:"page"`
	PerPage   int  `json:"per_page"`
	TotalPage int  `json:"total_page"`
	Total     int  `json:"total"`
}
