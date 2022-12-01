package command

import (
	"git.paygear.ir/giftino/account/internal/account/aggregate"
	"git.paygear.ir/giftino/account/internal/account/entity"
)

type CreateAccountCommand struct {
	BaseCommand
	Data   entity.CreateAccountRequest
	Result chan error
}

func NewCreateAccountCommand(id string, accountEntity entity.CreateAccountRequest) CreateAccountCommand {
	return CreateAccountCommand{
		BaseCommand: BaseCommand{
			Type:          AccountCreateType,
			AggregateID:   id,
			AggregateType: aggregate.AccountAggregate,
		},
		Data:   accountEntity,
		Result: make(chan error),
	}
}

// DeleteAccountCommand DeleteAccountBecauseUpdateQrFailedCommand
type DeleteAccountCommand struct {
	BaseCommand
	AccountId string
}

func NewDeleteAccountCommand(id string, accountUUID string) DeleteAccountCommand {
	return DeleteAccountCommand{
		BaseCommand: BaseCommand{
			Type:          AccountDeleteType,
			AggregateID:   id,
			AggregateType: aggregate.AccountAggregate,
		},
		AccountId: accountUUID,
	}
}

// UpdateAccountCommand UpdateAccountBecauseUpdateQrFailedCommand
type UpdateAccountCommand struct {
	BaseCommand
	Data entity.UpdateAccountRequest
}

func NewUpdateAccountCommand(id string, data entity.UpdateAccountRequest) UpdateAccountCommand {
	return UpdateAccountCommand{
		BaseCommand: BaseCommand{
			Type:          AccountUpdateType,
			AggregateID:   id,
			AggregateType: aggregate.AccountAggregate,
		},
		Data: data,
	}
}
