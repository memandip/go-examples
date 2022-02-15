package mysql

import "fmt"

func GenerateSelectQuery(selection []string, tableName string, conditions map[string]interface{}) string {
	selectFields := ""
	whereClause := ""

	for i, s := range selection {
		var isLast bool
		if (i + 1) == len(selection) {
			isLast = true
		} else {
			isLast = false
		}
		selectFields += fmt.Sprintf("%v", s)
		if !isLast {
			selectFields += ", "
		}
	}

	if len(selectFields) == 0 {
		selectFields = "*"
	}

	first := true
	for k, v := range conditions {
		if !first {
			whereClause += " AND "
		}
		whereClause += fmt.Sprintf("%v='%v'", k, v)
		if first {
			first = false
		}
	}

	query := fmt.Sprintf("SELECT %v FROM %v", selectFields, tableName)

	if len(whereClause) > 0 {
		query += fmt.Sprintf(" WHERE %v", whereClause)
	}

	return query
}
