package utils

import (
	"testing"
)

func TestNewSelectFilteringOptions_DifferentListsLength(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	NewSelectFilteringOptions([]string{"test", "test"}, []any{1}, "", false)
}

func TestMakeSelectFiltering_EmptyFilter(t *testing.T) {
	tests := []*selectFilteringOptions{
		NewSelectFilteringOptions([]string{}, []any{}, "", false),
		NewSelectFilteringOptions([]string{}, []any{}, "", true),
	}

	for _, tt := range tests {
		actual := MakeSelectFiltering(tt)
		if actual != "" {
			t.Fatalf("filter should be empty, got %s", actual)
		}
	}
}

func TestMakeSelectFiltering_OneWhere(t *testing.T) {
	expected := "WHERE id = ?"
	opt := NewSelectFilteringOptions([]string{"id"}, []any{1}, "", false)
	actual := MakeSelectFiltering(opt)

	if actual != expected {
		t.Fatalf("expected '%s', got '%s'", expected, actual)
	}
}

func TestMakeSelectFiltering_MultipleWhere(t *testing.T) {
	expected := "WHERE id = ? AND name = ?"
	opt := NewSelectFilteringOptions(
		[]string{"id", "name"},
		[]any{1, "test"},
		"",
		false,
	)
	actual := MakeSelectFiltering(opt)

	if actual != expected {
		t.Fatalf("expected '%s', got '%s'", expected, actual)
	}
}

func TestMakeSelectFiltering_OrderBy(t *testing.T) {
	expected := "ORDER BY name"
	opt := NewSelectFilteringOptions([]string{}, []any{}, "name", false)
	actual := MakeSelectFiltering(opt)

	if actual != expected {
		t.Fatalf("expected '%s', got '%s'", expected, actual)
	}
}

func TestMakeSelectFiltering_OrderByDesc(t *testing.T) {
	expected := "ORDER BY name DESC"
	opt := NewSelectFilteringOptions([]string{}, []any{}, "name", true)
	actual := MakeSelectFiltering(opt)

	if actual != expected {
		t.Fatalf("expected '%s', got '%s'", expected, actual)
	}
}

func TestMakeSelectFiltering_AllTogether(t *testing.T) {
	expected := "WHERE id = ? AND name = ? ORDER BY value DESC"
	opt := NewSelectFilteringOptions([]string{"id", "name"}, []any{1, "test"}, "value", true)
	actual := MakeSelectFiltering(opt)

	if actual != expected {
		t.Fatalf("expected '%s', got '%s'", expected, actual)
	}
}
