package utils

import (
	"errors"
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

	// Support wrapped errors, not only direct type assertions.
	actualAppErr, actualIsAppErr := errors.AsType[*customErrors.AppError](actual)
	expectedAppErr, expectedIsAppErr := errors.AsType[*customErrors.AppError](expected)

	if !actualIsAppErr || !expectedIsAppErr {
		// If not AppErrors, compare their error messages
		return actual.Error() == expected.Error()
	}

	if actualAppErr.Message != expectedAppErr.Message || actualAppErr.StatusCode != expectedAppErr.StatusCode {
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

	// non-SQLite error
	if expectedAppErr.Err != nil {
		return actualAppErr.Err == expectedAppErr.Err
	}

	return true
}
