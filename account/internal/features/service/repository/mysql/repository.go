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
	logger *log.Logger
}

func New(db *sql.DB, logger *log.Logger) domain.ServiceRepository {
	return &repository{
		db:     db,
		dbx:    sqlx.NewDb(db, "mysql"),
		logger: logger,
	}
}

func (r *repository) CreateService(ctx context.Context, req entity.CollectionRequest) (int64, error) {
	create := sqlbuilder.NewStruct(req.Model)
	sb := create.InsertInto(req.TableName, &req.Model)
	query, args := sb.Build()
	r.logger.Info(query, args)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	res, err := r.dbx.ExecContext(ctx, query, args...)
	if err != nil {
		r.logger.Printf("Error %s when inserting row into order table", err)
		return 0, utils.NewMysqlInternalServerError(err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		r.logger.Printf("Error %s when finding rows affected", err)
		return 0, utils.NewMysqlInternalServerError(err)
	}

	r.logger.Printf("%d service created ", rows)
	id, _ := res.LastInsertId()

	return id, nil
}

func (r *repository) ReadOneService(ctx context.Context, request entity.CollectionRequest) (entity.ServiceResponse, error) {
	selectBuilder := sqlbuilder.NewStruct(entity.ServiceResponse{})
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
		return entity.ServiceResponse{}, utils.NewMysqlInternalServerError(err)
	}

	defer rows.Close()
	for rows.Next() {
		var serviceModel entity.ServiceModel
		err = rows.StructScan(&serviceModel)
		if err != nil {
			r.logger.Info(err)
			return entity.ServiceResponse{}, utils.NewMysqlInternalServerError(err)
		}
		return mapToResponse(serviceModel), nil
	}

	return entity.ServiceResponse{}, utils.NewNotFoundError("service", "filters", request.Filters)
}

func (r *repository) ReadManyServices(ctx context.Context, request entity.CollectionRequest) ([]entity.ServiceResponse, error) {
	var resp []entity.ServiceResponse
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
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var service entity.ServiceModel
		err = rows.StructScan(&service)
		if err != nil {
			r.logger.Info(err)
			return nil, err
		}

		resp = append(resp, mapToResponse(service))
	}

	return resp, nil //utils.NewNotFoundError("service not found")
}

func (r *repository) UpdateService(ctx context.Context, request entity.UpdateCollectionRequest) error {
	log.Println("-------------------------------- UpsertService Repository ---------------------------------------")
	defer log.Println("-------------------------------- UpsertService Repository ---------------------------------------")
	updateStruct := sqlbuilder.NewStruct(request.Model)
	sb := updateStruct.Update(request.TableName, &request.Model)
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
	result, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}

	r.logger.Info(result.LastInsertId())
	r.logger.Info(result.RowsAffected())
	return nil
}

func (r repository) TotalServices(ctx context.Context, request entity.CollectionRequest) (int64, error) {
	sb := sqlbuilder.Select("COUNT(*)").From(request.TableName)
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
		return 0, err
	}

	defer rows.Close()
	// iterate over rows
	var total int64
	for rows.Next() {
		err = rows.Scan(
			&total)

		if err != nil {
			r.logger.Info(err)
			return 0, err
		}

		return total, nil
	}

	return 0, nil
}

func mapToResponse(model entity.ServiceModel) entity.ServiceResponse {
	return entity.ServiceResponse{
		ID:          model.ID.Int64,
		Name:        model.Name.String,
		Path:        model.Path.String,
		Function:    model.Function.String,
		Method:      model.Method.String,
		IsActive:    model.IsActive.Bool,
		Description: model.Description.String,
		CreatedAt:   model.CreatedAt.Time,
	}
}
