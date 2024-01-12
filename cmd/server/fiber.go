package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prismelabs/prismeanalytics/internal/config"
	"github.com/prismelabs/prismeanalytics/internal/handlers"
	"github.com/prismelabs/prismeanalytics/internal/middlewares"
)

// ProvideFiber is a wire provider for fiber.App.
func ProvideFiber(
	cfg config.Config,
	viewsEngine fiber.Views,
	loggerMiddleware middlewares.Logger,
	accessLogMiddleware middlewares.AccessLog,
	requestIdMiddleware middlewares.RequestId,
	staticMiddleware middlewares.Static,
) *fiber.App {
	fiberCfg := fiber.Config{
		ServerHeader:          "prisme",
		StrictRouting:         true,
		AppName:               "Prisme Analytics",
		DisableStartupMessage: true,
		ErrorHandler: func(_ *fiber.Ctx, _ error) error {
			// Errors are handled manually by a middleware.
			return nil
		},
		Views:       viewsEngine,
		ViewsLayout: "layouts/empty",
	}
	if cfg.Server.TrustProxy {
		fiberCfg.EnableIPValidation = false
		fiberCfg.ProxyHeader = fiber.HeaderXForwardedFor
	} else {
		fiberCfg.EnableIPValidation = true
		fiberCfg.ProxyHeader = ""
	}

	app := fiber.New(fiberCfg)

	app.Use(fiber.Handler(requestIdMiddleware))
	app.Use(fiber.Handler(accessLogMiddleware))
	app.Use(fiber.Handler(loggerMiddleware))

	app.Use("/static", fiber.Handler(staticMiddleware))

	app.Get("/sign_up", handlers.GetSignUp)
	app.Post("/sign_up", handlers.PostSignUp)

	return app
}
