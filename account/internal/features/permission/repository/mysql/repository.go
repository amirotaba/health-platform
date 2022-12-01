package mysql

import (
	"context"
	"database/sql"
	"errors"
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

func New(db *sql.DB, logger *log.Logger) domain.PermissionRepository {
	return &repository{
		db:     db,
		dbx:    sqlx.NewDb(db, "mysql"),
		logger: logger,
	}
}

func (r repository) CreatePermission(ctx context.Context, req entity.CollectionRequest) (id int64, err error) {
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

func (r repository) ReadPermission(ctx context.Context, id int64) (entity.PermissionEntity, error) {
	// query permission data
	rows, err := r.db.Query("select * from permission where id=?;", id)

	if err != nil {
		return entity.PermissionEntity{}, err
	}

	defer rows.Close()
	// declare empty permission variable
	var permissions []entity.PermissionEntity
	// iterate over rows
	for rows.Next() {
		var permission entity.PermissionModel
		err = rows.Scan(
			&permission.ID,
			&permission.Name,
			&permission.IsActive,
			&permission.Description,
			&permission.CreatedAt,
			&permission.UpdatedAt,
			&permission.DeletedAt)

		permissions = append(permissions, modelToEntity(permission))
		if err != nil {
			return entity.PermissionEntity{}, err
		}
	}

	if len(permissions) < 1 {
		return entity.PermissionEntity{}, errors.New("permission not found")
	}
	return permissions[0], nil
}

func (r repository) ReadManyPermission(ctx context.Context, request entity.CollectionRequest) ([]entity.PermissionEntity, error) {
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

	var resp []entity.PermissionEntity
	defer rows.Close()
	for rows.Next() {
		var model entity.PermissionModel
		err = rows.StructScan(&model)
		if err != nil {
			r.logger.Info(err)
			return nil, err
		}

		resp = append(resp, model.ToEntity())
	}

	return resp, nil
}

func (r repository) NewReadPermission(ctx context.Context, request entity.CollectionRequest, result interface{}) error {
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

	return utils.NewNotFoundError("permission", "filters", request.Filters)
}

func (r repository) ReadPermissions(ctx context.Context) ([]entity.PermissionEntity, error) {
	// query permission data
	rows, err := r.db.Query("select * from permissions;")

	if err != nil {
		return []entity.PermissionEntity{}, err
	}

	defer rows.Close()
	// declare empty permission variable
	var permissions []entity.PermissionEntity
	// iterate over rows
	for rows.Next() {
		var permission entity.PermissionModel
		err = rows.Scan(
			&permission.ID,
			&permission.Name,
			&permission.IsActive,
			&permission.Description,
			&permission.CreatedAt,
			&permission.UpdatedAt,
			&permission.DeletedAt)

		if err != nil {
			return []entity.PermissionEntity{}, err
		}

		permissions = append(permissions, modelToEntity(permission))
	}

	return permissions, nil
}

func (r repository) ReadPermissionServices(ctx context.Context, request entity.CollectionRequest, permissionID int64) ([]entity.PermissionServicesResponse, error) {
	var resp []entity.PermissionServicesResponse
	selectBuilder := sqlbuilder.NewStruct(entity.PermissionServicesResponse{})
	sb := selectBuilder.SelectFrom(domain.PermissionServiceTable)
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
		var model entity.PermissionServicesModel
		err = rows.StructScan(&model)
		if err != nil {
			r.logger.Info(err)
			return nil, err
		}

		resp = append(resp, mapToResponse(model))
	}

	return resp, err
}

func (r repository) UpdatePermission(ctx context.Context, permission entity.PermissionUpdateRequest) error {
	stmt, err := r.db.Prepare("update permission set name=?, is_active=?, description=?, updated_at=? where id=?")
	if err != nil {
		return err
	}

	defer stmt.Close()
	// execute
	res, err := stmt.Exec(permission.Name, permission.IsActive, permission.Description, time.Now(), permission.ID)
	if err != nil {
		return err
	}

	a, err := res.RowsAffected()
	if err != nil {
		return err
	}

	fmt.Println(a)
	return nil
}

func (r repository) DeletePermission(ctx context.Context, request entity.CollectionRequest) error {
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

func mapToResponse(model entity.PermissionServicesModel) entity.PermissionServicesResponse {
	return entity.PermissionServicesResponse{
		ID:           model.ID.Int64,
		PermissionID: model.PermissionID.Int64,
		ServiceID:    model.ServiceID.Int64,
		Description:  model.Description.String,
		CreatedAt:    model.CreatedAt.Time,
	}
}
