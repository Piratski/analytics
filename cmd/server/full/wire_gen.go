// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package full

import (
	"github.com/prismelabs/prismeanalytics/cmd/server/wired"
	"github.com/prismelabs/prismeanalytics/internal/clickhouse"
	"github.com/prismelabs/prismeanalytics/internal/grafana"
	"github.com/prismelabs/prismeanalytics/internal/handlers"
	"github.com/prismelabs/prismeanalytics/internal/middlewares"
	"github.com/prismelabs/prismeanalytics/internal/postgres"
	"github.com/prismelabs/prismeanalytics/internal/services/auth"
	"github.com/prismelabs/prismeanalytics/internal/services/eventstore"
	grafana2 "github.com/prismelabs/prismeanalytics/internal/services/grafana"
	"github.com/prismelabs/prismeanalytics/internal/services/ipgeolocator"
	"github.com/prismelabs/prismeanalytics/internal/services/sessions"
	"github.com/prismelabs/prismeanalytics/internal/services/sourceregistry"
	"github.com/prismelabs/prismeanalytics/internal/services/uaparser"
	"github.com/prismelabs/prismeanalytics/internal/services/users"
)

// Injectors from wire.go:

func Initialize(logger wired.BootstrapLogger) wired.App {
	server := wired.ProvideServerConfig(logger)
	eventsCors := middlewares.ProvideEventsCors()
	eventsRateLimiter := middlewares.ProvideEventsRateLimiter(server)
	favicon := middlewares.ProvideFavicon()
	getIndex := handlers.ProvideGetIndex()
	getSignIn := handlers.ProvideGetSignIn()
	getSignUp := handlers.ProvideGetSignUp()
	views := wired.ProvideFiberViewsEngine(server)
	logLogger := wired.ProvideLogger(server)
	middlewaresLogger := middlewares.ProvideLogger(logLogger)
	accessLog := middlewares.ProvideAccessLog(server, logLogger)
	requestId := middlewares.ProvideRequestId(server)
	static := middlewares.ProvideStatic(server)
	minimalFiber := wired.ProvideMinimalFiber(server, views, middlewaresLogger, accessLog, requestId, static)
	notFound := handlers.ProvideNotFound()
	configClickhouse := wired.ProvideClickhouseConfig(logger)
	ch := clickhouse.ProvideCh(logLogger, configClickhouse)
	service := eventstore.ProvideClickhouseService(ch, logLogger)
	sourceregistryService := sourceregistry.ProvideEnvVarService(logLogger)
	uaparserService := uaparser.ProvideService()
	ipgeolocatorService := ipgeolocator.ProvideMmdbService(logLogger)
	postPageViewEvent := handlers.ProvidePostEventsPageViews(service, sourceregistryService, uaparserService, ipgeolocatorService)
	configPostgres := wired.ProvidePostgresConfig(logger)
	pg := postgres.ProvidePg(logLogger, configPostgres)
	usersService := users.ProvideService(pg)
	authService := auth.ProvideService(usersService)
	sessionsService := sessions.ProvideService()
	postSignIn := handlers.ProvidePostSignIn(authService, sessionsService)
	postSignUp := handlers.ProvidePostSignUp(usersService, sessionsService)
	withSession := middlewares.ProvideWithSession(sessionsService)
	app := ProvideFiber(eventsCors, eventsRateLimiter, favicon, getIndex, getSignIn, getSignUp, minimalFiber, notFound, postPageViewEvent, postSignIn, postSignUp, withSession)
	configGrafana := wired.ProvideGrafanaConfig(logger)
	client := grafana.ProvideClient(configGrafana)
	grafanaService := grafana2.ProvideService(client, configClickhouse)
	setup := ProvideSetup(logLogger, grafanaService)
	wiredApp := wired.ProvideApp(server, app, logLogger, setup)
	return wiredApp
}
