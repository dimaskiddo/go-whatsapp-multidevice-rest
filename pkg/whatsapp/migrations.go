package whatsapp

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func runDeviceWebhooksMigrations(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS device_webhooks (
        jid TEXT PRIMARY KEY REFERENCES whatsmeow_device(jid) ON DELETE CASCADE,
        webhook_url TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`
	_, err := db.ExecContext(context.Background(), query)
	return err
}
