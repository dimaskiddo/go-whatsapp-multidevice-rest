package database

import (
	"database/sql"
	"fmt"
)

type upgradeFunc func(*sql.Tx, *DatabaseContainer) error

// Upgrades is a list of functions that will upgrade a database to the latest version.
//
// This may be of use if you want to manage the database fully manually, but in most cases you
// should just call DatabaseContainer.Upgrade to let the library handle everything.
var Upgrades = [...]upgradeFunc{upgradeV1}

const (
	AppMigrationTable = "app_migration_version"
)

func (c *DatabaseContainer) getVersion() (int, error) {
	_, err := c.db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (version INTEGER)", AppMigrationTable))
	if err != nil {
		return -1, err
	}

	version := 0
	row := c.db.QueryRow(fmt.Sprintf("SELECT version FROM %s LIMIT 1", AppMigrationTable))
	if row != nil {
		_ = row.Scan(&version)
	}
	return version, nil
}

func (c *DatabaseContainer) setVersion(tx *sql.Tx, version int) error {
	_, err := tx.Exec(fmt.Sprintf("DELETE FROM %s", AppMigrationTable))
	if err != nil {
		return err
	}
	_, err = tx.Exec(fmt.Sprintf("INSERT INTO %s (version) VALUES ($1)", AppMigrationTable), version)
	return err
}

func (c *DatabaseContainer) Upgrade() error {
	if c.dialect == "sqlite" {
		var foreignKeysEnabled bool
		err := c.db.QueryRow("PRAGMA foreign_keys").Scan(&foreignKeysEnabled)
		if err != nil {
			return fmt.Errorf("failed to check if foreign keys are enabled: %w", err)
		} else if !foreignKeysEnabled {
			return fmt.Errorf("foreign keys are not enabled")
		}
	}

	version, err := c.getVersion()
	if err != nil {
		return err
	}

	for ; version < len(Upgrades); version++ {
		var tx *sql.Tx
		tx, err = c.db.Begin()
		if err != nil {
			return err
		}

		migrateFunc := Upgrades[version]
		err = migrateFunc(tx, c)
		if err != nil {
			_ = tx.Rollback()
			return err
		}

		if err = c.setVersion(tx, version+1); err != nil {
			return err
		}

		if err = tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

func upgradeV1(tx *sql.Tx, c *DatabaseContainer) error {
	var err error
	postgresAppCallback := `CREATE TABLE app_callback_responses (
		id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
		callback_url TEXT NOT NULL,
		status VARCHAR(50) NOT NULL,
		response JSONB,
		error_message TEXT,
		timestamp TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);`
	sqliteAppCallback := `CREATE TABLE app_callback_responses (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		callback_url TEXT NOT NULL,
		status TEXT NOT NULL,
		response TEXT,
		error_message TEXT,
		timestamp TEXT DEFAULT CURRENT_TIMESTAMP
	);`
	if c.dialect == "postgres" {
		_, err = tx.Exec(postgresAppCallback)
	} else if c.dialect == "sqlite" {
		_, err = tx.Exec(sqliteAppCallback)
	}
	if err != nil {
		return err
	}
	return nil
}
