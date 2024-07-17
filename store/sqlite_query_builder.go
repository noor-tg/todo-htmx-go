package store

import (
	"fmt"
	"strings"
)

func FilterBy(key string, operator string, filter string, query *string, queryArgs *[]interface{}) {
	if filter != "" {
		if !ContainWhere(*query) {
			*query += " WHERE"
		} else {
			*query += " AND"
		}
		if operator == "LIKE" {
			*queryArgs = append(*queryArgs, fmt.Sprintf("%%%s%%", filter))
		} else {
			*queryArgs = append(*queryArgs, filter)
		}
		*query += fmt.Sprintf(" %s %s ?", key, operator)
	}
}

func ContainWhere(query string) bool {
	if strings.Contains(query, "WHERE") {
		return true
	}
	return false
}
