package aggregate

import (
	"context"
	"log"

	"git.paygear.ir/giftino/account/internal/account/entity"
	"git.paygear.ir/giftino/account/internal/account/event"
)

// ChannelRule will save in database as event store
type ChannelRule struct {
	BaseAggregate
	ChannelRuleData entity.ChannelRuleEntity
}

// CategoryChannelRule will save in database as event store
type CategoryChannelRule struct {
	entity.ChannelRuleCategory
	ChannelRules []entity.ChannelRuleEntity
}

func (o *ChannelRule) AggregateType() string {
	return ChannelRuleAggregate
}

func (o *ChannelRule) Restore(ctx context.Context, guid string, store event.Store) error {
	events, err := store.Load(guid)
	if err != nil {
		return err
	}

	o.Events = events
	err = o.ApplyEvents(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (o *ChannelRule) ApplyEvents(ctx context.Context) error {
	for _, e := range o.Events {
		ev := e.(event.BaseEvent)
		switch ev.Type {
		default:
			log.Println(ev)
		}
	}

	o.Version = uint64(len(o.Events))
	return nil
}

func NewChannelRule(id string) Aggregate {
	agg := new(ChannelRule)
	agg.UUID = id
	return agg
}
