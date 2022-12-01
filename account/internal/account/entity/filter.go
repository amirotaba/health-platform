package entity

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
)

type Collection interface {
	IBuildQuery() (string, []interface{})
}

type FilterType int64

const (
	GT FilterType = iota
	GE
	LT
	LE
	EQ
	IN
	NE
)

type SortType int64

const (
	Ascending  SortType = 1
	Descending          = -1
)

func (e *SortType) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	value, ok := map[string]SortType{"asc": Ascending, "desc": Descending}[s]
	if !ok {
		return errors.New("invalid SortType value")
	}

	*e = value
	return nil
}

func (e *SortType) MarshalJSON() ([]byte, error) {
	value, ok := map[SortType]string{Ascending: "asc", Descending: "desc"}[*e]
	if !ok {
		return nil, errors.New("invalid SortType value")
	}
	return json.Marshal(value)
}

type CollectionRequest struct {
	Page            int            `json:"page"`
	PerPage         int            `json:"per_page"`
	ShowDeleted     bool           `json:"show_deleted"`
	ShowActivated   bool           `json:"show_activated"`
	Model           interface{}    `json:"model"`
	TableName       string         `json:"table_name"`
	OrderBy         OrderBy        `json:"order_by"`
	Sort            []SortedField  `json:"sort"`
	Searches        []Search       `json:"searches"`
	Filters         []Filter       `json:"filters"`
	PaginateData    Paginate       `json:"paginate"`
	JoinWithOptions JoinWithOption `json:"joinWithOptions"`
	Tables          []Table        `json:"tables"`
	Result          *sqlx.Rows     `json:"result"`
	Sql             string         `json:"sql"`
	Query           string         `json:"query"`
	Args            []interface{}  `json:"args"`
	DB              *sqlx.DB       `json:"db"`
}

func NewCollectionRequest(tableName string) CollectionRequest {
	result := CollectionRequest{TableName: tableName}
	return result
}

func (collection CollectionRequest) EQ(field string, value interface{}) CollectionRequest {
	collection.Filters = append(collection.Filters, Filter{
		TableName: collection.TableName,
		Field:     field,
		Type:      EQ,
		Value:     value,
	})

	return collection
}

func (collection CollectionRequest) LT(field string, value interface{}) CollectionRequest {
	collection.Filters = append(collection.Filters, Filter{
		TableName: collection.TableName,
		Field:     field,
		Type:      LT,
		Value:     value,
	})

	return collection
}

func (collection CollectionRequest) LE(field string, value interface{}) CollectionRequest {
	collection.Filters = append(collection.Filters, Filter{
		TableName: collection.TableName,
		Field:     field,
		Type:      LE,
		Value:     value,
	})

	return collection
}

func (collection CollectionRequest) GT(field string, value interface{}) CollectionRequest {
	collection.Filters = append(collection.Filters, Filter{
		TableName: collection.TableName,
		Field:     field,
		Type:      GE,
		Value:     value,
	})

	return collection
}

func (collection CollectionRequest) GE(field string, value interface{}) CollectionRequest {
	collection.Filters = append(collection.Filters, Filter{
		TableName: collection.TableName,
		Field:     field,
		Type:      GE,
		Value:     value,
	})

	return collection
}

func (collection CollectionRequest) NE(field string, value interface{}) CollectionRequest {
	collection.Filters = append(collection.Filters, Filter{
		TableName: collection.TableName,
		Field:     field,
		Type:      NE,
		Value:     value,
	})

	return collection
}

func (collection CollectionRequest) IN(field string, values ...interface{}) CollectionRequest {
	collection.Filters = append(collection.Filters, Filter{
		TableName: collection.TableName,
		Field:     field,
		Type:      IN,
		Values:    values,
	})

	return collection
}

func (collection CollectionRequest) SQL(raw string) CollectionRequest {
	collection.Sql = raw
	return collection
}

