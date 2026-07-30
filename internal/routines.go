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
			// Check Every Authenticated JID
			for jid, client := range pkgWhatsApp.WhatsAppClient {
				if client == nil || client.Store == nil || client.Store.ID == nil {
					continue
				}

				// Get Real JID from Datastore
				realJID := client.Store.ID.User

				// Mask JID for Logging Information
				maskJID := realJID[0:len(realJID)-4] + "xxxx"

				// Print Log Show Information of Device Checking
				log.Print(nil).Info("Checking WhatsApp Client for " + maskJID)

				// Check WhatsAppClient Registered JID with Authenticated MSISDN
				if jid != realJID {
					// Print Log Show Information to Force Log-out Device
					log.Print(nil).Info("Logging out WhatsApp Client for " + maskJID + " Due to Missmatch Authentication")

					// Logout WhatsAppClient Device
					_ = pkgWhatsApp.WhatsAppLogout(jid)
					delete(pkgWhatsApp.WhatsAppClient, jid)
					continue
				}

				// Connection health check
				isConnected := client.IsConnected()
				isLoggedIn := client.IsLoggedIn()

				if !isConnected && isLoggedIn {
					// Client has valid session but disconnected — try reconnecting
					log.Print(nil).Warnf("[%s] Client disconnected, attempting auto-reconnect...", maskJID)
					err := pkgWhatsApp.WhatsAppReconnect(jid)
					if err != nil {
						log.Print(nil).Errorf("[%s] Auto-reconnect failed: %v", maskJID, err)
					} else {
						log.Print(nil).Infof("[%s] Auto-reconnect successful", maskJID)
					}
				} else if !isConnected && !isLoggedIn {
					log.Print(nil).Warnf("[%s] Client not connected and not logged in (needs manual re-login)", maskJID)
				} else {
					log.Print(nil).Debugf("[%s] Client health OK", maskJID)
				}
			}
		}
	})

	// Additional health check every 30 seconds (on the 30s mark)
	cron.AddFunc("30 * * * * *", func() {
		if len(pkgWhatsApp.WhatsAppClient) > 0 {
			log.Print(nil).Debug("Running connection health check for all clients...")
			failedJIDs := pkgWhatsApp.WhatsAppCheckAllConnections()
			if len(failedJIDs) > 0 {
				log.Print(nil).Warnf("Health check: %d client(s) failed to reconnect", len(failedJIDs))
			}
		}
	})

	cron.Start()
}
