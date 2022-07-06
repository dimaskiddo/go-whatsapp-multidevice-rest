package internal

import (
	"github.com/robfig/cron/v3"

	"github.com/dimaskiddo/go-whatsapp-multidevice-rest/pkg/log"
	pkgWhatsApp "github.com/dimaskiddo/go-whatsapp-multidevice-rest/pkg/whatsapp"
)

func Routines(cron *cron.Cron) {
	log.Print(nil).Info("Running Routine Tasks")

	cron.AddFunc("0 * * * * *", func() {
		// If WhatsAppClient Connection is more than 0
		if len(pkgWhatsApp.WhatsAppClient) > 0 {
			// Check Every Authenticated MSISDN
			for jid, client := range pkgWhatsApp.WhatsAppClient {
				// Get Real JID from Datastore
				realJID := client.Store.ID.User

				// Check WhatsAppClient Registered JID with Authenticated MSISDN
				if jid != realJID {
					// Mask JID for Logging Information
					maskJID := realJID[0:len(realJID)-4] + "xxxx"

					log.Print(nil).Info("Logging out WhatsApp Client for " + maskJID + " Due to Missmatch Authentication")

					// Logout WhatsAppClient Device
					_ = pkgWhatsApp.WhatsAppLogout(jid)
					delete(pkgWhatsApp.WhatsAppClient, jid)
				}
			}
		}
	})

	cron.Start()
}
