package auth

import (
	"github.com/dimaskiddo/go-whatsapp-multidevice-rest/pkg/env"
)

var AuthBasicUsername string
var AuthBasicPassword string

var AuthJWTSecret string
var AuthJWTExpiredHour int

func init() {
	var err error

	AuthBasicUsername, err = env.GetEnvString("AUTH_BASIC_USERNAME")
	if err != nil {
		AuthBasicUsername = "administrator"
	}

	AuthBasicPassword, err = env.GetEnvString("AUTH_BASIC_PASSWORD")
	if err != nil {
		AuthBasicPassword = "83e4060e-78e1-4fe5-9977-aeeccd46a2b8"
	}

	AuthJWTSecret, err = env.GetEnvString("AUTH_JWT_SECRET")
	if err != nil {
		AuthJWTSecret = "9e4eb4cf-be25-4a29-bba3-fefb5a30f6ab"
	}

	AuthJWTExpiredHour, err = env.GetEnvInt("AUTH_JWT_EXPIRED_HOUR")
	if err != nil {
		AuthJWTExpiredHour = 24
	}
}
