package utils

import (
	"slices"
	"strings"
)

type WhereClause struct {
	Column string
	Values []any
	Like   bool
}

type OrderByClause struct {
	Column     string
	Descending bool
}

type SelectFilteringOptions struct {
	Where   []WhereClause
	OrderBy []OrderByClause
}

func (s *SelectFilteringOptions) WhereColumns() []string {
	ret := []string{}
	for _, w := range s.Where {
		ret = append(ret, w.Column)
	}
	return ret
}

func (s *SelectFilteringOptions) WhereValues() []any {
	ret := []any{}
	for _, w := range s.Where {
		ret = append(ret, w.Values...)
	}
	return ret
}

func MakeSelectFiltering(opt *SelectFilteringOptions) string {
	filter := ""

	if len(opt.Where) > 0 {
		clauses := make([]string, 0, len(opt.Where))
		for _, w := range opt.Where {
			if w.Column == "" {
				panic("MakeSelectFiltering: column string cannot be empty")
			}
			if len(w.Values) == 0 {
				panic("MakeSelectFiltering: values list cannot be empty")
			}

			pattern := "?"
			if w.Like {
				pattern = "%?%"
			}

			if len(w.Values) > 1 {
				clause := "(" + strings.Join(
					slices.Repeat(
						[]string{"WHERE " + w.Column + " = " + pattern},
						len(w.Values),
					),
					" OR ",
				) + ")"
				clauses = append(clauses, clause)
			} else {
				clause := "WHERE " + w.Column + " = " + pattern
				clauses = append(clauses, clause)
			}
		}
		filter = strings.Join(clauses, " AND ")
	}

	if len(opt.OrderBy) > 0 {
		clauses := make([]string, 0, len(opt.OrderBy))
		// build strings like 'column1', 'column2 DESC', 'column3'
		for _, o := range opt.OrderBy {
			if o.Column == "" {
				panic("MakeSelectFiltering: column string cannot be empty")
			}
			clause := o.Column
			if o.Descending {
				clause += " DESC"
			}
			clauses = append(clauses, clause)
		}
		// build a string like 'ORDER BY column1, column2 DESC, column3'
		filter += " ORDER BY " + strings.Join(clauses, ", ")
	}

	return strings.TrimSpace(filter)
}
