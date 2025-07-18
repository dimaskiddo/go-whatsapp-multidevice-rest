package whatsapp

import "database/sql"

type whatsAppRepository struct {
	DB *sql.DB
}

type WhatsAppRepository interface {
	SetWebhook(jid string, webhookURL string) error
	GetWebhook(jid string) (string, error)
	DeleteWebhook(jid string) error
}

func NewWhatsappRepository(db *sql.DB) WhatsAppRepository {
	return &whatsAppRepository{DB: db}
}

func (r *whatsAppRepository) SetWebhook(jid string, webhookURL string) error {
	_, err := r.DB.Exec(`
		INSERT INTO device_webhooks (jid, webhook_url)
		VALUES ($1, $2)
		ON CONFLICT (jid)
		DO UPDATE SET webhook_url = EXCLUDED.webhook_url, updated_at = CURRENT_TIMESTAMP
	`, jid, webhookURL)
	return err
}

func (r *whatsAppRepository) GetWebhook(jid string) (string, error) {
	var webhookURL string
	err := r.DB.QueryRow("SELECT webhook_url FROM device_webhooks WHERE jid = $1", jid).Scan(&webhookURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return webhookURL, nil
}

func (r *whatsAppRepository) DeleteWebhook(jid string) error {
	_, err := r.DB.Exec("DELETE FROM device_webhooks WHERE jid = $1", jid)
	return err
}
