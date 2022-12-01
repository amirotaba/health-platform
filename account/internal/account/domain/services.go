package domain

import log "github.com/sirupsen/logrus"

type AccountServices struct {
	AccountUsecase     AccountUsecase
	ChannelUsecase     ChannelUsecase
	PermissionUsecase  PermissionUsecase
	RoleUsecase        RoleUsecase
	ServiceUsecase     ServiceUsecase
	AccountTypeUsecase AccountTypeUsecase
	ChannelRuleUsecase ChannelRuleUsecase
	AuthUsecase        AuthUsecase
	OtpUsecase         OtpUsecase
	TokenUsecase       TokenUsecase
	InventoryGrpc      InventoryGrpcPort
	RolePermission     RolePermissionUsecase
	PermissionServices PermissionServiceUsecase
	ChannelAccount     ChannelAccountUsecase
	Middleware         Middleware
	HttpPort           string
	GrpcPort           string
	Logger             *log.Logger
}
