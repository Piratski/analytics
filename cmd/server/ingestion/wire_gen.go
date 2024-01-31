// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package ingestion

import (
	"github.com/prismelabs/prismeanalytics/cmd/server/wired"
	"github.com/prismelabs/prismeanalytics/internal/clickhouse"
	"github.com/prismelabs/prismeanalytics/internal/handlers"
	"github.com/prismelabs/prismeanalytics/internal/middlewares"
	"github.com/prismelabs/prismeanalytics/internal/services/eventstore"
	"github.com/prismelabs/prismeanalytics/internal/services/ipgeolocator"
	"github.com/prismelabs/prismeanalytics/internal/services/sourceregistry"
	"github.com/prismelabs/prismeanalytics/internal/services/uaparser"
)

// Injectors from wire.go:

func Initialize(logger wired.BootstrapLogger) wired.App {
	server := wired.ProvideServerConfig(logger)
	eventsCors := middlewares.ProvideEventsCors()
	eventsRateLimiter := middlewares.ProvideEventsRateLimiter(server)
	logLogger := wired.ProvideLogger(server)
	accessLog := middlewares.ProvideAccessLog(server, logLogger)
	errorHandler := middlewares.ProvideErrorHandler()
	healhCheck := handlers.ProvideHealthCheck()
	middlewaresLogger := middlewares.ProvideLogger(logLogger)
	requestId := middlewares.ProvideRequestId(server)
	static := middlewares.ProvideStatic(server)
	minimalFiber := wired.ProvideMinimalFiber(accessLog, server, errorHandler, healhCheck, middlewaresLogger, requestId, static)
	configClickhouse := wired.ProvideClickhouseConfig(logger)
	ch := clickhouse.ProvideCh(logLogger, configClickhouse)
	service := eventstore.ProvideClickhouseService(ch, logLogger)
	sourceregistryService := sourceregistry.ProvideEnvVarService(logLogger)
	uaparserService := uaparser.ProvideService()
	ipgeolocatorService := ipgeolocator.ProvideMmdbService(logLogger)
	postPageViewEvent := handlers.ProvidePostEventsPageViews(service, sourceregistryService, uaparserService, ipgeolocatorService)
	app := ProvideFiber(eventsCors, eventsRateLimiter, minimalFiber, postPageViewEvent)
	setup := wired.ProvideSetup()
	wiredApp := wired.ProvideApp(server, app, logLogger, setup)
	return wiredApp
}
