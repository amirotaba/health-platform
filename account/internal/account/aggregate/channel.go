package aggregate

import (
	"context"
	"encoding/json"
	"git.paygear.ir/giftino/account/internal/account/entity"
	"git.paygear.ir/giftino/account/internal/account/event"
)

// Channel will save in database as event store
type Channel struct {
	BaseAggregate
	ChannelData entity.ChannelEntity
	RolesData   []entity.RoleEntity
	IsAdmin     bool
}

func (o *Channel) AggregateType() string {
	return ChannelAggregate
}

func (o *Channel) Restore(ctx context.Context, guid string, store event.Store) error {
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

func (o *Channel) ApplyEvents(ctx context.Context) error {
	for _, e := range o.Events {
		ev := e.(event.BaseEvent)
		switch ev.Type {
		case event.ChannelCreatedEventType:
			var data entity.CreateChannelRequest
			err := json.Unmarshal([]byte(ev.Data), &data)
			if err != nil {
				return err
			}
		case event.ChannelUpdatedEventType:
			var data entity.UpdateChannelRequest
			err := json.Unmarshal([]byte(ev.Data), &data)
			if err != nil {
				return err
			}

		}
	}

	o.Version = uint64(len(o.Events))
	return nil
}

func NewChannel(id string) Aggregate {
	agg := new(Channel)
	agg.UUID = id
	return agg
}
