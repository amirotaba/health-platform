package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"

	"git.paygear.ir/giftino/account/internal/account/domain"
	"git.paygear.ir/giftino/account/internal/account/entity"
	"git.paygear.ir/giftino/account/internal/utils"
)

type repository struct {
	db     *sql.DB
	dbx    *sqlx.DB
	logger log.Logger
}

func New(db *sql.DB) domain.OtpRepository {
	return &repository{
		db:  db,
		dbx: sqlx.NewDb(db, "mysql"),
	}
}

func (r *repository) CreateOtp(ctx context.Context, in entity.OtpRequest) (int64, error) {
	create := sqlbuilder.NewStruct(entity.OtpRequest{})
	sb := create.InsertInto(domain.AccountOtpTable, &in)
	query, args := sb.Build()
	r.logger.Info(query, args)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	res, err := r.dbx.ExecContext(ctx, query, args...)
	if err != nil {
		r.logger.Printf("Error %s when inserting row into otp table", err)
		return 0, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		r.logger.Printf("Error %s when finding rows affected", err)
		return 0, err
	}

	r.logger.Printf("%d otp created ", rows)
	id, _ := res.LastInsertId()

	return id, nil
}

func (r *repository) ReadOtp(ctx context.Context, request entity.CollectionRequest, result interface{}) error {
	selectBuilder := sqlbuilder.NewStruct(request.Model)
	sb := selectBuilder.SelectFrom(request.TableName)
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
		return err
	}

	defer rows.Close()
	for rows.Next() {
		err = rows.StructScan(result)
		if err != nil {
			r.logger.Info(err)
			return err
		}

		return nil
	}

	return utils.NewNotFoundError("otp", "filters", request.Filters)
}
