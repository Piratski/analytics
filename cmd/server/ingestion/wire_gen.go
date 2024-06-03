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
	"github.com/prismelabs/analytics/pkg/services/saltmanager"
	"github.com/prismelabs/analytics/pkg/services/sessionstorage"
	"github.com/prismelabs/analytics/pkg/services/teardown"
	"github.com/prismelabs/analytics/pkg/services/uaparser"
	"github.com/prismelabs/analytics/pkg/wired"
)

// Injectors from wire.go:

func Initialize(logger wired.BootstrapLogger) wired.App {
	eventsCors := middlewares.ProvideEventsCors()
	server := wired.ProvideServerConfig(logger)
	storage := wired.ProvideFiberStorage()
	eventsRateLimiter := middlewares.ProvideEventsRateLimiter(server, storage)
	zerologLogger := wired.ProvideLogger(server)
	accessLog := middlewares.ProvideAccessLog(server, zerologLogger)
	errorHandler := middlewares.ProvideErrorHandler()
	config := wired.ProvideMinimalFiberConfig(server)
	healhCheck := handlers.ProvideHealthCheck()
	requestId := middlewares.ProvideRequestId(server)
	static := middlewares.ProvideStatic(server)
	registry := wired.ProvidePrometheusRegistry()
	metrics := middlewares.ProvideMetrics(registry)
	service := teardown.ProvideService()
	minimalFiber := wired.ProvideMinimalFiber(accessLog, errorHandler, config, healhCheck, zerologLogger, requestId, static, metrics, service)
	originregistryService := originregistry.ProvideEnvVarService(zerologLogger)
	nonRegisteredOriginFilter := middlewares.ProvideNonRegisteredOriginFilter(originregistryService)
	eventstoreConfig := eventstore.ProvideConfig()
	configClickhouse := wired.ProvideClickhouseConfig(logger)
	driver := clickhouse.ProvideEmbeddedSourceDriver(zerologLogger)
	ch := clickhouse.ProvideCh(zerologLogger, configClickhouse, driver)
	eventstoreService := eventstore.ProvideService(eventstoreConfig, ch, zerologLogger, registry, service)
	saltmanagerService := saltmanager.ProvideService(zerologLogger)
	sessionstorageConfig := sessionstorage.ProvideConfig()
	sessionstorageService := sessionstorage.ProvideService(zerologLogger, sessionstorageConfig)
	postEventsCustom := handlers.ProvidePostEventsCustom(zerologLogger, eventstoreService, saltmanagerService, sessionstorageService)
	uaparserService := uaparser.ProvideService(zerologLogger, registry)
	ipgeolocatorService := ipgeolocator.ProvideMmdbService(zerologLogger, registry)
	postEventsPageview := handlers.ProvidePostEventsPageViews(zerologLogger, eventstoreService, uaparserService, ipgeolocatorService, saltmanagerService, sessionstorageService)
	app := ProvideFiber(eventsCors, eventsRateLimiter, minimalFiber, nonRegisteredOriginFilter, postEventsCustom, postEventsPageview)
	promhttpLogger := wired.ProvidePromHttpLogger(server, zerologLogger)
	setup := wired.ProvideSetup()
	wiredApp := wired.ProvideApp(app, server, zerologLogger, promhttpLogger, registry, setup, service)
	return wiredApp
}