func (collection CollectionRequest) BuildQuery() CollectionRequest {
	selectBuilder := sqlbuilder.NewStruct(collection.Model)
	sb := selectBuilder.SelectFrom(collection.TableName)
	for _, filter := range collection.Filters {
		switch filter.Type {
		//case domain.KeyFromDate:
		//	sb.Where(sb.GE(domain.KeyFromDate, filter.Value[0]))
		//case domain.KeyToDate:
		//	sb.Where(sb.LE(domain.KeyToDate, filter.Value[0]))
		case EQ:
			sb.Where(sb.E(filter.Field, filter.Value))
		case GE:
			sb.Where(sb.GE(filter.Field, filter.Value))
		case LE:
			sb.Where(sb.LE(filter.Field, filter.Value))
		case GT:
			sb.Where(sb.G(filter.Field, filter.Value))
		case LT:
			sb.Where(sb.L(filter.Field, filter.Value))
		case NE:
			sb.Where(sb.NE(filter.Field, filter.Value))
		default:
			sb.Where(sb.In(filter.Field, filter.Values...))
		}
	}

	var searches []string
	for _, search := range collection.Searches {
		searches = append(searches, sb.Like(search.Field, search.Value))
	}

	if len(searches) > 0 {
		sb.Where(searches...)
	}

	if collection.PaginateData.Has {
		sb.Offset(collection.PaginateData.Offset).Limit(collection.PaginateData.Limit)
	}

	if collection.OrderBy.Has {
		sb.OrderBy(collection.OrderBy.Field)
		if collection.OrderBy.Ascending {
			sb.Asc()
		} else {
			sb.Desc()
		}
	}

	//query, args := sb.Build()
	//r.logger.Info(query, args)
	collection.Query, collection.Args = sb.Build()
	return collection
}

func (collection CollectionRequest) IBuildQuery() (string, []interface{}) {
	selectBuilder := sqlbuilder.NewStruct(collection.Model)
	sb := selectBuilder.SelectFrom(collection.TableName)
	for _, filter := range collection.Filters {
		switch filter.Type {
		case EQ:
			sb.Where(sb.E(filter.Field, filter.Value))
		case GE:
			sb.Where(sb.GE(filter.Field, filter.Value))
		case LE:
			sb.Where(sb.LE(filter.Field, filter.Value))
		case GT:
			sb.Where(sb.G(filter.Field, filter.Value))
		case LT:
			sb.Where(sb.L(filter.Field, filter.Value))
		case NE:
			sb.Where(sb.NE(filter.Field, filter.Value))
		default:
			sb.Where(sb.In(filter.Field, filter.Values...))
		}
	}

	var searches []string
	for _, search := range collection.Searches {
		searches = append(searches, sb.Like(search.Field, search.Value))
	}

	if len(searches) > 0 {
		sb.Where(searches...)
	}

	if collection.PaginateData.Has {
		sb.Offset(collection.PaginateData.Offset).Limit(collection.PaginateData.Limit)
	}

	if collection.OrderBy.Has {
		sb.OrderBy(collection.OrderBy.Field)
		if collection.OrderBy.Ascending {
			sb.Asc()
		} else {
			sb.Desc()
		}
	}

	//query, args := sb.Build()
	//r.logger.Info(query, args)
	//collection.Query, collection.Args = sb.Build()
	//return collection
	return sb.Build()
}

func (collection CollectionRequest) SetModel(model interface{}) CollectionRequest {
	collection.Model = model
	return collection
}

func (collection CollectionRequest) Scan(in interface{}) {
	defer collection.Result.Close()
	for collection.Result.Next() {
		err := collection.Result.Scan(in)
		if err != nil {
			return
		}

		return
	}
}

func (collection CollectionRequest) Get(dst interface{}) error {
	return collection.DB.Get(dst, collection.Query, collection.SQL)
}

func (collection CollectionRequest) StructScan(dest interface{}) error {
	return collection.Result.StructScan(dest)
}

func (collection CollectionRequest) MapScan(dest map[string]interface{}) error {
	return collection.Result.MapScan(dest)
}

func (collection CollectionRequest) SliceScan() ([]interface{}, error) {
	return collection.Result.SliceScan()
}

type DeleteCollectionRequest struct {
	Filters     []Filter `json:"filters"`
	OR          bool     `json:"or"`
	ShowDeleted bool     `json:"show_deleted"`
	TableName   string   `json:"table_name"`
}

func NewDeleteCollectionRequest(tableName string) DeleteCollectionRequest {
	return DeleteCollectionRequest{
		TableName: tableName,
	}
}

func (collection DeleteCollectionRequest) EQ(field string, value interface{}) DeleteCollectionRequest {
	collection.Filters = append(collection.Filters, Filter{
		TableName: collection.TableName,
		Field:     field,
		Type:      EQ,
		Value:     value,
	})

	return collection
}

