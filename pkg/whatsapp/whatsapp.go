package whatsapp

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	qrCode "github.com/skip2/go-qrcode"
	"google.golang.org/protobuf/proto"

	"go.mau.fi/whatsmeow"
	waproto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
)

var WhatsAppClient = make(map[string]*whatsmeow.Client)

func WhatsAppInit(jid string) (*whatsmeow.Client, error) {
	// Prepare SQLite Database File and Connection Address
	dbFileName := "dbs/" + jid + ".db"
	dbAddress := fmt.Sprintf("file:%s?_foreign_keys=on", dbFileName)

	// Create and Connect to SQLite Database
	datastore, err := sqlstore.New("sqlite3", dbAddress, nil)
	if err != nil {
		return nil, errors.New("Failed to Connect SQLite Database")
	}

	// Get First WhatsApp Device from SQLite Database
	device, err := datastore.GetFirstDevice()
	if err != nil {
		return nil, errors.New("Failed to Load WhatsApp Device")
	}

	// Set Client Properties
	store.CompanionProps.Os = proto.String("Go WhatsApp MultiDevice REST")
	store.CompanionProps.PlatformType = waproto.CompanionProps_DESKTOP.Enum()

	// Create New Client Connection
	client := whatsmeow.NewClient(device, nil)

	// Return Client Connection
	return client, nil
}

func WhatAppConnect(jid string) error {
	if WhatsAppClient[jid] == nil {
		// Initialize New WhatsApp Client
		client, err := WhatsAppInit(jid)
		if err != nil {
			return err
		}

		// Set Created WhatsApp Client to Map
		WhatsAppClient[jid] = client
	}

	return nil
}

func WhatsAppGenerateQR(qrChan <-chan whatsmeow.QRChannelItem) (string, string) {
	qrChanCode := make(chan string)
	qrChanTimeout := make(chan string)
	qrChanBase64 := make(chan string)

	// Get QR Code Data and Timeout
	go func() {
		for evt := range qrChan {
			if evt.Event == "code" {
				qrChanCode <- evt.Code
				qrChanTimeout <- evt.Timeout.String()
			}
		}
	}()

	// Generate QR Code Data to PNG Base64 Format
	go func() {
		select {
		case tmp := <-qrChanCode:
			png, _ := qrCode.Encode(tmp, qrCode.Medium, 256)
			qrChanBase64 <- base64.StdEncoding.EncodeToString(png)
		}
	}()

	// Return QR Code and Timeout Information
	return <-qrChanBase64, <-qrChanTimeout
}

func WhatsAppLogin(jid string) (string, string, error) {
	if WhatsAppClient[jid] != nil {
		// Make Sure WebSocket Connection is Disconnected
		WhatsAppClient[jid].Disconnect()

		if WhatsAppClient[jid].Store.ID == nil {
			// Device ID is not Exist
			// Generate QR Code
			qrChanGenerate, _ := WhatsAppClient[jid].GetQRChannel(context.Background())

			// Connect WebSocket while Initialize QR Code Data to be Sent
			err := WhatsAppClient[jid].Connect()
			if err != nil {
				return "", "", err
			}

			// Get Generated QR Code and Timeout Information
			qrImage, qrTimeout := WhatsAppGenerateQR(qrChanGenerate)

			// Return QR Code in Base64 Format and Timeout Information
			return "data:image/png;base64," + qrImage, qrTimeout, nil
		} else {
			// Device ID is Exist
			// Reconnect WebSocket
			err := WhatsAppClient[jid].Connect()
			if err != nil {
				return "", "", err
			}
		}
	}

	// Return Error WhatsApp Client is Not Valid
	return "", "", errors.New("WhatsApp Client is not Valid")
}

func WhatsAppLogout(jid string) error {
	if WhatsAppClient[jid] != nil {
		// Logout WhatsApp Client and Disconnect from WebSocket
		err := WhatsAppClient[jid].Logout()
		if err != nil {
			return err
		}

		// Remove SQLite Database File
		_ = os.Remove("dbs/" + jid + ".db")

		// Free WhatsApp Client Map
		WhatsAppClient[jid] = nil
		delete(WhatsAppClient, jid)

		return nil
	}

	// Return Error WhatsApp Client is Not Valid
	return errors.New("WhatsApp Client is not Valid")
}

func WhatsAppCreateUserJID(jid string) types.JID {
	return types.NewJID(jid, types.DefaultUserServer)
}

func WhatsAppCreateGroupJID(gjid string) types.JID {
	return types.NewJID(gjid, types.GroupServer)
}

func WhatsAppSendText(jid string, rjid string, message string) error {
	if WhatsAppClient[jid] != nil {
		if WhatsAppClient[jid].IsConnected() && WhatsAppClient[jid].IsLoggedIn() {
			_, err := WhatsAppClient[jid].SendMessage(WhatsAppCreateUserJID(rjid), "", &waproto.Message{
				Conversation: proto.String(message),
			})

			if err != nil {
				return err
			}

			return nil
		}
	}

	// Return Error WhatsApp Client is Not Valid
	return errors.New("WhatsApp Client is not Valid")
}
