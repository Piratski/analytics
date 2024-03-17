//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/prismelabs/analytics/pkg/clickhouse"
	"github.com/prismelabs/analytics/pkg/wired"
)

func Initialize(logger wired.BootstrapLogger) App {
	wire.Build(
		ProvideApp,
		ProvideConfig,
		ProvideLogger,
		wired.ProvideClickhouseConfig,
		clickhouse.ProvideCh,
		clickhouse.ProvideEmbeddedSourceDriver,
	)

	return App{}
}
