package dbutils

import (
	"database/sql"
	"log/slog"

	_ "modernc.org/sqlite"
)

func OpenDb(path string) (*sql.DB, error) {
	// enable foreign keys constraint
	dbOpts := "_pragma=foreign_keys(1)&" +
		// force driver to store timestamps as integer
		"_time_integer_format=unix_nano&" +
		// tell driver to convert integer to time.Time when reading from db
		"_inttotime=1"

	dsn := "file:" + path + "?" + dbOpts

	slog.Debug("opening db", "path", path, "dsn", dsn)

	return sql.Open("sqlite", dsn)
}
