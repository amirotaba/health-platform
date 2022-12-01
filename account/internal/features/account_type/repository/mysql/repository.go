package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"git.paygear.ir/giftino/account/internal/account/entity"
	"time"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"

	"git.paygear.ir/giftino/account/internal/account/domain"
	"git.paygear.ir/giftino/account/internal/utils"
)

type repository struct {
	db     *sql.DB
	dbx    *sqlx.DB
	logger *log.Logger
}

func New(db *sql.DB, logger *log.Logger) domain.AccountTypeRepository {
	return &repository{
		db:     db,
		dbx:    sqlx.NewDb(db, "mysql"),
		logger: logger,
	}
}

func (r repository) CreateAccountType(ctx context.Context, accountType *domain.AccountTypeEntity) error {
	query := `INSERT INTO account_type ( name, is_active, description) VALUES (?,?,?);`
	ctx, cancelFunc := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancelFunc()
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement \n", err)
		return err
	}

	defer stmt.Close()
	res, err := stmt.ExecContext(ctx,
		accountType.Name,
		accountType.IsActive,
		accountType.Description)
	if err != nil {
		log.Printf("Error %s when inserting row into accountType table \n", err)
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected \n", err)
		return err
	}

	log.Printf("%d accountType created \n", rows)
	_, err = res.LastInsertId()

	return nil
}

func (r repository) ReadAccountType(ctx context.Context, id int64) (domain.AccountTypeEntity, error) {
	// query accountType data
	rows, err := r.db.Query("select * from account_type where id=?;", id)
	if err != nil {
		return domain.AccountTypeEntity{}, err
	}

	defer rows.Close()
	// declare empty accountType variable
	// iterate over rows
	for rows.Next() {
		var accountType domain.AccountTypeModel
		err = rows.Scan(
			&accountType.ID,
			&accountType.Name,
			&accountType.IsActive,
			&accountType.Description,
			&accountType.CreatedAt,
			&accountType.UpdatedAt,
			&accountType.DeletedAt)

		if err != nil {
			return domain.AccountTypeEntity{}, err
		}

		return modelToEntity(accountType), nil
	}

	return domain.AccountTypeEntity{}, errors.New("not found")
}

func (r repository) ReadAccountTypes(ctx context.Context) ([]domain.AccountTypeEntity, error) {
	// query accountType data
	log.Printf("Read accountTypes started \n")
	rows, err := r.db.Query("select * from account_type;")
	if err != nil {
		return []domain.AccountTypeEntity{}, err
	}

	defer rows.Close()
	// declare empty accountType variable
	var accountTypes []domain.AccountTypeEntity
	// iterate over rows
	for rows.Next() {
		var accountType domain.AccountTypeModel

		err = rows.Scan(
			&accountType.ID,
			&accountType.Name,
			&accountType.IsActive,
			&accountType.Description,
			&accountType.CreatedAt,
			&accountType.UpdatedAt,
			&accountType.DeletedAt)

		if err != nil {
			return []domain.AccountTypeEntity{}, err
		}
		accountTypes = append(accountTypes, modelToEntity(accountType))
	}

	return accountTypes, nil
}

func (r repository) ReadOneAccountType(ctx context.Context, req entity.CollectionRequest) (domain.AccountTypeEntity, error) {
	var resp domain.AccountTypeEntity
	selectBuilder := sqlbuilder.NewStruct(domain.AccountTypeModel{})
	sb := selectBuilder.SelectFrom(domain.AccountTypeTable)
	for _, filter := range req.Filters {
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
		return resp, err
	}

	defer rows.Close()
	for rows.Next() {
		var model domain.AccountTypeModel
		err = rows.StructScan(&model)
		if err != nil {
			r.logger.Info(err)
			return resp, err
		}

		return modelToEntity(model), err
	}

	return resp, utils.NewNotFoundError("account type", "filters", req.Filters)
}

func (r repository) UpdateAccountType(ctx context.Context, accountType *domain.AccountTypeUpdateRequest) error {
	query := `update account_type set name=?, is_active=?, description=?, updated_at=? where id=?`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()
	// execute
	res, err := stmt.Exec(accountType.Name, accountType.IsActive, accountType.Description, time.Now(), accountType.ID)
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

func modelToEntity(model domain.AccountTypeModel) domain.AccountTypeEntity {
	return domain.AccountTypeEntity{
		ID:          model.ID.Int64,
		Name:        model.Name.String,
		IsActive:    model.IsActive.Bool,
		Description: model.Description.String,
		CreatedAt:   model.CreatedAt.Time,
		UpdatedAt:   model.UpdatedAt.Time,
		DeletedAt:   model.DeletedAt.Time,
	}
}
