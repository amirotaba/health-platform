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

func New(db *sql.DB, logger *log.Logger) domain.AccountRepository {
	return &repository{db: db, dbx: sqlx.NewDb(db, "mysql"), logger: logger}
}

func (r repository) CreateAccountNew(ctx context.Context, request entity.InsertCollectionRequest) (int64, error) {
	create := sqlbuilder.NewStruct(request.Model)
	sb := create.InsertInto(domain.AccountTable, &request.Model)
	query, args := sb.Build()
	r.logger.Info(query, args)
	res, err := r.dbx.Exec(query, args...)
	if err != nil {
		r.logger.Printf("Error %s when inserting row into account table", err)
		return 0, utils.NewMysqlInternalServerError(err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		r.logger.Printf("Error %s when finding rows affected", err)
		return 0, utils.NewMysqlInternalServerError(err)
	}

	r.logger.Printf("%d new account created ", rows)
	id, _ := res.LastInsertId()

	return id, nil
}

//
//func (r repository) ReadAccount(ctx context.Context, id int64) (entity.AccountEntity, error) {
//	// query account data
//	rows, err := r.db.Query("select * from account where id=?;", id)
//	if err != nil {
//		return entity.AccountEntity{}, err
//	}
//
//	defer rows.Close()
//	// declare empty account variable
//	var account entity.AccountModel
//	// iterate over rows
//	for rows.Next() {
//		err = rows.Scan(
//			&account.ID,
//			&account.UUID,
//			&account.FirstName,
//			&account.LastName,
//			&account.DisplayName,
//			&account.Password,
//			&account.Email,
//			&account.PhoneNumber,
//			&account.Address,
//			&account.TypeId,
//			&account.IsActive,
//			&account.ExpireAt,
//			&account.LastLogin,
//			&account.Description,
//			&account.CreatedAt,
//			&account.UpdatedAt,
//			&account.DeletedAt)
//
//		if err != nil {
//			return entity.AccountEntity{}, err
//		}
//		return account.ToEntity(), nil
//	}
//
//	return entity.AccountEntity{}, utils.NewNotFoundError("account", "id", id) //errors.New(fmt.Sprintf("account with id %v not found", id))
//}

func (r repository) ReadOneAccount(ctx context.Context, request entity.CollectionRequest) (entity.AccountEntity, error) {
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
		return entity.AccountEntity{}, err
	}

	defer rows.Close()
	var account entity.AccountModel
	for rows.Next() {
		err = rows.StructScan(&account)
		if err != nil {
			r.logger.Info(err)
			return entity.AccountEntity{}, err
		}

		return account.ToEntity(), nil
	}

	return entity.AccountEntity{}, utils.NewNotFoundError("account", "filters", request.Filters)
}

func (r repository) ReadAccountByOwnerPhoneNumber(ctx context.Context, phoneNumber string) (entity.AccountEntity, error) {
	// query account data
	log.Println("ReadAccountByOwnerPhoneNumber")
	log.Println(fmt.Sprintf("select * from account where phone_number=%s;", phoneNumber))
	rows, err := r.db.Query("select * from account where phone_number=?;", phoneNumber)
	if err != nil {
		return entity.AccountEntity{}, err
	}

	defer rows.Close()
	var account entity.AccountModel
	// iterate over rows
	for rows.Next() {
		err = rows.Scan(
			&account.ID,
			&account.UUID,
			&account.FirstName,
			&account.LastName,
			&account.DisplayName,
			&account.Password,
			&account.Email,
			&account.PhoneNumber,
			&account.Address,
			&account.ImageUrl,
			&account.TypeId,
			&account.IsActive,
			&account.ExpireAt,
			&account.LastLogin,
			&account.Description,
			&account.CreatedAt,
			&account.UpdatedAt,
			&account.DeletedAt)

		if err != nil {
			return entity.AccountEntity{}, err
		}
		return account.ToEntity(), nil
	}

	return entity.AccountEntity{}, errors.New("not found")
}

func (r repository) ReadAccounts(ctx context.Context) ([]entity.AccountEntity, error) {
	// query account data
	rows, err := r.db.Query("select * from account;")
	if err != nil {
		return []entity.AccountEntity{}, err
	}

	defer rows.Close()
	// declare empty account variable
	var accounts []entity.AccountEntity
	// iterate over rows
	for rows.Next() {
		var account entity.AccountModel
		err = rows.Scan(
			&account.ID,
			&account.UUID,
			&account.FirstName,
			&account.LastName,
			&account.DisplayName,
			&account.Password,
			&account.Email,
			&account.PhoneNumber,
			&account.Address,
			&account.TypeId,
			&account.IsActive,
			&account.ExpireAt,
			&account.LastLogin,
			&account.Description,
			&account.CreatedAt,
			&account.UpdatedAt,
			&account.DeletedAt)

		log.Println(account)
		if err != nil {
			log.Println(err)
			return []entity.AccountEntity{}, err
		}

		accounts = append(accounts, account.ToEntity())
	}

	log.Println("list of account : ", accounts)
	return accounts, nil
}

func (r repository) ReadManyAccounts(ctx context.Context, request entity.CollectionRequest) ([]entity.AccountEntity, error) {
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

	var resp []entity.AccountEntity
	defer rows.Close()
	for rows.Next() {
		var model entity.AccountModel
		err = rows.StructScan(&model)
		if err != nil {
			r.logger.Info(err)
			return nil, err
		}

		resp = append(resp, model.ToEntity())
	}

	return resp, nil
}

func (r repository) UpdateAccount(ctx context.Context, account entity.UpdateAccountRequest) error {
	query := `update account set first_name=?, 
                                 last_name=?,
                                 display_name=?,
                                 email=?,
                                 owner_phone_number=?,
                                 address=?,
                                 is_active=?,
                                 description=?,
                                 updated_at=? where id=?`

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()
	// execute
	res, err := stmt.Exec(account.FirstName, account.LastName, account.DisplayName, account.Email,
		account.PhoneNumber, account.Address, account.IsActive, account.Description, time.Now(), account.ID)
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

func (r repository) UpdateAccountNew(ctx context.Context, request entity.UpdateCollectionRequest) error {
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

func (r repository) UpdateAccountPassword(ctx context.Context, req entity.ResetPassRequest, id int64) error {
	query := `update account set password=?, updated_at=? where id=?`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()
	// execute
	res, err := stmt.Exec(req.NewPass, time.Now(), id)
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

func (r repository) UpdateAccountLastLogin(ctx context.Context, phoneNumber string, lastLogin time.Time) error {
	query := `update account set last_login=? where phone_number=?`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()
	// execute
	res, err := stmt.Exec(lastLogin, phoneNumber)
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
