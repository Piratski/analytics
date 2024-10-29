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
	server := wired.ProvideServerConfig(logger)
	apiEventsTimeout := middlewares.ProvideApiEventsTimeout(server)
	eventsCors := middlewares.ProvideEventsCors()
	storage := wired.ProvideFiberStorage()
	eventsRateLimiter := middlewares.ProvideEventsRateLimiter(server, storage)
	config := eventstore.ProvideConfig()
	zerologLogger := wired.ProvideLogger(server)
	configClickhouse := wired.ProvideClickhouseConfig(logger)
	driver := clickhouse.ProvideEmbeddedSourceDriver(zerologLogger)
	ch := clickhouse.ProvideCh(zerologLogger, configClickhouse, driver)
	registry := wired.ProvidePrometheusRegistry()
	service := teardown.ProvideService()
	eventstoreService := eventstore.ProvideService(config, ch, zerologLogger, registry, service)
	saltmanagerService := saltmanager.ProvideService(zerologLogger)
	sessionstorageConfig := sessionstorage.ProvideConfig()
	sessionstorageService := sessionstorage.ProvideService(zerologLogger, sessionstorageConfig, registry)
	getNoscriptEventsCustom := handlers.ProvideGetNoscriptEventsCustom(eventstoreService, saltmanagerService, sessionstorageService)
	uaparserService := uaparser.ProvideService(zerologLogger, registry)
	ipgeolocatorService := ipgeolocator.ProvideMmdbService(zerologLogger, registry)
	getNoscriptEventsPageviews := handlers.ProvideGetNoscriptEventsPageviews(zerologLogger, eventstoreService, uaparserService, ipgeolocatorService, saltmanagerService, sessionstorageService)
	accessLog := middlewares.ProvideAccessLog(server, zerologLogger)
	errorHandler := middlewares.ProvideErrorHandler()
	fiberConfig := wired.ProvideMinimalFiberConfig(server)
	healhCheck := handlers.ProvideHealthCheck()
	requestId := middlewares.ProvideRequestId(server)
	static := middlewares.ProvideStatic(server)
	metrics := middlewares.ProvideMetrics(registry)
	minimalFiber := wired.ProvideMinimalFiber(accessLog, errorHandler, fiberConfig, healhCheck, zerologLogger, requestId, static, metrics, service)
	originregistryService := originregistry.ProvideEnvVarService(zerologLogger)
	nonRegisteredOriginFilter := middlewares.ProvideNonRegisteredOriginFilter(originregistryService)
	noscriptHandlersCache := middlewares.ProvideNoscriptHandlersCache()
	postEventsCustom := handlers.ProvidePostEventsCustom(eventstoreService, saltmanagerService, sessionstorageService)
	postEventsClicksFileDownload := handlers.ProvidePostEventsClicksFileDownload(eventstoreService, saltmanagerService, sessionstorageService)
	postEventsClicksOutboundLink := handlers.ProvidePostEventsClicksOutboundLink(eventstoreService, saltmanagerService, sessionstorageService)
	postEventsPageviews := handlers.ProvidePostEventsPageViews(zerologLogger, eventstoreService, uaparserService, ipgeolocatorService, saltmanagerService, sessionstorageService)
	app := wired.ProvideFiber(apiEventsTimeout, eventsCors, eventsRateLimiter, getNoscriptEventsCustom, getNoscriptEventsPageviews, minimalFiber, nonRegisteredOriginFilter, noscriptHandlersCache, postEventsCustom, postEventsClicksFileDownload, postEventsClicksOutboundLink, postEventsPageviews)
	promhttpLogger := wired.ProvidePromHttpLogger(server, zerologLogger)
	setup := wired.ProvideSetup()
	wiredApp := wired.ProvideApp(app, server, zerologLogger, promhttpLogger, registry, setup, service)
	return wiredApp
}
