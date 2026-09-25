package utils

import (
	"cmp"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/mattn/go-sqlite3"
	"github.com/zouipo/yumsday/backend/internal/migration"
)

// TimesApproximatelyEqual checks if two time values are approximately equal within a specified tolerance.
func TimesApproximatelyEqual(t1, t2 time.Time, tolerance time.Duration) bool {
	diff := t1.Sub(t2)
	return math.Abs(diff.Seconds()) <= tolerance.Seconds()
}

// CompareErrors compares two errors to check if they are equivalent AppErrors.
// It compares the error message and underlying error if it is an sqlite error.
// Returns true if the errors are equivalent.
func CompareErrors(actual, expected error) bool {
	if actual == nil && expected == nil {
		return true
	}

	if (actual == nil) != (expected == nil) {
		return false
	}

	if actual.Error() != expected.Error() {
		return false
	}

	// Compare sqlite extended codes when both wrapped errors are sqlite3.Error.
	actualSQLErr, actualIsSQLErr := errors.AsType[sqlite3.Error](actual)
	expectedSQLErr, expectedIsSQLErr := errors.AsType[sqlite3.Error](expected)

	if actualIsSQLErr && expectedIsSQLErr {
		return actualSQLErr.ExtendedCode == expectedSQLErr.ExtendedCode
	}

	// If actual is sqlite3.Error but expected is an error code constant (ErrNoExtended or ErrNo),
	// compare the actual error's ExtendedCode with the expected constant
	if actualIsSQLErr {
		if errNoExt, ok := errors.AsType[sqlite3.ErrNoExtended](expected); ok {
			return actualSQLErr.ExtendedCode == errNoExt
		}
		if errNo, ok := errors.AsType[sqlite3.ErrNo](expected); ok {
			return actualSQLErr.ExtendedCode == sqlite3.ErrNoExtended(errNo)
		}
	}

	return true
}

func SortSliceByFieldName[T any](s []T, sortBy string, descending bool) []T {
	sorted := append([]T{}, s...)
	sortWords := strings.Split(sortBy, ".")

	sort.Slice(sorted, func(i, j int) bool {
		return compareFieldsByName(sorted[i], sorted[j], sortWords, descending)
	})

	return sorted
}

// compareFieldsByName compare t1 and t2 based on the values of the fields
// designated by sortWords' content.
// sortWords can contain multiple words to compare fields deep in the struct nesting.
// compareFieldsByName tries to reproduce the database ordering.
// See https://sqlite.org/datatype3.html#sortorder
func compareFieldsByName[T any](t1 T, t2 T, sortWords []string, descending bool) bool {
	a := reflect.ValueOf(t1).FieldByName(sortWords[0])
	b := reflect.ValueOf(t2).FieldByName(sortWords[0])
	var res bool

	// Dereference pointers recursively (handles **int, ***int, etc.)
	for a.Kind() == reflect.Pointer {
		if a.IsNil() {
			// nil is always "smaller" than anything else, including another nil
			// So if a is nil then it's smaller, even if b is also nil
			// a is the first argument so ascending => true; descending => false
			return !descending
		}
		a = a.Elem()
	}
	for b.Kind() == reflect.Pointer {
		if b.IsNil() {
			// if b is nil it is considered smaller
			// b is the second argument so we apply the opposite logic from a
			// i.e. ascending => false; descending => true
			return descending
		}
		b = b.Elem()
	}

	if a.Kind() != b.Kind() {
		panic(errors.New("cannot compare values of different kind, " +
			"check that sortWords leads to a field of the same kind in both values"))
	}

	switch a.Kind() {
	case reflect.String:
		res = cmp.Less(a.String(), b.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		res = cmp.Less(a.Int(), b.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		res = cmp.Less(a.Uint(), b.Uint())
	case reflect.Float32, reflect.Float64:
		res = cmp.Less(a.Float(), b.Float())
	case reflect.Struct:
		if len(sortWords) == 1 {
			panic(fmt.Errorf("cannot compare structs '%s'. sortWords must lead to a comparable value", sortWords[0]))
		}
		return compareFieldsByName(a.Interface(), b.Interface(), sortWords[1:], descending)
	default:
		panic(fmt.Errorf("unhandled kind %v", a.Kind()))
	}

	return res != descending
}

func SetUpTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", "file::memory:?_foreign_keys=on")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	// Apply migrations using the migration package
	migrationsFS := os.DirFS("../../data/migrations")
	err = migration.Migrate(db, migrationsFS)
	if err != nil {
		t.Fatalf("failed to apply migrations: %v", err)
	}

	testScript, _ := os.ReadFile("../../data/test.sql")
	_, err = db.Exec(string(testScript))
	if err != nil {
		t.Fatalf("failed to run test.sql: %v", err)
	}

	return db
}
