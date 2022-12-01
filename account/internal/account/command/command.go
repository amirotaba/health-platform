package command

import (
	"git.paygear.ir/giftino/account/internal/account/entity"
	"git.paygear.ir/giftino/account/internal/account/event"
)

const (
	AccountCreateType = "create:account"
	AccountDeleteType = "delete:account"
	ChannelCreateType = "create:channel"
	ChannelChargeType = "charge:channel"
	AccountUpdateType = "update:account"
)

type Command interface {
	GetAggregateID() string
	GetAggregateType() string
	GetCommandType() string
}

//Message is a message command implementation interface
type Message interface {
	Command
	Payload() interface{}
	CommandType() string
}

type Handler interface {
	ApplyChange(events []event.Event) (entity.AccountEntity, error)
	HandleCommand()
	Restore(uuid string) (entity.AccountEntity, error)
}

type BaseCommand struct {
	Type          string
	AggregateID   string
	AggregateType string
}

func (c BaseCommand) GetAggregateID() string {
	return c.AggregateID
}

func (c BaseCommand) GetAggregateType() string {
	return c.AggregateType
}

func (c BaseCommand) GetCommandType() string {
	return c.Type
}