func (collection DeleteCollectionRequest) LT(field string, value interface{}) DeleteCollectionRequest {
	collection.Filters = append(collection.Filters, Filter{
		TableName: collection.TableName,
		Field:     field,
		Type:      LT,
		Value:     value,
	})

	return collection
}

func (collection DeleteCollectionRequest) LE(field string, value interface{}) DeleteCollectionRequest {
	collection.Filters = append(collection.Filters, Filter{
		TableName: collection.TableName,
		Field:     field,
		Type:      LE,
		Value:     value,
	})

	return collection
}

func (collection DeleteCollectionRequest) GT(field string, value interface{}) DeleteCollectionRequest {
	collection.Filters = append(collection.Filters, Filter{
		TableName: collection.TableName,
		Field:     field,
		Type:      GE,
		Value:     value,
	})

	return collection
}

func (collection DeleteCollectionRequest) GE(field string, value interface{}) DeleteCollectionRequest {
	collection.Filters = append(collection.Filters, Filter{
		TableName: collection.TableName,
		Field:     field,
		Type:      GE,
		Value:     value,
	})

	return collection
}

func (collection DeleteCollectionRequest) NE(field string, value interface{}) DeleteCollectionRequest {
	collection.Filters = append(collection.Filters, Filter{
		TableName: collection.TableName,
		Field:     field,
		Type:      NE,
		Value:     value,
	})

	return collection
}

func (collection DeleteCollectionRequest) IN(field string, values ...interface{}) DeleteCollectionRequest {
	collection.Filters = append(collection.Filters, Filter{
		TableName: collection.TableName,
		Field:     field,
		Type:      IN,
		Values:    values,
	})

	return collection
}

type UpdateCollectionRequest struct {
	Filters   []Filter    `json:"filters"`
	TableName string      `json:"table_name"`
	Model     interface{} `json:"model"`
}

func NewUpdateCollectionRequest(tableName string) UpdateCollectionRequest {
	return UpdateCollectionRequest{TableName: tableName}
}

func (collection UpdateCollectionRequest) EQ(field string, value interface{}) UpdateCollectionRequest {
	collection.Filters = append(collection.Filters, Filter{
		TableName: collection.TableName,
		Field:     field,
		Type:      EQ,
		Value:     value,
	})

	return collection
}

func (collection UpdateCollectionRequest) LT(field string, value interface{}) UpdateCollectionRequest {
	collection.Filters = append(collection.Filters, Filter{
		TableName: collection.TableName,
		Field:     field,
		Type:      LT,
		Value:     value,
	})

	return collection
}

func (collection UpdateCollectionRequest) LE(field string, value interface{}) UpdateCollectionRequest {
	collection.Filters = append(collection.Filters, Filter{
		TableName: collection.TableName,
		Field:     field,
		Type:      LE,
		Value:     value,
	})

	return collection
}

func (collection UpdateCollectionRequest) GT(field string, value interface{}) UpdateCollectionRequest {
	collection.Filters = append(collection.Filters, Filter{
		TableName: collection.TableName,
		Field:     field,
		Type:      GE,
		Value:     value,
	})

	return collection
}

func (collection UpdateCollectionRequest) GE(field string, value interface{}) UpdateCollectionRequest {
	collection.Filters = append(collection.Filters, Filter{
		TableName: collection.TableName,
		Field:     field,
		Type:      GE,
		Value:     value,
	})

	return collection
}

func (collection UpdateCollectionRequest) NE(field string, value interface{}) UpdateCollectionRequest {
	collection.Filters = append(collection.Filters, Filter{
		TableName: collection.TableName,
		Field:     field,
		Type:      NE,
		Value:     value,
	})

	return collection
}

func (collection UpdateCollectionRequest) IN(field string, values ...interface{}) UpdateCollectionRequest {
	collection.Filters = append(collection.Filters, Filter{
		TableName: collection.TableName,
		Field:     field,
		Type:      IN,
		Values:    values,
	})

	return collection
}

func (collection UpdateCollectionRequest) SetModel(model interface{}) UpdateCollectionRequest {
	collection.Model = model
	return collection
}

