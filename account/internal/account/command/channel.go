package command

import (
	"git.paygear.ir/giftino/account/internal/account/aggregate"
	"git.paygear.ir/giftino/account/internal/account/entity"
)

type CreateChannelCommand struct {
	BaseCommand
	Data   entity.CreateChannelRequest
	Result chan error
}

func NewCreateChannelCommand(id string, accountEntity entity.CreateChannelRequest) CreateChannelCommand {
	return CreateChannelCommand{
		BaseCommand: BaseCommand{
			Type:          ChannelCreateType,
			AggregateID:   id,
			AggregateType: aggregate.ChannelAggregate,
		},
		Data:   accountEntity,
		Result: make(chan error),
	}
}

// DeleteChannelCommand DeleteChannelBecauseUpdateQrFailedCommand
type DeleteChannelCommand struct {
	BaseCommand
	ChannelId string
}

// ChargeChannelCommand UpdateChannelBecauseUpdateQrFailedCommand
type ChargeChannelCommand struct {
	BaseCommand
	Data entity.ChargeChannelRequest
}

func NewChargeChannelCommand(id string, data entity.ChargeChannelRequest) ChargeChannelCommand {
	return ChargeChannelCommand{
		BaseCommand: BaseCommand{
			Type:          ChannelChargeType,
			AggregateID:   id,
			AggregateType: aggregate.ChannelAggregate,
		},
		Data: data,
	}
}
