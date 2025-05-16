package database

import (
	"database/sql"
	"fmt"
)

type DatabaseContainer struct {
	db      *sql.DB
	dialect string
}

// New connects to the given SQL database and wraps it in a Container.
//
// Only SQLite and Postgres are currently fully supported.
//
// The logger can be nil and will default to a no-op logger.
//
// When using SQLite, it's strongly recommended to enable foreign keys by adding `?_foreign_keys=true`:
//
//	container, err := sqlstore.New("sqlite3", "file:yoursqlitefile.db?_foreign_keys=on", nil)
func New(dialect, address string) (*DatabaseContainer, error) {
	db, err := sql.Open(dialect, address)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	container := NewWithDB(db, dialect)
	err = container.Upgrade()
	if err != nil {
		return nil, fmt.Errorf("failed to upgrade database: %w", err)
	}
	return container, nil
}

func NewWithDB(db *sql.DB, dialect string) *DatabaseContainer {
	return &DatabaseContainer{
		db:      db,
		dialect: dialect,
	}
}

const (
	insertWebhookResponseQuery = `
		INSERT INTO app_callback_responses (callback_url,status,response,error_message)
		VALUES ($1, $2, $3, $4)
	`
)

func (c *DatabaseContainer) StoreResponse(res *AppWebhookResponse) error {
	_, err := c.db.Exec(insertWebhookResponseQuery,
		res.CallbackUrl,
		res.Status,
		res.Response,
		res.ErrorMessage,
	)
	return err
}
