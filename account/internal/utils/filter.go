package utils

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func BuildFilterString(query string, filters map[string][]string, allowedFilters map[string]string) (string, []interface{}, error) {
	filterString := "WHERE 1=1"
	var inputArgs []interface{}

	// filter key is the url name of the filter used as the lookup for the allowed filters list
	for filterKey, filterValList := range filters {
		if realFilterName, ok := allowedFilters[filterKey]; ok {
			if len(filterValList) == 0 {
				continue
			}

			filterString = fmt.Sprintf("%s AND %s IN (?)", filterString, realFilterName)
			inputArgs = append(inputArgs, filterValList)
		}
	}
	// template the where clause into the original query and then expand the IN clauses with sqlx
	query, args, err := sqlx.In(fmt.Sprintf(query, filterString), inputArgs...)
	if err != nil {
		return "", nil, err
	}
	// using postgres means we need to rebind the ? bindvars that sqlx.IN creates by default to $ bindvars
	// you can omit this if you are using mysql
	query = sqlx.Rebind(sqlx.DOLLAR, query)
	return query, args, nil
}
