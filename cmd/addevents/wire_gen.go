// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/prismelabs/analytics/pkg/clickhouse"
	"github.com/prismelabs/analytics/pkg/wired"
)

import (
	_ "embed"
)

// Injectors from wire.go:

func Initialize(logger wired.BootstrapLogger) App {
	zerologLogger := ProvideLogger()
	config := ProvideConfig()
	configClickhouse := wired.ProvideClickhouseConfig(logger)
	driver := clickhouse.ProvideEmbeddedSourceDriver(zerologLogger)
	ch := clickhouse.ProvideCh(zerologLogger, configClickhouse, driver)
	app := ProvideApp(zerologLogger, config, ch)
	return app
}
