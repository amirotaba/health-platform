package utils

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	"git.paygear.ir/giftino/account/internal/account/domain"
	"git.paygear.ir/giftino/account/internal/account/entity"
)

//JSONError is used for returning json valid error message
type JSONError struct {
	Message string `json:"message"`
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func Base64Decode(seg string) ([]byte, error) {
	seg = Base64ScapePadding(seg)
	return base64.URLEncoding.DecodeString(seg)
}

func Base64ScapePadding(seg string) string {
	seg = strings.TrimRight(seg, "=")
	if l := len(seg) % 4; l > 0 {
		seg += strings.Repeat("=", 4-l)
	}
	return seg
}

func base64ScapeCharacter(seg string) string {
	for i := 0; i < len(seg); i++ {
		if seg[i] == '_' {
			seg = seg[:i] + "/" + seg[i+1:]
		} else if seg[i] == '-' {
			seg = seg[:i] + "+" + seg[i+1:]
		}
	}
	return seg
}

func IntInSlice(arr []int64, n int64) bool {
	for _, i := range arr {
		if n == i {
			return true
		}
	}

	return false
}

func RemoveDuplicateInt(intSlice []int64) []int64 {
	allKeys := make(map[int64]struct{})
	var list []int64
	for _, item := range intSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = struct{}{}
			list = append(list, item)
		}
	}
	return list
}

func StringInSlice(arr []string, s string) bool {
	for _, i := range arr {
		if s == i {
			return true
		}
	}

	return false
}

func RoleToString(roles []entity.RoleEntity) []string {
	var s []string
	for _, role := range roles {
		s = append(s, role.Name)
	}
	return s
}

func GetSearchAndFilter(c echo.Context, searches []string, filters []string) ([]entity.Search, []entity.Filter) {
	params := c.QueryParams()
	var queryParams entity.FilterSearchRequest
	for _, search := range searches {
		if s, ok := params[search]; ok {
			queryParams.Searches = append(queryParams.Searches, entity.Search{Field: search, Value: "%" + s[0] + "%"})
		}
	}

	for _, filter := range filters {
		if f, ok := params[filter]; ok {
			var ffs []interface{}
			for _, ff := range f {
				ffs = append(ffs, ff)
			}

			queryParams.Filters = append(queryParams.Filters, entity.Filter{Field: filter, Type: entity.IN, Values: ffs})
		}
	}

	return queryParams.Searches, queryParams.Filters
}

func GetSearchAndFilterNew(c echo.Context, searches []entity.QuerySearch, filters []entity.QueryFilter) (queryParams entity.FilterSearchRequest, err error) {
	params := c.QueryParams()
	//var err error /
	//var queryParams entity.FilterSearchRequest
	for _, search := range searches {
		if s, ok := params[search.Field]; ok {
			queryParams.Searches = append(queryParams.Searches, entity.Search{Field: fmt.Sprintf("%s.%s", search.TableName, search.Field), Value: "%" + s[0] + "%"})
		}
	}

	for _, filter := range filters {
		if f, ok := params[filter.Field]; ok {
			var ffs []interface{}
			for _, ff := range f {
				ffs = append(ffs, ff)
			}

			queryParams.Filters = append(queryParams.Filters, entity.Filter{Field: fmt.Sprintf("%s.%s", filter.TableName, filter.Field), Type: entity.EQ, Value: ffs[0]})
		}
	}

	paginate := false
	if c.QueryParam(domain.KeyPage) != domain.KeyEmptyString {
		queryParams.Page, err = strconv.Atoi(c.QueryParam(domain.KeyPage))
		if err != nil {
			return //c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
		}

		paginate = true
	} else {
		queryParams.Page = 1
	}

	if c.QueryParam(domain.KeyPerPage) != domain.KeyEmptyString {
		queryParams.PerPage, err = strconv.Atoi(c.QueryParam(domain.KeyPerPage))
		if err != nil {
			return // c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
		}

		paginate = true
	} else {
		queryParams.PerPage = 10
	}

	queryParams.Paginate = entity.Paginate{
		Has:    paginate,
		Limit:  queryParams.PerPage,
		Offset: queryParams.PerPage * (queryParams.Page - 1),
	}

	if c.QueryParam(domain.KeyASC) != domain.KeyEmptyString {
		queryParams.Sort = entity.OrderBy{
			Has:       true,
			Field:     c.QueryParam(domain.KeyASC),
			Ascending: true,
		}
	} else if c.QueryParam(domain.KeyDESC) != domain.KeyEmptyString {
		queryParams.Sort = entity.OrderBy{
			Has:       true,
			Field:     c.QueryParam(domain.KeyDESC),
			Ascending: false,
		}
	} else {
		queryParams.Sort = entity.OrderBy{
			Has:       true,
			Field:     "created_at",
			Ascending: false,
		}
	}

	return //queryParams.Searches, queryParams.Filters
}

func GetOneFilterNew(c echo.Context, filter entity.QueryFilter) (entity.Filter, error) {
	params := c.QueryParams()
	//var err error /
	//var queryParams entity.FilterSearchRequest
	if f, ok := params[filter.Field]; ok {
		var ffs []interface{}
		for _, ff := range f {
			ffs = append(ffs, ff)
		}

		return entity.Filter{Field: fmt.Sprintf("%s.%s", filter.TableName, filter.Field), Type: entity.IN, Values: ffs}, nil
	}

	return entity.Filter{}, nil //queryParams.Searches, queryParams.Filters
}

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	//log.Println(snake)
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	//log.Println(snake)
	return strings.ToLower(snake)
}

func GenerateSelectStatement(tables []entity.Table) string {
	var fields []string
	for _, table := range tables {
		t := reflect.TypeOf(table.Model)
		//log.Println(strcase.ToSnake(t.Name()))
		//table := strcase.ToSnake(t.Name())
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			var db, as string
			if value, ok := f.Tag.Lookup("db"); ok {
				if value != "-" {
					//fmt.Println("Tag found : ", value)
					db = value
				}
			}

			if value, ok := f.Tag.Lookup("as"); ok {
				as = value
			}

			fields = append(fields, fmt.Sprintf("%s.%s as %s", table, db, as))
		}
	}

	return strings.Join(fields, ", ")
}
