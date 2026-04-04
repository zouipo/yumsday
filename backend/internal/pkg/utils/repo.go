package utils

import (
	"fmt"
	"strings"
)

type SelectFilteringOptions struct {
	whereColumns      []string
	whereValues       []any
	orderByColumn     string
	orderByDescending bool
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
		whereColumns:      whereColumns,
		whereValues:       whereValues,
		orderByColumn:     orderByColumn,
		orderByDescending: orderByDescending,
	}
}

func MakeSelectFiltering(opt *SelectFilteringOptions) string {
	filter := ""

	if len(opt.whereColumns) > 0 {
		filter = "WHERE "
		for i := range opt.whereColumns {
			filter += fmt.Sprintf("%v = ? ", opt.whereColumns[i])
			if i < len(opt.whereColumns)-1 {
				filter += "AND "
			}
		}
	}

	if opt.orderByColumn != "" {
		filter += fmt.Sprintf("ORDER BY %v ", opt.orderByColumn)
		if opt.orderByDescending {
			filter += "DESC "
		}
	}

	return strings.TrimSpace(filter)
}
