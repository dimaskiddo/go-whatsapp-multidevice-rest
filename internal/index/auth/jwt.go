package model

import (
	"github.com/golang-jwt/jwt"
)

type AuthJWTClaims struct {
	Data AuthJWTClaimsPayload `json:"dat"`
	jwt.StandardClaims
}

type AuthJWTClaimsPayload struct {
	MSISDN string `json:"msisdn"`
}
