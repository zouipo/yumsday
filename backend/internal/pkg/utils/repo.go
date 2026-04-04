package utils

import (
	"fmt"
	"strings"
)

type SelectFilteringOptions struct {
	WhereColumns      []string
	WhereValues       []any
	OrderByColumn     string
	OrderByDescending bool
}

func NewSelectFilteringOptions(
	whereColumns []string,
	whereValues []any,
	orderByColumn string,
	orderByDescending bool,
) *SelectFilteringOptions {
	if len(whereColumns) != len(whereValues) {
		panic("SelectFilteringOptions: columns and values list must have the same length")
	}
	return &SelectFilteringOptions{
		WhereColumns:      whereColumns,
		WhereValues:       whereValues,
		OrderByColumn:     orderByColumn,
		OrderByDescending: orderByDescending,
	}
}

func MakeSelectFiltering(opt *SelectFilteringOptions) string {
	filter := ""

	if len(opt.WhereColumns) > 0 {
		filter = "WHERE "
		for i := range opt.WhereColumns {
			filter += fmt.Sprintf("%v = ? ", opt.WhereColumns[i])
			if i < len(opt.WhereColumns)-1 {
				filter += "AND "
			}
		}
	}

	if opt.OrderByColumn != "" {
		filter += fmt.Sprintf("ORDER BY %v ", opt.OrderByColumn)
		if opt.OrderByDescending {
			filter += "DESC "
		}
	}

	return strings.TrimSpace(filter)
}
