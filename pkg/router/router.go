package router

import (
	"github.com/dimaskiddo/go-whatsapp-multidevice-rest/pkg/env"
)

var BaseURL, CORSOrigin, BodyLimit string
var GZipLevel int

func init() {
	var err error

	BaseURL, err = env.GetEnvString("API_BASE_URL")
	if err != nil {
		BaseURL = "/"
	}

	CORSOrigin, err = env.GetEnvString("HTTP_CORS_ORIGIN")
	if err != nil {
		CORSOrigin = "*"
	}

	BodyLimit, err = env.GetEnvString("HTTP_BODY_LIMIT_SIZE")
	if err != nil {
		BodyLimit = "8M"
	}

	GZipLevel, err = env.GetEnvInt("HTTP_GZIP_LEVEL")
	if err != nil {
		GZipLevel = 3
	}
}
