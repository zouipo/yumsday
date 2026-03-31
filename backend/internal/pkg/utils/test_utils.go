package utils

import (
	"cmp"
	"errors"
	"fmt"
	"math"
	"reflect"
	"sort"
	"time"

	"github.com/mattn/go-sqlite3"
)

// Ptr returns a pointer of type T initialized to v.
func Ptr[T any](v T) *T {
	return &v
}

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
	sort.Slice(sorted, func(i, j int) bool {
		a := reflect.ValueOf(s[i]).FieldByName(sortBy)
		b := reflect.ValueOf(s[j]).FieldByName(sortBy)
		var res int

		switch a.Kind() {
		case reflect.String:
			res = cmp.Compare(a.String(), b.String())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			res = cmp.Compare(a.Int(), b.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			res = cmp.Compare(a.Uint(), b.Uint())
		case reflect.Float32, reflect.Float64:
			res = cmp.Compare(a.Float(), b.Float())
		default:
			panic(fmt.Errorf("unhandled kind %v", a.Kind()))
		}

		if descending {
			return res == 1
		} else {
			return res == -1
		}
	})

	return sorted
}
