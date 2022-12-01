package event

import "git.paygear.ir/giftino/account/internal/account/entity"

const (
	AccountCreatedEventType     = "create:account"
	AccountTypeUpdatedEventType = "updateType:account"
	AccountRoleUpdatedEventType = "updateRole:account"
	AccountUpdatedEventType     = "update:account"
	AccountDeletedEventType     = "delete:account"
)

type AccountCreatedEvent struct {
	BaseEvent
	Data entity.CreateAccountRequest
}

type AccountDepositedEvent struct {
	BaseEvent
}

type AccountDeletedEvent struct {
	BaseEvent
	AccountId string
}

type AccountUpdatedEvent struct {
	BaseEvent
	Data entity.UpdateAccountRequest
}

type AccountRoleUpdatedEvent struct {
	BaseEvent
	Data entity.UpdateAccountRole
}

type AccountTypeUpdatedEvent struct {
	BaseEvent
	Data entity.UpdateAccountType
}
