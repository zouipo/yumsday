package migration

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log/slog"
	"regexp"
	"strconv"
)

type migration struct {
	version int
	name    string
	script  string
}

// Applies the migration scripts contained in migrationFs on db.
// migrationsFs must contain scripts at its root.
// If subdirectories are found, they are ignored and a warning is logged.
//
// Scripts names must match the following format :
//
// ^(\d+)_(.+)\.sql$ => <version_number>_<script_name>.sql
//
// The following filename is a valid migration script name :
//
// 1_db-creation.sql
//
// The version of the last applied migration is stored in
// the _migration_version table of the database.
// On next calls, only migrations with a greater version than
// the one currently in db are applied.
func Migrate(db *sql.DB, migrationsFs fs.FS) error {
	slog.Info("Starting database migration")

	migrations, err := loadMigrations(migrationsFs)
	if err != nil {
		return fmt.Errorf("Failed to load migrations: %w", err)
	}

	err = initializeMigrationVersion(db)
	if err != nil {
		return fmt.Errorf("Failed to initialize migration version: %w", err)
	}

	currentVersion, err := getMigrationVersion(db)
	if err != nil {
		return fmt.Errorf("Failed to get current migration version: %w", err)
	}
	slog.Info("Current database migration version", "version", currentVersion)

	err = performMigrations(db, migrations, currentVersion)
	if err != nil {
		return fmt.Errorf("Failed to perform migrations: %w", err)
	}

	slog.Info("Database migration completed successfully")
	return nil
}

func loadMigrations(scriptsFs fs.FS) ([]migration, error) {
	slog.Debug("Loading migrations")
	var migrations []migration
	migrationFiles, err := fs.ReadDir(scriptsFs, ".")
	if err != nil {
		return nil, err
	}

	for _, file := range migrationFiles {
		slog.Debug("Visiting file", "file", file.Name())

		if file.IsDir() {
			slog.Warn("Directory found in migrations", "directory", file.Name())
			continue
		}

		re := regexp.MustCompile(`^(\d+)_(.+)\.sql$`)
		matches := re.FindStringSubmatch(file.Name())

		if len(matches) != 3 {
			return nil, fmt.Errorf("Cannot parse migration version and name from filename %s", file.Name())
		}

		var version int
		var name string

		// No need to check for error, regex ensures it's a valid integer
		version, _ = strconv.Atoi(matches[1])

		name = matches[2]

		scriptBytes, err := fs.ReadFile(scriptsFs, file.Name())
		if err != nil {
			return nil, err
		}

		m := migration{
			version: version,
			name:    name,
			script:  string(scriptBytes),
		}
		migrations = append(migrations, m)
	}

	return migrations, nil
}

func performMigrations(db *sql.DB, migrations []migration, currentVersion int) error {
	for _, m := range migrations {
		if m.version <= currentVersion {
			slog.Debug("Skipping already applied migration", "version", m.version, "name", m.name)
			continue
		}

		slog.Info("Applying migration", "version", m.version, "name", m.name)
		_, err := db.Exec(m.script)
		if err != nil {
			return fmt.Errorf("Failed to apply migration %d_%s: %w", m.version, m.name, err)
		}

		_, err = db.Exec(`UPDATE _migration_version SET version = ?;`, m.version)
		if err != nil {
			return fmt.Errorf("Failed to update migration version to %d: %w", m.version, err)
		}
		slog.Info("Successfully applied migration", "version", m.version, "name", m.name)
	}
	return nil
}

func initializeMigrationVersion(db *sql.DB) error {
	version, err := getMigrationVersion(db)
	if err == nil {
		slog.Debug("Migration version table already initialized", "version", version)
		return nil
	}

	slog.Info("Initializing migration version table to -1")

	_, err = db.Exec(`CREATE TABLE _migration_version 
	(version INTEGER NOT NULL UNIQUE PRIMARY KEY);`)
	if err != nil {
		return fmt.Errorf("Failed to create migration version table: %w", err)
	}

	_, err = db.Exec(`INSERT INTO _migration_version VALUES (-1);`)
	if err != nil {
		return fmt.Errorf("Failed to initialize migration version: %w", err)
	}

	return nil
}

func getMigrationVersion(db *sql.DB) (int, error) {
	var version int
	err := db.QueryRow(`SELECT version FROM _migration_version LIMIT 1;`).Scan(&version)
	if err != nil {
		return 0, fmt.Errorf("Failed to get current migration version: %w", err)
	}
	return version, nil
}
