package database

import "encoding/json"

type AppWebhookResponse struct {
	CallbackUrl  string          `json:"callback_url"`
	Status       string          `json:"status"`
	Response     json.RawMessage `json:"response"`
	ErrorMessage string          `json:"error_message"`
}
