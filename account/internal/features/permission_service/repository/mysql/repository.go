package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
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

func New(db *sql.DB, logger *log.Logger) domain.PermissionServiceRepository {
	return &repository{
		db:     db,
		dbx:    sqlx.NewDb(db, "mysql"),
		logger: logger,
	}
}

func (r repository) CreatePermissionServices(ctx context.Context, req entity.InsertCollectionRequest) (id int64, err error) {
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

	r.logger.Printf("%d permission(service) created ", rows)
	id, _ = res.LastInsertId()
	return
}

func (r repository) ReadPermissionServices(ctx context.Context, request entity.CollectionRequest) ([]entity.PermissionServicesResponse, error) {
	var resp []entity.PermissionServicesResponse
	selectBuilder := sqlbuilder.NewStruct(entity.PermissionServicesResponse{})
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
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var model entity.PermissionServicesResponseModel
		err = rows.StructScan(&model)
		if err != nil {
			r.logger.Info(err)
			return nil, err
		}

		resp = append(resp, mapToResponse(model))
	}

	return resp, err
}

func (r repository) ReadPermissionServicesN(ctx context.Context, request entity.CollectionRequest) ([]entity.PermissionServicesResponse, error) {
	var resp []entity.PermissionServicesResponse
	selectBuilder := sqlbuilder.NewSelectBuilder() //sqlbuilder.NewStruct(entity.PermissionServicesResponse{})
	sb := selectBuilder.Select(CreateSelectStmt(entity.PermissionServicesResponse{}, "db")).From(request.TableName)
	//selectBuilder := sqlbuilder.NewStruct(entity.PermissionServicesResponse{})
	//sb := selectBuilder.SelectFrom(domain.PermissionServiceTable)
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

	//parent := request.JoinWithOptions.ParentTable

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

	query, args := sb.Build()
	//query = strings.ReplaceAll(query, domain.PermissionServiceTable+"."+domain.PermissionTable, domain.PermissionTable)
	//query = strings.ReplaceAll(query, domain.PermissionServiceTable+"."+domain.ServicesTable, domain.ServicesTable)
	r.logger.Info(query, args)
	rows, err := r.dbx.Queryx(query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		//model := make(map[string]interface{})
		var model entity.PermissionServicesResponseModel
		//err = rows.MapScan(model)
		err = rows.StructScan(&model)
		if err != nil {
			r.logger.Info(err)
			return nil, err
		}

		resp = append(resp, mapToResponse(model))
	}

	return resp, err
}

func (r repository) TotalPermissionServices(ctx context.Context, request entity.CollectionRequest) (int64, error) {
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

	//parent := request.JoinWithOptions.ParentTable

	for _, join := range request.JoinWithOptions.Joins {
		sb.JoinWithOption(sqlbuilder.JoinOption(join.JoinType), join.LeftTable, fmt.Sprintf("%s.%s=%s.%s", join.LeftTable, join.ON.LeftON, join.RightTable, join.ON.RightON))
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

func (r repository) DeletePermissionServices(ctx context.Context, request entity.DeleteCollectionRequest) error {
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
	r.logger.Info(query, args)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	res, err := r.dbx.ExecContext(ctx, query, args...)
	if err != nil {
		r.logger.Printf("Error %s when inserting row into permission table", err)
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		r.logger.Printf("Error %s when finding rows affected", err)
		return err
	}

	r.logger.Printf("%d permission(service) deleted ", rows)

	return nil
}

func CreateSelectStmt(i interface{}, tagName string) string {
	t := reflect.TypeOf(i)

	// Get the type and kind of our user variable
	fmt.Println("Type:", t.Name())
	fmt.Println("Kind:", t.Kind())

	var tags []string
	// Iterate over all available fields and read the tag value
	for i := 0; i < t.NumField(); i++ {
		// Get the field, returns https://golang.org/pkg/reflect/#StructField
		field := t.Field(i)

		// Get the field tag value
		tag := field.Tag.Get(tagName)

		fmt.Printf("%d. %v (%v), tag: '%v'\n", i+1, field.Name, field.Type.Name(), tag)
		// Skip if tag is not defined or ignored
		if tag == "" || tag == "-" {
			continue
		}

		tags = append(tags, tag)
	}

	return strings.Join(tags, ", ")
}

func modelToEntity(model entity.PermissionModel) entity.PermissionEntity {
	return entity.PermissionEntity{
		ID:          model.ID.Int64,
		Name:        model.Name.String,
		IsActive:    model.IsActive.Bool,
		Description: model.Description.String,
		CreatedAt:   model.CreatedAt.Time,
		UpdatedAt:   model.UpdatedAt.Time,
		DeletedAt:   model.DeletedAt.Time,
	}
}

func mapToResponse(model entity.PermissionServicesResponseModel) entity.PermissionServicesResponse {
	return entity.PermissionServicesResponse{
		ID:           model.ID.Int64,
		PermissionID: model.PermissionID.Int64,
		//PermissionName: model.PermissionName.String,
		ServiceID: model.ServiceID.Int64,
		//ServiceName:      model.ServiceName.String,
		Description: model.Description.String,
		CreatedAt:   model.CreatedAt.Time,
	}
}
