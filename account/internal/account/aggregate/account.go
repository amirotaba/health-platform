package aggregate

import (
	"context"
	"encoding/json"
	"errors"

	"git.paygear.ir/giftino/account/internal/account/entity"
	"git.paygear.ir/giftino/account/internal/account/event"
)

// Account will save in database as event store
type Account struct {
	BaseAggregate
	AccountData entity.AccountEntity
	RolesData   []entity.RoleEntity
	IsAdmin     bool
}

func (o *Account) AggregateType() string {
	return AccountAggregate
}

func (o *Account) Restore(ctx context.Context, guid string, store event.Store) error {
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

func (o *Account) ApplyEvents(ctx context.Context) error {
	for _, e := range o.Events {
		ev := e.(event.BaseEvent)
		switch ev.Type {
		case event.AccountCreatedEventType:
			var data entity.CreateAccountRequest
			err := json.Unmarshal([]byte(ev.Data), &data)
			if err != nil {
				return err
			}

			o.AccountData = data.ToEntity()
		case event.AccountRoleUpdatedEventType:
			var data entity.UpdateAccountRole
			err := json.Unmarshal([]byte(ev.Data), &data)
			if err != nil {
				return err
			}

			//if o.Data.RoleId != data.OldRoleID {
			//	return errors.New("in load aggregate")
			//}

			//o.Data.RoleId = data.NewRoleID
		case event.AccountTypeUpdatedEventType:
			var data entity.UpdateAccountType
			err := json.Unmarshal([]byte(ev.Data), &data)
			if err != nil {
				return err
			}

			if o.AccountData.TypeId != data.OldTypeID {
				return errors.New("in load aggregate")
			}

			//o.Data.RoleId = data.NewTypeID
		case event.AccountUpdatedEventType:
			var data entity.UpdateAccountRequest
			err := json.Unmarshal([]byte(ev.Data), &data)
			if err != nil {
				return err
			}

			o.AccountData = data.ToEntity(o.AccountData)
		}
	}

	o.Version = uint64(len(o.Events))
	return nil
}

func NewAccount(id string) Aggregate {
	agg := new(Account)
	agg.UUID = id
	return agg
}
