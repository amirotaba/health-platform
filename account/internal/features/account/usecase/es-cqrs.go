package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"git.paygear.ir/giftino/account/internal/account/aggregate"
	"git.paygear.ir/giftino/account/internal/account/command"
	"git.paygear.ir/giftino/account/internal/account/domain"
	"git.paygear.ir/giftino/account/internal/account/entity"
	"git.paygear.ir/giftino/account/internal/account/event"
	"log"
)

func (u *usecase) ProcessCommand(ctx context.Context, cmd command.Command) ([]event.Event, error) {
	var err error
	accountAggregate := aggregate.NewAccount(cmd.GetAggregateID())
	err = accountAggregate.Restore(ctx, cmd.GetAggregateID(), u.eventStore)
	if err != nil {
		return nil, err
	}

	account := accountAggregate.(*aggregate.Account)
	if err != nil {
		return nil, err
	}

	switch c := cmd.(type) {
	case command.CreateAccountCommand:
		fmt.Println("------------------------------CreateAccountCommand------------")
		defer fmt.Println("------------------------------CreateAccountCommand------------")
		account.AccountData = c.Data.ToEntity()
		AccountEntity := mapAggregateToDTO(*account)
		createAccountRequest := entity.NewInsertCollectionRequest(domain.AccountTable, c.Data)
		AccountEntity.ID, err = u.accountRepo.CreateAccountNew(ctx, createAccountRequest)
		if err != nil {
			return nil, err
		}

		account.ID = AccountEntity.ID
		// todo define a projection that account upsert into it
		data, err := json.Marshal(c.Data)
		if err != nil {
			data = []byte(err.Error())
		}

		evt := event.AccountCreatedEvent{
			BaseEvent: event.BaseEvent{
				AggregateID:   c.AggregateID,
				AggregateType: c.GetAggregateType(),
				Type:          c.Type,
				Version:       int64(accountAggregate.GetVersion() + 1),
				Tag:           event.Tag,
				Data:          string(data),
			},
			Data: c.Data,
		}

		return []event.Event{evt}, err
	case command.UpdateAccountCommand:
		fmt.Println("------------------------------CreateAccountCommand------------")
		defer fmt.Println("------------------------------CreateAccountCommand------------")
		//account.Data = c.Data.ToEntity()
		AccountEntity := mapAggregateToDTO(*account)
		err = u.accountRepo.UpdateAccount(ctx, c.Data)
		if err != nil {
			return nil, err
		}

		account.ID = AccountEntity.ID
		// todo define a projection that account upsert into it
		data, err := json.Marshal(c.Data)
		if err != nil {
			data = []byte(err.Error())
		}

		evt := event.AccountUpdatedEvent{
			BaseEvent: event.BaseEvent{
				AggregateID:   c.AggregateID,
				AggregateType: c.GetAggregateType(),
				Type:          c.Type,
				Version:       int64(accountAggregate.GetVersion() + 1),
				Tag:           event.Tag,
				Data:          string(data),
			},
			Data: c.Data,
		}

		return []event.Event{evt}, err
	}

	return nil, errors.New("this command not handled")
}

func (u *usecase) HandleCommand() {
	// todo define an log store for save error in go routine
	for {
		cmd := <-u.CommandChan
		ctx := context.Background()
		err := u.handleCommand(ctx, cmd)
		if err != nil {
			log.Println(err)
		}
	}
}

func (u *usecase) handleCommand(ctx context.Context, cmd command.Command) error {
	events, err := u.ProcessCommand(ctx, cmd)
	if err != nil {
		return err
	}

	err = u.eventStore.Store(events)
	if err != nil {
		log.Println(err)
		return err
	}

	err = u.eventBus.Publish(events)
	if err != nil {
		return err
	}
	return nil
}

func (u usecase) GetCommandChan() chan command.Command {
	return u.CommandChan
}

func mapAggregateToDTO(accountAggregate aggregate.Account) entity.AccountEntity {
	//accountAggregate := agg.(*aggregate.AccountModel)
	return entity.AccountEntity{
		UUID:        accountAggregate.UUID,
		FirstName:   accountAggregate.AccountData.FirstName,
		LastName:    accountAggregate.AccountData.LastName,
		DisplayName: accountAggregate.AccountData.DisplayName,
		Password:    accountAggregate.AccountData.Password,
		Email:       accountAggregate.AccountData.Email,
		PhoneNumber: accountAggregate.AccountData.PhoneNumber,
		Address:     accountAggregate.AccountData.Address,
		TypeId:      accountAggregate.AccountData.TypeId,
		IsActive:    accountAggregate.AccountData.IsActive,
		ExpireAt:    accountAggregate.AccountData.ExpireAt,
		LastLogin:   accountAggregate.AccountData.LastLogin,
		Description: accountAggregate.AccountData.Description,
		CreatedAt:   accountAggregate.AccountData.CreatedAt,
		UpdatedAt:   accountAggregate.AccountData.UpdatedAt,
		DeletedAt:   accountAggregate.AccountData.DeletedAt,
	}
}
