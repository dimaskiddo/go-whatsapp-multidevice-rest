package auth

import (
	"github.com/dimaskiddo/go-whatsapp-multidevice-rest/pkg/env"
)

var AuthBasicUsername string
var AuthBasicPassword string

var AuthJWTSecret string
var AuthJWTExpiredHour int

func init() {
	AuthBasicUsername, _ = env.GetEnvString("AUTH_BASIC_USERNAME")
	AuthBasicPassword, _ = env.GetEnvString("AUTH_BASIC_PASSWORD")

	AuthJWTSecret, _ = env.GetEnvString("AUTH_JWT_SECRET")
	AuthJWTExpiredHour, _ = env.GetEnvInt("AUTH_JWT_EXPIRED_HOUR")
}
