// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package ingestion

import (
	"github.com/prismelabs/analytics/pkg/clickhouse"
	"github.com/prismelabs/analytics/pkg/handlers"
	"github.com/prismelabs/analytics/pkg/middlewares"
	"github.com/prismelabs/analytics/pkg/services/eventstore"
	"github.com/prismelabs/analytics/pkg/services/ipgeolocator"
	"github.com/prismelabs/analytics/pkg/services/originregistry"
	"github.com/prismelabs/analytics/pkg/services/teardown"
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
	requestId := middlewares.ProvideRequestId(server)
	static := middlewares.ProvideStatic(server)
	service := teardown.ProvideService()
	minimalFiber := wired.ProvideMinimalFiber(accessLog, errorHandler, config, healhCheck, zerologLogger, requestId, static, service)
	originregistryService := originregistry.ProvideEnvVarService(zerologLogger)
	nonRegisteredOriginFilter := middlewares.ProvideNonRegisteredOriginFilter(originregistryService)
	configClickhouse := wired.ProvideClickhouseConfig(logger)
	driver := clickhouse.ProvideEmbeddedSourceDriver(zerologLogger)
	ch := clickhouse.ProvideCh(zerologLogger, configClickhouse, driver)
	eventstoreService := eventstore.ProvideClickhouseService(ch, zerologLogger, service)
	uaparserService := uaparser.ProvideService(zerologLogger)
	ipgeolocatorService := ipgeolocator.ProvideMmdbService(zerologLogger)
	postEventsCustom := handlers.ProvidePostEventsCustom(eventstoreService, uaparserService, ipgeolocatorService)
	postEventsPageview := handlers.ProvidePostEventsPageViews(zerologLogger, eventstoreService, uaparserService, ipgeolocatorService)
	app := ProvideFiber(eventsCors, eventsRateLimiter, minimalFiber, nonRegisteredOriginFilter, postEventsCustom, postEventsPageview)
	setup := wired.ProvideSetup()
	wiredApp := wired.ProvideApp(server, app, zerologLogger, service, setup)
	return wiredApp
}
