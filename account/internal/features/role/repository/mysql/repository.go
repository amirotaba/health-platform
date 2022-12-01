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

func New(db *sql.DB, logger *log.Logger) domain.RoleRepository {
	return &repository{db: db, dbx: sqlx.NewDb(db, "mysql"), logger: logger}
}

func (r repository) CreateRole(ctx context.Context, request entity.InsertCollectionRequest) (id int64, err error) {
	create := sqlbuilder.NewStruct(request.Model)
	sb := create.InsertInto(domain.RoleTable, &request.Model)
	query, args := sb.Build()
	r.logger.Info(query, args)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	res, err := r.dbx.ExecContext(ctx, query, args...)
	if err != nil {
		r.logger.Printf("Error %s when inserting row into order table", err)
		return 0, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		r.logger.Printf("Error %s when finding rows affected", err)
		return 0, err
	}

	r.logger.Printf("%d token created ", rows)
	id, _ = res.LastInsertId()

	return id, nil
}

func (r repository) ReadRole(ctx context.Context, id int64) (entity.RoleEntity, error) {
	// todo : use join and left join for filling role permission
	// query role data
	/*
		select role.name, p.name from role
		         left join role_permission rp on role.id = rp.role_id
		         left join permission p on rp.permission_id = p.id
		          where role.id=1;
	*/
	rows, err := r.db.Query("select * from role where id=?;", id)

	if err != nil {
		return entity.RoleEntity{}, err
	}

	// declare empty role variable
	var role entity.RoleModel

	defer rows.Close()
	// iterate over rows
	for rows.Next() {
		err = rows.Scan(
			&role.ID,
			&role.Name,
			&role.IsActive,
			&role.Description,
			&role.CreatedAt,
			&role.UpdatedAt,
			&role.DeletedAt)

		if err != nil {
			return entity.RoleEntity{}, err
		}
		return modelToEntity(role), nil
	}

	return entity.RoleEntity{}, errors.New(fmt.Sprintf("role with id %d not found", id))
}

func (r repository) ReadRoles(ctx context.Context) ([]entity.RoleEntity, error) {
	// query role data
	log.Printf("Read roles started \n")
	rows, err := r.db.Query("select * from role;")

	if err != nil {
		return []entity.RoleEntity{}, err
	}

	defer rows.Close()
	// declare empty role variable
	var roles []entity.RoleEntity
	// iterate over rows
	for rows.Next() {
		var role entity.RoleModel

		err = rows.Scan(
			&role.ID,
			&role.Name,
			&role.IsActive,
			&role.Description,
			&role.CreatedAt,
			&role.UpdatedAt,
			&role.DeletedAt)

		if err != nil {
			return []entity.RoleEntity{}, err
		}
		roles = append(roles, modelToEntity(role))
	}

	return roles, nil
}

func (r repository) ReadOneRole(ctx context.Context, req entity.CollectionRequest) (entity.RoleEntity, error) {
	fmt.Println("*********************ReadOneRole*****************************************")
	defer fmt.Println("*********************ReadOneRole*****************************************")
	var resp entity.RoleEntity
	selectBuilder := sqlbuilder.NewStruct(entity.RoleModel{})
	sb := selectBuilder.SelectFrom(domain.RoleTable)
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
		var model entity.RoleModel
		err = rows.StructScan(&model)
		if err != nil {
			r.logger.Info(err)
			return resp, err
		}

		return modelToEntity(model), err
	}

	return resp, utils.NewNotFoundError("role", "filters", req.Filters)

}

func (r repository) ReadManyRole(ctx context.Context, req entity.CollectionRequest) ([]entity.RoleEntity, error) {
	fmt.Println("*********************ReadManyRole*****************************************")
	defer fmt.Println("*********************ReadManyRole*****************************************")
	var resp []entity.RoleEntity
	log.Println(req.Filters)
	selectBuilder := sqlbuilder.NewStruct(entity.RoleModel{})
	sb := selectBuilder.SelectFrom(domain.RoleTable)
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

	var searches []string
	for _, search := range req.Searches {
		searches = append(searches, sb.Like(search.Field, search.Value))
	}

	if len(searches) > 0 {
		sb.Where(searches...)
	}

	query, args := sb.Build()
	r.logger.Info(query, args)
	rows, err := r.dbx.Queryx(query, args...)
	if err != nil {
		return resp, err
	}

	defer rows.Close()
	for rows.Next() {
		var model entity.RoleModel
		err = rows.StructScan(&model)
		if err != nil {
			r.logger.Info(err)
			return resp, err
		}

		resp = append(resp, modelToEntity(model))
	}

	if len(resp) == 0 {
		return resp, utils.NewNotFoundError("role", "filters", req.Filters)
	}

	return resp, nil //utils.NewNotFoundError("role", "filters", req.Filters)
}

