package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/prismelabs/prismeanalytics/internal/config"
	"github.com/prismelabs/prismeanalytics/internal/log"
	"github.com/prismelabs/prismeanalytics/internal/middlewares"
)

func main() {
	cfg := config.FromEnv()

	logger := log.NewLogger("app", os.Stderr, cfg.Server.Debug)
	accessLogger := log.NewLogger("access_log", os.Stderr, cfg.Server.Debug)
	log.TestLoggers(logger)

	e := echo.New()
	if cfg.Server.TrustProxy {
		e.IPExtractor = echo.ExtractIPFromXFFHeader()
	} else {
		e.IPExtractor = echo.ExtractIPDirect()
	}
	e.HideBanner = true
	e.HidePort = true

	e.Use(middlewares.RequestId(cfg.Server))
	e.Use(middlewares.AccessLog(accessLogger))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	socket := "0.0.0.0:8000"
	logger.Info().Msgf("start listening for incoming requests on http://%v", socket)
	logger.Panic().Err(e.Start(socket)).Send()
}
