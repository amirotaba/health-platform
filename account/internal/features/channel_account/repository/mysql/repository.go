package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"

	"git.paygear.ir/giftino/account/internal/account/domain"
	"git.paygear.ir/giftino/account/internal/account/entity"
)

type repository struct {
	db     *sql.DB
	dbx    *sqlx.DB
	logger *log.Logger
}

func New(db *sql.DB, logger *log.Logger) domain.ChannelAccountRepository {
	return &repository{
		db:     db,
		dbx:    sqlx.NewDb(db, "mysql"),
		logger: logger,
	}
}

func (r repository) CreateChannelAccounts(ctx context.Context, req entity.InsertCollectionRequest) (id int64, err error) {
	create := sqlbuilder.NewStruct(req.Model)
	sb := create.InsertInto(req.TableName, &req.Model)
	query, args := sb.Build()
	r.logger.Info(query, args)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	res, err := r.dbx.ExecContext(ctx, query, args...)
	if err != nil {
		r.logger.Printf("Error %s when inserting row into permission table", err)
		return
	}

	rows, err := res.RowsAffected()
	if err != nil {
		r.logger.Printf("Error %s when finding rows affected", err)
		return
	}

	r.logger.Printf("%d assigned channel to account successfuly ", rows)
	id, _ = res.LastInsertId()
	return
}

func (r repository) ReadChannelAccounts(ctx context.Context, request entity.CollectionRequest) ([]entity.ChannelAccountEntity, error) {
	var resp []entity.ChannelAccountEntity
	selectBuilder := sqlbuilder.NewStruct(entity.ChannelAccountsResponse{})
	sb := selectBuilder.SelectFrom(domain.ChannelAccountTable)
	for _, filter := range request.Filters {
		switch filter.Type {
		//case domain.KeyFromDate:
		//	sb.Where(sb.GE(domain.KeyFromDate, filter.Value[0]))
		//case domain.KeyToDate:
		//	sb.Where(sb.LE(domain.KeyToDate, filter.Value[0]))
		case entity.EQ:
			sb.Where(sb.E(filter.Field, filter.Value))
		case entity.GE:
			sb.Where(sb.GE(filter.Field, filter.Value))
		case entity.LE:
			sb.Where(sb.LE(filter.Field, filter.Value))
		case entity.GT:
			sb.Where(sb.G(filter.Field, filter.Value))
		case entity.LT:
			sb.Where(sb.L(filter.Field, filter.Value))
		case entity.NE:
			sb.Where(sb.NE(filter.Field, filter.Value))
		default:
			sb.Where(sb.In(filter.Field, filter.Values...))
		}
	}

	if request.PaginateData.Has {
		sb.Offset(request.PaginateData.Offset).Limit(request.PaginateData.Limit)
	}

	if request.OrderBy.Has {
		sb.OrderBy(request.OrderBy.Field)
		if request.OrderBy.Ascending {
			sb.Asc()
		} else {
			sb.Desc()
		}
	}

	query, args := sb.Build()
	r.logger.Info(query, args)
	rows, err := r.dbx.Queryx(query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var model entity.ChannelAccountModel
		err = rows.StructScan(&model)
		if err != nil {
			r.logger.Info(err)
			return nil, err
		}

		resp = append(resp, model.ToEntity())
	}

	return resp, nil
}

func (r repository) UpdateChannelAccounts(ctx context.Context, request entity.UpdateCollectionRequest) error {
	updateStruct := sqlbuilder.NewStruct(request.Model)
	sb := updateStruct.Update(request.TableName, &request.Model)
	log.Println(request.Filters)
	for _, filter := range request.Filters {
		switch filter.Type {
		//case domain.KeyFromDate:
		//	sb.Where(sb.GE(domain.KeyFromDate, filter.Value[0]))
		//case domain.KeyToDate:
		//	sb.Where(sb.LE(domain.KeyToDate, filter.Value[0]))
		case entity.EQ:
			sb.Where(sb.E(filter.Field, filter.Value))
		case entity.GE:
			sb.Where(sb.GE(filter.Field, filter.Value))
		case entity.LE:
			sb.Where(sb.LE(filter.Field, filter.Value))
		case entity.GT:
			sb.Where(sb.G(filter.Field, filter.Value))
		case entity.LT:
			sb.Where(sb.L(filter.Field, filter.Value))
		case entity.NE:
			sb.Where(sb.NE(filter.Field, filter.Value))
		default:
			sb.Where(sb.In(filter.Field, filter.Values...))
		}
	}

	query, args := sb.Build()
	r.logger.Info(query, args)
	// execute
	res, err := r.dbx.Exec(query, args...)
	if err != nil {
		return err
	}

	a, _ := res.RowsAffected()

	fmt.Println(a)
	return nil
}

func (r repository) DeleteAccountsFromChannel(ctx context.Context, request entity.DeleteCollectionRequest) error {
	sb := sqlbuilder.DeleteFrom(request.TableName)
	for _, filter := range request.Filters {
		switch filter.Type {
		//case domain.KeyFromDate:
		//	sb.Where(sb.GE(domain.KeyFromDate, filter.Value[0]))
		//case domain.KeyToDate:
		//	sb.Where(sb.LE(domain.KeyToDate, filter.Value[0]))
		case entity.EQ:
			sb.Where(sb.E(filter.Field, filter.Value))
		case entity.GE:
			sb.Where(sb.GE(filter.Field, filter.Value))
		case entity.LE:
			sb.Where(sb.LE(filter.Field, filter.Value))
		case entity.GT:
			sb.Where(sb.G(filter.Field, filter.Value))
		case entity.LT:
			sb.Where(sb.L(filter.Field, filter.Value))
		case entity.NE:
			sb.Where(sb.NE(filter.Field, filter.Value))
		default:
			sb.Where(sb.In(filter.Field, filter.Values...))
		}
	}

	query, args := sb.Build()
	r.logger.Info(query)
	res, err := r.dbx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		r.logger.Printf("Error %s when finding rows affected", err)
		return err
	}

	r.logger.Printf("%d permission(service) deleted ", rows)
	id, err := res.LastInsertId()
	r.logger.Printf("permission(service) deleted with id %d", id)

	return nil
}
