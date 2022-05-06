package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/go-playground/validator/v10"

	"github.com/dimaskiddo/go-whatsapp-multidevice-rest/pkg/env"
	"github.com/dimaskiddo/go-whatsapp-multidevice-rest/pkg/log"
	"github.com/dimaskiddo/go-whatsapp-multidevice-rest/pkg/router"

	"github.com/dimaskiddo/go-whatsapp-multidevice-rest/internal"
)

type Server struct {
	Address string
	Port    string
}

type EchoValidator struct {
	Validator *validator.Validate
}

func (ev *EchoValidator) Validate(i interface{}) error {
	return ev.Validator.Struct(i)
}

func main() {
	var err error

	// Initialize Echo
	e := echo.New()

	// Router Recovery
	e.Use(middleware.Recover())

	// Router Compression
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: router.GZipLevel,
	}))

	// Router CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{router.CORSOrigin},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE},
	}))

	// Router Security
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		ContentTypeNosniff: "nosniff",
		XSSProtection:      "1; mode=block",
		XFrameOptions:      "SAMEORIGIN",
	}))

	// Router Body Size Limit
	e.Use(middleware.BodyLimitWithConfig(middleware.BodyLimitConfig{
		Limit: router.BodyLimit,
	}))

	// Router Cache
	e.Use(router.HttpCacheInMemory(
		router.CacheCapacity,
		router.CacheTTLSeconds,
	))

	// Router RealIP
	e.Use(router.HttpRealIP())

	// Router Validator
	e.Validator = &EchoValidator{
		Validator: validator.New(),
	}

	// Router Default Handler
	e.HTTPErrorHandler = router.HttpErrorHandler
	e.GET("/favicon.ico", router.ResponseNoContent)

	// Load Internal Routes
	internal.Routes(e)

	// Running Startup Tasks
	internal.Startup()

	// Get Server Configuration
	var serverConfig Server

	serverConfig.Address, err = env.GetEnvString("SERVER_ADDRESS")
	if err != nil {
		serverConfig.Address = "127.0.0.1"
	}

	serverConfig.Port, err = env.GetEnvString("SERVER_PORT")
	if err != nil {
		serverConfig.Port = "3000"
	}

	// Start Server
	go func() {
		err := e.Start(serverConfig.Address + ":" + serverConfig.Port)
		if err != nil && err != http.ErrServerClosed {
			log.Print(nil).Fatal(err.Error())
		}
	}()

	// Watch for Shutdown Signal
	sigShutdown := make(chan os.Signal, 1)
	signal.Notify(sigShutdown, os.Interrupt)
	signal.Notify(sigShutdown, syscall.SIGTERM)
	<-sigShutdown

	// Wait 5 Seconds Before Graceful Shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Try To Shutdown Server
	err = e.Shutdown(ctx)
	if err != nil {
		log.Print(nil).Fatal(err.Error())
	}
}
