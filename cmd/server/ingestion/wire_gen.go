// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package ingestion

import (
	"github.com/prismelabs/analytics/pkg/clickhouse"
	"github.com/prismelabs/analytics/pkg/handlers"
	"github.com/prismelabs/analytics/pkg/middlewares"
	"github.com/prismelabs/analytics/pkg/services/eventstore"
	"github.com/prismelabs/analytics/pkg/services/ipgeolocator"
	"github.com/prismelabs/analytics/pkg/services/sourceregistry"
	"github.com/prismelabs/analytics/pkg/services/uaparser"
	"github.com/prismelabs/analytics/pkg/wired"
)

// Injectors from wire.go:

func Initialize(logger wired.BootstrapLogger) wired.App {
	server := wired.ProvideServerConfig(logger)
	eventsCors := middlewares.ProvideEventsCors()
	eventsRateLimiter := middlewares.ProvideEventsRateLimiter(server)
	zerologLogger := wired.ProvideLogger(server)
	accessLog := middlewares.ProvideAccessLog(server, zerologLogger)
	errorHandler := middlewares.ProvideErrorHandler()
	config := wired.ProvideMinimalFiberConfig(server)
	healhCheck := handlers.ProvideHealthCheck()
	middlewaresLogger := middlewares.ProvideLogger(zerologLogger)
	requestId := middlewares.ProvideRequestId(server)
	static := middlewares.ProvideStatic(server)
	minimalFiber := wired.ProvideMinimalFiber(accessLog, errorHandler, config, healhCheck, middlewaresLogger, requestId, static)
	configClickhouse := wired.ProvideClickhouseConfig(logger)
	driver := clickhouse.ProvideEmbeddedSourceDriver(zerologLogger)
	ch := clickhouse.ProvideCh(zerologLogger, configClickhouse, driver)
	service := eventstore.ProvideClickhouseService(ch, zerologLogger)
	sourceregistryService := sourceregistry.ProvideEnvVarService(zerologLogger)
	uaparserService := uaparser.ProvideService(zerologLogger)
	ipgeolocatorService := ipgeolocator.ProvideMmdbService(zerologLogger)
	postEventsPageview := handlers.ProvidePostEventsPageViews(service, sourceregistryService, uaparserService, ipgeolocatorService)
	postEventsCustom := handlers.ProvidePostEventsCustom(service, sourceregistryService, uaparserService, ipgeolocatorService)
	app := ProvideFiber(eventsCors, eventsRateLimiter, minimalFiber, postEventsPageview, postEventsCustom)
	setup := wired.ProvideSetup()
	wiredApp := wired.ProvideApp(server, app, zerologLogger, setup)
	return wiredApp
}
