package internal

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/dimaskiddo/go-whatsapp-multidevice-rest/pkg/log"
	pkgWhatsApp "github.com/dimaskiddo/go-whatsapp-multidevice-rest/pkg/whatsapp"
)

func Startup() {
	log.Print(nil).Info("Running Startup Tasks")

	dbs, err := filepath.Glob("./dbs/*.db")
	if err != nil {
		log.Print(nil).Error("Error to Get Existing SQLite Database Files")
	}

	for _, db := range dbs {
		jid := strings.TrimSuffix(filepath.Base(db), path.Ext(db))
		maskJID := jid[0:len(jid)-4] + "xxxx"

		log.Print(nil).Info("Restoring WhatsApp Client for " + maskJID)

		err := pkgWhatsApp.WhatAppConnect(jid)
		if err != nil {
			log.Print(nil).Error(err.Error())
		}

		err = pkgWhatsApp.WhatsAppReconnect(jid)
		if err != nil {
			log.Print(nil).Error(err.Error())
		}
	}
}
