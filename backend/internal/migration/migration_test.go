package migration

import (
	"database/sql"
	"embed"
	"io/fs"
	"os"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations_test/valid
var validMigrations embed.FS

//go:embed migrations_test/invalid_version/*.sql
var invalidMigrations embed.FS

var expectedValidMigrations = []migration{
	{
		version: 0,
		name:    "create-user-table",
		script: `CREATE TABLE user (
    id INTEGER PRIMARY KEY NOT NULL UNIQUE,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL
);
`,
	},
	{
		version: 1,
		name:    "create-book-table",
		script: `CREATE TABLE book (
    id INTEGER PRIMARY KEY NOT NULL UNIQUE,
    title TEXT NOT NULL,
    author TEXT NOT NULL,
    published_year INTEGER,
    isbn TEXT UNIQUE,
    created_at TIMESTAMP NOT NULL
);
`,
	},
}

func TestLoadMigrations(t *testing.T) {
	tests := []struct {
		name           string
		scriptsFs      fs.FS
		migrationsRoot string
		expectedError  bool
		expectedCount  int
	}{
		{
			name:           "Valid migrations",
			scriptsFs:      validMigrations,
			migrationsRoot: "migrations_test/valid",
			expectedError:  false,
			expectedCount:  2,
		},
		{
			name:           "Invalid migration version",
			scriptsFs:      invalidMigrations,
			migrationsRoot: "migrations_test/invalid_version",
			expectedError:  true,
			expectedCount:  0,
		},
		{
			name:           "Invalid FS",
			scriptsFs:      os.DirFS("invalid"),
			migrationsRoot: "",
			expectedError:  true,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subFs := tt.scriptsFs
			var err error
			if tt.migrationsRoot != "" {
				subFs, err = fs.Sub(tt.scriptsFs, tt.migrationsRoot)
				if err != nil {
					t.Fatalf("Failed to create sub filesystem: %v", err)
				}
			}

			migrations, err := loadMigrations(subFs)
			if tt.expectedError {
				if err == nil {
					t.Fatalf("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if len(migrations) != tt.expectedCount {
				t.Fatalf("Expected %d migrations, got %d", tt.expectedCount, len(migrations))
			}

			for i, m := range migrations {
				em := expectedValidMigrations[i]
				if m.version != em.version || m.name != em.name || m.script != em.script {
					t.Errorf("Migration %d does not match expected.\nGot: %+v\nExpected: %+v", i, m, em)
				}
			}
		})
	}
}

func TestPerformMigrations(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}
	defer db.Close()

	err = initializeMigrationVersion(db)
	if err != nil {
		t.Fatalf("Failed to initialize migration version: %v", err)
	}

	err = performMigrations(db, expectedValidMigrations, -1)
	if err != nil {
		t.Fatalf("Failed to perform migrations: %v", err)
	}

	version, err := getMigrationVersion(db)
	if err != nil {
		t.Fatalf("Failed to get migration version: %v", err)
	}

	if version != 1 {
		t.Errorf("Expected migration version 1, got %d", version)
	}

	expectedValidMigrations = append(expectedValidMigrations, migration{
		version: 2,
		name:    "add-published-date-to-books",
		script:  `ALTER TABLE book ADD COLUMN published_date TEXT;`,
	})

	err = performMigrations(db, expectedValidMigrations, 1)
	if err != nil {
		t.Fatalf("Failed to perform migrations: %v", err)
	}

	version, err = getMigrationVersion(db)
	if err != nil {
		t.Fatalf("Failed to get migration version: %v", err)
	}

	if version != 2 {
		t.Errorf("Expected migration version 2, got %d", version)
	}
}

func TestPerformMigrations_InvalidScript(t *testing.T) {
	migrations := []migration{
		{
			version: 1,
			name:    "test",
			script:  "dummy",
		},
	}

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}
	defer db.Close()

	err = performMigrations(db, migrations, 0)
	if err == nil || !strings.Contains(err.Error(), "Failed to apply migration") {
		t.Errorf("Expected error")
	}
}

func TestInitializeMigrationVersion(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}
	defer db.Close()

	err = initializeMigrationVersion(db)
	if err != nil {
		t.Fatalf("Failed to initialize migration version: %v", err)
	}

	var version int
	err = db.QueryRow(`SELECT version FROM _migration_version;`).Scan(&version)
	if err != nil {
		t.Fatalf("Failed to query migration version: %v", err)
	}

	if version != -1 {
		t.Errorf("Expected initial migration version -1, got %d", version)
	}
}

func TestInitializeMigrationVersion_AlreadyInitialized(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}
	defer db.Close()

	err = initializeMigrationVersion(db)
	if err != nil {
		t.Fatalf("Failed to initialize migration version: %v", err)
	}

	err = initializeMigrationVersion(db)
	if err != nil {
		t.Fatalf("Expected no error")
	}
}

func TestInitiliazeMigrationVersion_CreateExists(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE _migration_version (aaa INTEGER)")
	if err != nil {
		t.Fatalf("Failed to initialize test db")
	}

	err = initializeMigrationVersion(db)
	if err == nil {
		t.Errorf("expected error")
	}
}

func TestGetMigrationVersion(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}
	defer db.Close()

	err = initializeMigrationVersion(db)
	if err != nil {
		t.Fatalf("Failed to initialize migration version: %v", err)
	}

	version, err := getMigrationVersion(db)
	if err != nil {
		t.Fatalf("Failed to get migration version: %v", err)
	}

	if version != -1 {
		t.Errorf("Expected migration version -1, got %d", version)
	}
}
