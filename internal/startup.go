package internal

import (
	"github.com/dimaskiddo/go-whatsapp-multidevice-rest/pkg/log"
	pkgWhatsApp "github.com/dimaskiddo/go-whatsapp-multidevice-rest/pkg/whatsapp"
)

func Startup() {
	log.Print(nil).Info("Running Startup Tasks")

	devices, err := pkgWhatsApp.WhatsAppDatastore.GetAllDevices()
	if err != nil {
		log.Print(nil).Error("Failed to Load WhatsApp Client Devices from Datastore")
	}

	for _, device := range devices {
		jid := pkgWhatsApp.WhatsAppDecomposeJID(device.ID.String())
		maskJID := jid[0:len(jid)-4] + "xxxx"

		log.Print(nil).Info("Restoring WhatsApp Client for " + maskJID)

		err := pkgWhatsApp.WhatsAppInitClient(device, jid)
		if err != nil {
			log.Print(nil).Error(err.Error())
		}

		err = pkgWhatsApp.WhatsAppReconnect(jid)
		if err != nil {
			log.Print(nil).Error(err.Error())
		}
	}
}
