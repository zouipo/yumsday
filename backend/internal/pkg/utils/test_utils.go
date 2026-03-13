package utils

import (
	"math"
	"time"

	"github.com/mattn/go-sqlite3"
	customErrors "github.com/zouipo/yumsday/backend/internal/error"
)

// TimesApproximatelyEqual checks if two time values are approximately equal within a specified tolerance.
func TimesApproximatelyEqual(t1, t2 time.Time, tolerance time.Duration) bool {
	diff := t1.Sub(t2)
	return math.Abs(diff.Seconds()) <= tolerance.Seconds()
}

// CompareErrors compares two errors to check if they are equivalent AppErrors.
// It compares the Message, StatusCode, and underlying Err fields.
func CompareErrors(actual, expected error) bool {
	if actual == nil && expected == nil {
		return true
	}

	if (actual == nil) != (expected == nil) {
		return false
	}

	// Cast both to *AppError
	actualAppErr, actualIsAppErr := actual.(*customErrors.AppError)
	expectedAppErr, expectedIsAppErr := expected.(*customErrors.AppError)

	if actualIsAppErr && expectedIsAppErr {
		if actualAppErr.Message != expectedAppErr.Message || actualAppErr.StatusCode != expectedAppErr.StatusCode {
			return false
		}

		// Cast both into sqlite3.Error to compare their ExtendedCode if possible
		actualSQLErr, actualIsSQLErr := actualAppErr.Err.(sqlite3.Error)
		expectedSQLErr, expectedIsSQLErr := expectedAppErr.Err.(sqlite3.Error)

		if actualIsSQLErr && expectedIsSQLErr {
			return actualSQLErr.ExtendedCode == expectedSQLErr.ExtendedCode
		}

		// If actual is sqlite3.Error but expected is an error code constant (ErrNoExtended or ErrNo),
		// compare the actual error's ExtendedCode with the expected constant
		if actualIsSQLErr {
			if errNoExt, ok := expectedAppErr.Err.(sqlite3.ErrNoExtended); ok {
				return actualSQLErr.ExtendedCode == errNoExt
			}
			if errNo, ok := expectedAppErr.Err.(sqlite3.ErrNo); ok {
				return actualSQLErr.ExtendedCode == sqlite3.ErrNoExtended(errNo)
			}
		}

		// non-SQLite error
		if expectedAppErr.Err != nil {
			return actualAppErr.Err == expectedAppErr.Err
		}
		return true
	}

	// If not AppErrors, compare their error messages
	return actual.Error() == expected.Error()
}
