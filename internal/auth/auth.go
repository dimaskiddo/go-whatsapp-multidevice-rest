package auth

import (
	"encoding/json"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"

	typAuth "github.com/dimaskiddo/go-whatsapp-multidevice-rest/internal/auth/types"

	"github.com/dimaskiddo/go-whatsapp-multidevice-rest/pkg/auth"
	"github.com/dimaskiddo/go-whatsapp-multidevice-rest/pkg/router"
)

// Auth
// @Summary     Generate Authentication Token
// @Description Get Authentication Token
// @Tags        Root
// @Produce     json
// @Success     200
// @Security    BasicAuth
// @Router      /auth [get]
func Auth(c echo.Context) error {
	var reqAuthBasicInfo typAuth.RequestAuthBasicInfo
	var resAuthJWTData typAuth.ResponseAuthJWTData

	// Parse Basic Auth Information from Rewrited Body Request
	// By Basic Auth Middleware
	_ = json.NewDecoder(c.Request().Body).Decode(&reqAuthBasicInfo)

	// Create JWT Claims
	var jwtClaims *typAuth.AuthJWTClaims
	if auth.AuthJWTExpiredHour > 0 {
		jwtClaims = &typAuth.AuthJWTClaims{
			typAuth.AuthJWTClaimsPayload{
				JID: reqAuthBasicInfo.Username,
			},
			jwt.StandardClaims{
				IssuedAt:  time.Now().Unix(),
				ExpiresAt: time.Now().Add(time.Hour * time.Duration(auth.AuthJWTExpiredHour)).Unix(),
			},
		}
	} else {
		jwtClaims = &typAuth.AuthJWTClaims{
			typAuth.AuthJWTClaimsPayload{
				JID: reqAuthBasicInfo.Username,
			},
			jwt.StandardClaims{
				IssuedAt: time.Now().Unix(),
			},
		}
	}

	// Create JWT Token
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)

	// Generate Encoded JWT Token
	jwtTokenEncoded, err := jwtToken.SignedString([]byte(auth.AuthJWTSecret))
	if err != nil {
		return router.ResponseInternalError(c, "")
	}

	// Set Encoded JWT Token as Response Data
	resAuthJWTData.Token = jwtTokenEncoded

	// Return JWT Token in JSON Response
	return router.ResponseSuccessWithData(c, "Successfully Authenticated", resAuthJWTData)
}
