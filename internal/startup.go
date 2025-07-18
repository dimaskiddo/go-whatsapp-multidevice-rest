package internal

import (
	"github.com/dimaskiddo/go-whatsapp-multidevice-rest/internal/util"
	"github.com/dimaskiddo/go-whatsapp-multidevice-rest/pkg/log"
	pkgWhatsApp "github.com/dimaskiddo/go-whatsapp-multidevice-rest/pkg/whatsapp"
)

func Startup() {
	log.Print(nil).Info("Running Startup Tasks")

	// Load All WhatsApp Client Devices from Datastore
	devices, err := pkgWhatsApp.WhatsAppDatastore.GetAllDevices()
	if err != nil {
		log.Print(nil).Error("Failed to Load WhatsApp Client Devices from Datastore")
	}

	// Do Reconnect for Every Device in Datastore
	for _, device := range devices {
		// Get JID from Datastore
		jid := pkgWhatsApp.WhatsAppDecomposeJID(device.ID.User)

		// Mask JID for Logging Information
		maskedJID := util.MaskedJID(jid)

		// Print Restore Log
		log.Print(nil).Info("Restoring WhatsApp Client for " + maskedJID)

		// Initialize WhatsApp Client
		pkgWhatsApp.WhatsAppInitClient(device, jid)

		// Reconnect WhatsApp Client WebSocket
		err = pkgWhatsApp.WhatsAppReconnect(jid)
		if err != nil {
			log.Print(nil).Error(err.Error())
		}
	}
}
