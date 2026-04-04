package utils

import (
	"slices"
	"testing"
)

func TestWhereColumns(t *testing.T) {
	opt := &SelectFilteringOptions{
		Where: []WhereClause{
			{Column: "test", Values: []any{1}},
			{Column: "test2", Values: []any{1}},
		},
	}
	expected := []string{"test", "test2"}

	actual := opt.WhereColumns()

	if !slices.Equal(actual, expected) {
		t.Fatalf("expected where columns '%v', got '%v'", expected, actual)
	}
}

func TestWhereValues(t *testing.T) {
	opt := &SelectFilteringOptions{
		Where: []WhereClause{
			{Column: "id", Values: []any{1, 2, 3}},
			{Column: "test", Values: []any{"aoh", 4}},
		},
	}
	expected := []any{1, 2, 3, "aoh", 4}

	actual := opt.WhereValues()

	if !slices.Equal(actual, expected) {
		t.Fatalf("expected where values '%v', got '%v'", expected, actual)
	}
}

func TestMakeSelectFiltering_EmptyFilter(t *testing.T) {
	actual := MakeSelectFiltering(&SelectFilteringOptions{})

	if actual != "" {
		t.Fatalf("filter should be empty, got %s", actual)
	}
}

func TestMakeSelectFiltering_OneWhere(t *testing.T) {
	expected := "WHERE id = ?"
	opt := &SelectFilteringOptions{
		Where: []WhereClause{
			{Column: "id", Values: []any{1}},
		},
	}

	actual := MakeSelectFiltering(opt)

	if actual != expected {
		t.Fatalf("expected '%s', got '%s'", expected, actual)
	}
}

func TestMakeSelectFiltering_MultipleWhere(t *testing.T) {
	expected := "WHERE id IN (?, ?, ?) AND test = ?"
	opt := &SelectFilteringOptions{
		Where: []WhereClause{
			{Column: "id", Values: []any{1, 2, 3}},
			{Column: "test", Values: []any{"aoh"}},
		},
	}

	actual := MakeSelectFiltering(opt)

	if actual != expected {
		t.Fatalf("expected '%s', got '%s'", expected, actual)
	}
}

func TestMakeSelectFiltering_OrderBy(t *testing.T) {
	expected := "ORDER BY name"
	opt := &SelectFilteringOptions{
		OrderBy: []OrderByClause{
			{Column: "name"},
		},
	}

	actual := MakeSelectFiltering(opt)

	if actual != expected {
		t.Fatalf("expected '%s', got '%s'", expected, actual)
	}
}

func TestMakeSelectFiltering_OrderByDesc(t *testing.T) {
	expected := "ORDER BY name DESC"
	opt := &SelectFilteringOptions{
		OrderBy: []OrderByClause{
			{Column: "name", Descending: true},
		},
	}

	actual := MakeSelectFiltering(opt)

	if actual != expected {
		t.Fatalf("expected '%s', got '%s'", expected, actual)
	}
}

func TestMakeSelectFiltering_MultipleOrderBy(t *testing.T) {
	expected := "ORDER BY name DESC, test"
	opt := &SelectFilteringOptions{
		OrderBy: []OrderByClause{
			{Column: "name", Descending: true},
			{Column: "test"},
		},
	}

	actual := MakeSelectFiltering(opt)

	if actual != expected {
		t.Fatalf("expected '%s', got '%s'", expected, actual)
	}
}

func TestMakeSelectFiltering_AllTogether(t *testing.T) {
	expected := "WHERE id IN (?, ?, ?) AND test = ? AND test2 IN (?, ?) ORDER BY value DESC, value2"
	opt := &SelectFilteringOptions{
		Where: []WhereClause{
			{Column: "id", Values: []any{1, "2", 3.0}},
			{Column: "test", Values: []any{4}},
			{Column: "test2", Values: []any{"yes", "no"}},
		},
		OrderBy: []OrderByClause{
			{Column: "value", Descending: true},
			{Column: "value2"},
		},
	}

	actual := MakeSelectFiltering(opt)

	if actual != expected {
		t.Fatalf("expected '%s', got '%s'", expected, actual)
	}
}
