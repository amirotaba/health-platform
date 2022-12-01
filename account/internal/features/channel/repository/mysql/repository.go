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
	"git.paygear.ir/giftino/account/internal/utils"
)

type repository struct {
	db     *sql.DB
	dbx    *sqlx.DB
	logger *log.Logger
}

func New(db *sql.DB, logger *log.Logger) domain.ChannelRepository {
	return &repository{
		db:     db,
		dbx:    sqlx.NewDb(db, "mysql"),
		logger: logger,
	}
}

func (r repository) CreateChannel(ctx context.Context, req entity.InsertCollectionRequest) (id int64, err error) {
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

	r.logger.Printf("%d channel created succesfully created ", rows)
	id, _ = res.LastInsertId()
	return
}

func (r repository) ReadOneChannel(ctx context.Context, request entity.CollectionRequest) (entity.ChannelEntity, error) {
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

	query, args := sb.Build()
	r.logger.Info(query, args)
	rows, err := r.dbx.Queryx(query, args...)
	if err != nil {
		return entity.ChannelEntity{}, err
	}

	defer rows.Close()
	var account entity.ChannelModel
	for rows.Next() {
		err = rows.StructScan(&account)
		if err != nil {
			r.logger.Info(err)
			return entity.ChannelEntity{}, err
		}

		return account.ToEntity(), nil
	}

	return entity.ChannelEntity{}, utils.NewNotFoundError("channel", "filters", request.Filters)
}

func (r repository) ReadManyChannels(ctx context.Context, request entity.CollectionRequest) ([]entity.ChannelEntity, error) {
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

	var searches []string
	for _, search := range request.Searches {
		searches = append(searches, sb.Like(search.Field, search.Value))
	}

	if len(searches) > 0 {
		sb.Where(searches...)
	}

	//parent := request.JoinWithOptions.RootTable

	for _, join := range request.JoinWithOptions.Joins {
		sb.JoinWithOption(sqlbuilder.JoinOption(join.JoinType), join.LeftTable, fmt.Sprintf("%s.%s=%s.%s", join.LeftTable, join.ON.LeftON, join.RightTable, join.ON.RightON))
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

	//for _, groupBy := range request.GroupsBy {
	//	sb.GroupBy(fmt.Sprintf("%s.%s", groupBy.TableName, groupBy.Field))
	//}

	query, args := sb.Build()
	r.logger.Info(query, args)
	rows, err := r.dbx.Queryx(query, args...)
	if err != nil {
		return nil, err
	}

	var resp []entity.ChannelEntity
	defer rows.Close()
	for rows.Next() {
		var model entity.ChannelModel
		err = rows.StructScan(&model)
		if err != nil {
			r.logger.Info(err)
			return nil, err
		}

		resp = append(resp, model.ToEntity())
	}

	return resp, nil
}

func (r repository) ReadTotal(ctx context.Context, request entity.CollectionRequest) entity.CollectionRequest {
	fmt.Println("***************************************ReadOne********************")
	defer fmt.Println("***************************************ReadOne********************")
	var sb *sqlbuilder.SelectBuilder
	if request.Sql == "" {
		selectBuilder := sqlbuilder.NewStruct(request.Model)
		sb = selectBuilder.SelectFrom(request.TableName)
	} else {
		sb = sqlbuilder.Select(request.Sql).From(request.TableName)
	}

	for _, filter := range request.Filters {
		r.logger.Info(filter.Value)
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
		default:
			sb.Where(sb.In(filter.Field, filter.Values...))
		}
	}

	var searches []string
	for _, search := range request.Searches {
		searches = append(searches, sb.Like(search.Field, search.Value))
	}

	if len(searches) > 0 {
		sb.Where(searches...)
	}

	//if request.PaginateData.Has {
	//	sb.Offset(request.PaginateData.Offset).Limit(request.PaginateData.Limit)
	//}

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
	//var model entity.FactorModel
	rows, err := r.dbx.Queryx(query, args...)
	if err != nil {
		r.logger.Info(err.Error())
		return request
	}

	//defer rows.Close()
	//for rows.Next() {
	//	err = rows.StructScan(&model)
	//	if err != nil {
	//		r.logger.Info(err)
	//		return request, utils.NewMysqlInternalServerError(err)
	//	}
	//
	//	return request, nil
	//}

	request.Result = rows
	//return request, utils.NewNotFoundError("Factor", "id", request.Filters)
	return request
}

func (r repository) UpdateChannelNew(ctx context.Context, request entity.UpdateCollectionRequest) error {
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
