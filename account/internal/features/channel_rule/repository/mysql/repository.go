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

func New(db *sql.DB, logger *log.Logger) domain.ChannelRuleRepository {
	return &repository{
		db:     db,
		dbx:    sqlx.NewDb(db, "mysql"),
		logger: logger,
	}
}

func (r repository) CreateChannelRule(ctx context.Context, req entity.InsertCollectionRequest) (id int64, err error) {
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

	r.logger.Printf("%d channel rule created ", rows)
	id, _ = res.LastInsertId()
	return
}

func (r repository) ReadChannelRule(ctx context.Context, id int64) (entity.ChannelRuleEntity, error) {
	// query channelRule data
	rows, err := r.db.Query("select * from channel_rule where id=?;", id)

	if err != nil {
		return entity.ChannelRuleEntity{}, err
	}

	defer rows.Close()
	// declare empty channelRule variable
	var channelRule entity.ChannelRuleModel
	// iterate over rows
	for rows.Next() {
		err = rows.Scan(
			&channelRule.ID,
			&channelRule.ChannelId,
			&channelRule.TagId,
			&channelRule.Price,
			&channelRule.IsActive,
			&channelRule.Description,
			&channelRule.CreatedAt,
			&channelRule.UpdatedAt,
			&channelRule.DeletedAt)

		if err != nil {
			return entity.ChannelRuleEntity{}, err
		}
	}

	return modelToEntity(channelRule), nil
}

func (r repository) ReadChannelRules(ctx context.Context, request entity.CollectionRequest) (resp []entity.ChannelRuleEntity, err error) {
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
		return
	}

	defer rows.Close()
	for rows.Next() {
		var model entity.ChannelRuleModel
		err = rows.StructScan(&model)
		if err != nil {
			r.logger.Info(err)
			return
		}

		resp = append(resp, modelToEntity(model))
	}

	return
}

func (r repository) UpdateChannelRule(ctx context.Context, channelRule entity.ChannelRuleUpdateRequest) error {
	stmt, err := r.db.Prepare("update channel_rule set channel_id=?, tag_id=?, price=?, is_active=?, updated_at=? where id=?")
	if err != nil {
		return err
	}

	defer stmt.Close()
	// execute
	res, err := stmt.Exec(channelRule.ChannelId, channelRule.TagId,
		channelRule.Price, channelRule.IsActive, time.Now(), channelRule.ID)
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

func (r repository) DeleteChannelRule(ctx context.Context, request entity.DeleteCollectionRequest) error {
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
	result, err := r.dbx.Exec(query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	r.logger.Infof("%d rows Affected:", rowsAffected)
	return nil
}

func modelToEntity(model entity.ChannelRuleModel) entity.ChannelRuleEntity {
	return entity.ChannelRuleEntity{
		ID:          model.ID.Int64,
		ChannelId:   model.ChannelId.Int64,
		TagId:       model.TagId.Int64,
		Price:       model.Price.Float64,
		IsActive:    model.IsActive.Bool,
		Description: model.Description.String,
		CreatedAt:   model.CreatedAt.Time,
		UpdatedAt:   model.UpdatedAt.Time,
		DeletedAt:   model.DeletedAt.Time,
	}
}

func modelToResponse(model entity.ChannelRuleModel) entity.ChannelRuleResponse {
	return entity.ChannelRuleResponse{
		ID:          model.ID.Int64,
		ChannelId:   model.ChannelId.Int64,
		TagId:       model.TagId.Int64,
		Price:       model.Price.Float64,
		IsActive:    model.IsActive.Bool,
		Description: model.Description.String,
		CreatedAt:   model.CreatedAt.Time,
		UpdatedAt:   model.UpdatedAt.Time,
	}
}
