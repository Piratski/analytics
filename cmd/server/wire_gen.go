// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/prismelabs/prismeanalytics/internal/log"
	"github.com/prismelabs/prismeanalytics/internal/middlewares"
)

// Injectors from wire.go:

func initialize(logger log.Logger) App {
	config := ProvideConfig(logger)
	standardLogger := ProvideStandardLogger(config)
	accessLogger := ProvideAccessLogger(config, standardLogger)
	authMiddleware := middlewares.ProvideAuthMiddleware(config)
	echo := ProvideEcho(config, accessLogger, authMiddleware)
	app := ProvideApp(config, echo, standardLogger)
	return app
}
