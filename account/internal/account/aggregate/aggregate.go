package aggregate

import (
	"context"

	"git.paygear.ir/giftino/account/internal/account/event"
)

const (
	AccountAggregate     = "account"
	ChannelAggregate     = "channel"
	ChannelRuleAggregate = "channel_rule"
)

type Aggregate interface {
	GetID() string
	IncrementVersion()
	GetVersion() uint64
	AggregateType() string
	Restore(ctx context.Context, guid string, store event.Store) error
	ApplyEvents(ctx context.Context) error
}

type BaseAggregate struct {
	ID      int64         `json:"id,omitempty"`
	UUID    string        `json:"uuid,omitempty"`
	Version uint64        `json:"version,omitempty"`
	Events  []event.Event `json:"events,omitempty"`
}

func (agg *BaseAggregate) GetID() string {
	return agg.UUID
}

func (agg *BaseAggregate) IncrementVersion() {
	agg.Version += 1
}

func (agg *BaseAggregate) GetVersion() uint64 {
	return agg.Version
}