func (r repository) ReadRolePermissions(ctx context.Context, roleId int64) ([]entity.PermissionEntity, error) {
	// query role data
	log.Printf("Read roles started \n")
	query := `select 
                     rp.permission_id,
                     p.name,
                     p.is_active,
                     p.description,
                     p.created_at,
                     p.updated_at,
                     p.deleted_at
              from role_permission rp     
                  left join role r on r.id = rp.role_id
                  left join permission p on rp.permission_id = p.id
              where rp.role_id=?;`
	rows, err := r.db.Query(query, roleId)

	if err != nil {
		return []entity.PermissionEntity{}, err
	}

	defer rows.Close()
	// declare empty role variable
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

		permissions = append(permissions, permissionModelToEntity(permission))
	}

	return permissions, nil
}

func (r repository) UpdateRole(ctx context.Context, role entity.RoleUpdateRequest) error {
	stmt, err := r.db.Prepare("update role set name=?, is_active=?, description=?, updated_at=? where id=?")
	if err != nil {
		return err
	}

	defer stmt.Close()
	// execute
	res, err := stmt.Exec(role.Name, role.IsActive, role.Description, time.Now(), role.ID)
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

func (r repository) DeleteRole(ctx context.Context, request entity.DeleteCollectionRequest) error {
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

//
//func (r repository) UpdatePermissionsToRole(ctx context.Context, permissions entity.PermissionsToRoleRequest) error {
//	query := "INSERT INTO role_permission(role_id, permission_id) VALUES "
//	var values []interface{}
//
//	for _, row := range permissions.Permissions {
//		query += "(?, ?),"
//		values = append(values, permissions.RoleId, row)
//	}
//
//	//trim the last ,
//	query = query[0 : len(query)-1]
//	ctx, cancelFunc := context.WithTimeout(context.Background(), 25*time.Second)
//
//	defer cancelFunc()
//	stmt, err := r.db.PrepareContext(ctx, query)
//	if err != nil {
//		log.Printf("Error %s when preparing SQL statement", err)
//		return err
//	}
//
//	defer stmt.Close()
//	res, err := stmt.ExecContext(ctx, values...)
//	if err != nil {
//		log.Printf("Error %s when inserting row into role table", err)
//		return err
//	}
//
//	rows, err := res.RowsAffected()
//	if err != nil {
//		log.Printf("Error %s when finding rows affected", err)
//		return err
//	}
//
//	log.Printf("%d role permission updated created ", rows)
//	_, err = res.LastInsertId()
//
//	return nil
//}
//
//func (r repository) DeletePermissionsFromRole(ctx context.Context, permissions entity.PermissionsToRoleRequest) error {
//	query := fmt.Sprintf("DELETE FROM role_permission WHERE role_id=%d AND permission_id IN ", permissions.RoleId)
//	query += "("
//	for permission, _ := range permissions.Permissions {
//		query += fmt.Sprintf("%d,", permission)
//	}
//
//	//trim the last ,
//	query = query[0 : len(query)-1]
//	query += ")"
//
//	ctx, cancelFunc := context.WithTimeout(context.Background(), 25*time.Second)
//
//	defer cancelFunc()
//	stmt, err := r.db.PrepareContext(ctx, query)
//	if err != nil {
//		log.Printf("Error %s when preparing SQL statement", err)
//		return err
//	}
//
//	defer stmt.Close()
//	res, err := stmt.ExecContext(ctx)
//	if err != nil {
//		log.Printf("Error %s when inserting row into role table", err)
//		return err
//	}
//
//	rows, err := res.RowsAffected()
//	if err != nil {
//		log.Printf("Error %s when finding rows affected", err)
//		return err
//	}
//
//	log.Printf("%d role permission updated created ", rows)
//	_, err = res.LastInsertId()
//
//	return nil
//}

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

func permissionModelToEntity(model entity.PermissionModel) entity.PermissionEntity {
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
