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

func New(db *sql.DB, logger *log.Logger) domain.AccountRoleRepository {
	return &repository{
		db:     db,
		dbx:    sqlx.NewDb(db, "mysql"),
		logger: logger,
	}
}

func (r repository) CreateAccountRoles(ctx context.Context, req entity.InsertCollectionRequest) (id int64, err error) {
	create := sqlbuilder.NewStruct(req.Model)
	sb := create.InsertInto(req.TableName, &req.Model)
	query, args := sb.Build()
	r.logger.Info(query, args)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	res, err := r.dbx.ExecContext(ctx, query, args...)
	if err != nil {
		r.logger.Printf("Error %s when inserting row into role table", err)
		return
	}

	rows, err := res.RowsAffected()
	if err != nil {
		r.logger.Printf("Error %s when finding rows affected", err)
		return
	}

	r.logger.Printf("%d role(service) created ", rows)
	id, _ = res.LastInsertId()
	return
}

func (r repository) ReadAccountRoles(ctx context.Context, request entity.CollectionRequest) ([]entity.AccountRolesResponse, error) {
	var resp []entity.AccountRolesResponse
	selectBuilder := sqlbuilder.NewStruct(entity.AccountRolesResponse{})
	sb := selectBuilder.SelectFrom(domain.AccountRoleTable)
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
		var model entity.AccountRoleModel
		err = rows.StructScan(&model)
		if err != nil {
			r.logger.Info(err)
			return nil, err
		}

		resp = append(resp, mapToResponse(model))
	}

	return resp, err
}

func (r repository) DeleteRolesFromAccount(ctx context.Context, roles entity.CreateRolesToAccountRequest) error {
	query := fmt.Sprintf("DELETE FROM account_role WHERE account_id=%d AND role_id IN ", roles.AccountId)
	query += "("
	for role, _ := range roles.Roles {
		query += fmt.Sprintf("%d,", role)
	}

	//trim the last ,
	query = query[0 : len(query)-1]
	query += ")"

	ctx, cancelFunc := context.WithTimeout(context.Background(), 25*time.Second)

	defer cancelFunc()
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement \n", err)
		return err
	}

	defer stmt.Close()
	res, err := stmt.ExecContext(ctx)
	if err != nil {
		log.Printf("Error %s when inserting row into account table \n", err)
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected \n", err)
		return err
	}

	log.Printf("%d account role updated created \n", rows)
	_, err = res.LastInsertId()

	return nil
}

func (r repository) UpdateAccountRoles(ctx context.Context, roles entity.CreateRolesToAccountRequest) error {
	query := "INSERT INTO account_role(account_id, role_id) VALUES "
	var values []interface{}

	for _, row := range roles.Roles {
		query += "(?, ?),"
		values = append(values, roles.AccountId, row)
	}

	//trim the last ,
	query = query[0 : len(query)-1]
	ctx, cancelFunc := context.WithTimeout(context.Background(), 25*time.Second)

	defer cancelFunc()
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement \n", err)
		return err
	}

	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, values...)
	if err != nil {
		log.Printf("Error %s when inserting row into account table \n", err)
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected \n", err)
		return err
	}

	log.Printf("%d account role updated created \n", rows)
	_, err = res.LastInsertId()

	return nil
}

func (r repository) UpdateAccountRolesN(ctx context.Context, request entity.CollectionRequest) error {
	return nil
}

func (r repository) DeleteRolesFromAccountN(ctx context.Context, request entity.DeleteCollectionRequest) error {
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

	r.logger.Printf("%d role(service) deleted ", rows)
	id, err := res.LastInsertId()
	r.logger.Printf("role(service) deleted with id %d", id)

	return nil
}

func (r repository) TotalAccountRoles(ctx context.Context, request entity.CollectionRequest) (int64, error) {
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

func modelToEntity(model entity.RoleModel) entity.RoleEntity {
	return entity.RoleEntity{
		ID:          model.ID.Int64,
		Name:        model.Name.String,
		IsActive:    model.IsActive.Bool,
		Description: model.Description.String,
		CreatedAt:   model.CreatedAt.Time,
		UpdatedAt:   model.UpdatedAt.Time,
		DeletedAt:   model.DeletedAt.Time,
	}
}

func mapToResponse(model entity.AccountRoleModel) entity.AccountRolesResponse {
	return entity.AccountRolesResponse{
		ID:          model.ID.Int64,
		RoleID:      model.RoleID.Int64,
		AccountID:   model.AccountID.Int64,
		Description: model.Description.String,
		CreatedAt:   model.CreatedAt.Time,
	}
}