func (collection UpdateCollectionRequest) IBuildQuery() (string, []interface{}) {
	selectBuilder := sqlbuilder.NewStruct(collection.Model)
	sb := selectBuilder.Update(collection.TableName, collection.Model)
	for _, filter := range collection.Filters {
		switch filter.Type {
		case EQ:
			sb.Where(sb.E(filter.Field, filter.Value))
		case GE:
			sb.Where(sb.GE(filter.Field, filter.Value))
		case LE:
			sb.Where(sb.LE(filter.Field, filter.Value))
		case GT:
			sb.Where(sb.G(filter.Field, filter.Value))
		case LT:
			sb.Where(sb.L(filter.Field, filter.Value))
		case NE:
			sb.Where(sb.NE(filter.Field, filter.Value))
		default:
			sb.Where(sb.In(filter.Field, filter.Values...))
		}
	}

	return sb.Build()
}

type InsertCollectionRequest struct {
	TableName string      `json:"table_name"`
	Model     interface{} `json:"model"`
}

func NewInsertCollectionRequest(tableName string, model interface{}) InsertCollectionRequest {
	return InsertCollectionRequest{
		TableName: tableName,
		Model:     model,
	}
}

type Filter struct {
	TableName string        `json:"table_name"`
	Field     string        `json:"field"`
	Type      FilterType    `json:"type"` // "=", ">", "<", ">=", "<=", "<>", "LIKE", "NOT LIKE", "IN", "NOT IN", "IS NULL"
	Value     interface{}   `json:"value"`
	Values    []interface{} `json:"values"`
	OR        bool          `json:"or"`
}

type Search struct {
	TableName string      `json:"table_name"`
	Field     string      `json:"field"`
	Value     interface{} `json:"value"`
}

type QueryFilter struct {
	TableName string     `json:"table_name"`
	Field     string     `json:"field"`
	Type      FilterType `json:"type"` // "=", ">", "<", ">=", "<=", "<>", "LIKE", "NOT LIKE", "IN", "NOT IN", "IS NULL"

}

type QuerySearch struct {
	TableName string `json:"table_name"`
	Field     string `json:"field"`
}

func NewFilter(tableName, field string, values ...interface{}) Filter {
	return Filter{
		TableName: tableName,
		Field:     field,
		Value:     values,
	}
}

func NewSearch(tableName, field string, values ...interface{}) Search {
	return Search{
		TableName: tableName,
		Field:     field,
		Value:     values,
	}
}

type SortedField struct {
	Field string   `json:"field"`
	Type  SortType `json:"type"` // SortType: "asc" 1 - Ascending, "desc" -1 - Descending
}

func (collection *CollectionRequest) AddSortedField(sortField string, sortType SortType) {
	sortedField := SortedField{
		Field: sortField,
		Type:  sortType,
	}

	collection.Sort = append(collection.Sort, sortedField)
}

type Paginate struct {
	Has    bool `json:"has"`
	Limit  int  `json:"limit"`
	Offset int  `json:"offset"`
}

type OrderBy struct {
	Has       bool
	Field     string
	Ascending bool
}

type FilterSearchRequest struct {
	Filters  []Filter `json:"filters"`
	Searches []Search `json:"searches"`
	Page     int      `json:"page"`
	PerPage  int      `json:"per_page"`
	Paginate Paginate `json:"paginate"`
	Sort     OrderBy  `json:"sort"`
}

type JoinWithOption struct {
	RootTable string `json:"root_table"`
	Joins     []Join `json:"joins"`
}

type Join struct {
	JoinType   string `json:"join_type"`
	LeftTable  string `json:"left_table"`
	RightTable string `json:"right_table"`
	ON         ON     `json:"on"`
}

type ON struct {
	LeftON  string `json:"parent_on"`
	RightON string `json:"child_on"`
}

type Table struct {
	Model interface{}
	Name  string
}

type Transaction struct {
	DB          *sqlx.DB
	Result      sql.Result
	Collections []Collection
}

func (t Transaction) Execute(ctx context.Context) error {
	tx, err := t.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	for _, collection := range t.Collections {
		query, args := collection.IBuildQuery()
		t.Result, err = tx.ExecContext(ctx, query, args)
		if err != nil {
			rollBackErr := tx.Rollback()
			if rollBackErr != nil {
				return rollBackErr
			}

			return err
		}
	}

	commitErr := tx.Commit()
	if commitErr != nil {
		return commitErr
	}
	return nil
}
