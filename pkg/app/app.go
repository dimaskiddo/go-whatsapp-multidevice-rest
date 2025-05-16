package app

import (
	"log"
	"time"

	"github.com/dimaskiddo/go-whatsapp-multidevice-rest/pkg/app/database"
	"github.com/dimaskiddo/go-whatsapp-multidevice-rest/pkg/app/http"
	"github.com/dimaskiddo/go-whatsapp-multidevice-rest/pkg/env"
)

var (
	AppWebhookURL       string
	AppWebhookBasicAuth string
	AppDatabase         *database.DatabaseContainer
	AppRequest          *http.HttpClient
)

func init() {
	var err error

	dbType, err := env.GetEnvString("WHATSAPP_DATASTORE_TYPE")
	if err != nil {
		log.Fatal("Error Parse Environment Variable for Application Datastore Type")
	}

	dbURI, err := env.GetEnvString("WHATSAPP_DATASTORE_URI")
	if err != nil {
		log.Fatal("Error Parse Environment Variable for Application Datastore URI")
	}

	// Initialize App Client Datastore
	initDB(dbType, dbURI)

	appWebhookUrl, err := env.GetEnvString("APP_WEBHOOK_URL_TARGET")
	if err != nil {
		log.Fatal("Error Parse Environment Variable for App Webhook URL Target")
	}
	AppWebhookURL = appWebhookUrl

	appWebhookBasicAuth, err := env.GetEnvString("APP_WEBHOOK_BASIC_AUTH")
	if err != nil {
		AppWebhookBasicAuth = ""
	}
	AppWebhookBasicAuth = appWebhookBasicAuth

	// Initialize App HTTP Request
	initHttpRequest()
}

func initDB(dbType string, dbURI string) {
	// Initialize App Client Datastore
	appDb, err := database.New(dbType, dbURI)
	if err != nil {
		log.Fatal("Error Connect Application Datastore: ", err)
	}
	AppDatabase = appDb
}

func initHttpRequest() {
	// Initialize App HTTP Request
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	if AppWebhookBasicAuth != "" {
		headers["Authorization"] = "Basic " + AppWebhookBasicAuth
	}

	client := http.NewHttpClient(http.HttpClientOptions{
		Timeout: 30 * time.Second,
		Headers: headers,
	})

	AppRequest = client
}
