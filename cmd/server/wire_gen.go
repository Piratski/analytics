// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/prismelabs/prismeanalytics/internal/clickhouse"
	"github.com/prismelabs/prismeanalytics/internal/handlers"
	"github.com/prismelabs/prismeanalytics/internal/middlewares"
	"github.com/prismelabs/prismeanalytics/internal/postgres"
	"github.com/prismelabs/prismeanalytics/internal/services/auth"
	"github.com/prismelabs/prismeanalytics/internal/services/eventstore"
	"github.com/prismelabs/prismeanalytics/internal/services/sessions"
	"github.com/prismelabs/prismeanalytics/internal/services/sourceregistry"
	"github.com/prismelabs/prismeanalytics/internal/services/uaparser"
	"github.com/prismelabs/prismeanalytics/internal/services/users"
)

// Injectors from wire.go:

func initialize(logger BootstrapLogger) App {
	config := ProvideConfig(logger)
	views := ProvideFiberViewsEngine(config)
	logLogger := ProvideLogger(config)
	middlewaresLogger := middlewares.ProvideLogger(logLogger)
	server := config.Server
	accessLog := middlewares.ProvideAccessLog(server, logLogger)
	requestId := middlewares.ProvideRequestId(server)
	static := middlewares.ProvideStatic(server)
	service := sessions.ProvideService()
	withSession := middlewares.ProvideWithSession(service)
	eventsCors := middlewares.ProvideEventsCors()
	favicon := middlewares.ProvideFavicon()
	getSignUp := handlers.ProvideGetSignUp()
	configPostgres := config.Postgres
	pg := postgres.ProvidePg(logLogger, configPostgres)
	usersService := users.ProvideService(pg)
	postSignUp := handlers.ProvidePostSignUp(usersService, service)
	getSignIn := handlers.ProvideGetSignIn()
	authService := auth.ProvideService(usersService)
	postSignIn := handlers.ProvidePostSignIn(authService, service)
	getIndex := handlers.ProvideGetIndex()
	notFound := handlers.ProvideNotFound()
	configClickhouse := config.Clickhouse
	ch := clickhouse.ProvideCh(logLogger, configClickhouse)
	eventstoreService := eventstore.ProvideClickhouseService(ch, logLogger)
	sourceregistryService := sourceregistry.ProvideEnvVarService(logLogger)
	uaparserService := uaparser.ProvideService()
	postPageViewEvent := handlers.ProvidePostEventsPageViews(eventstoreService, sourceregistryService, uaparserService)
	app := ProvideFiber(config, views, middlewaresLogger, accessLog, requestId, static, withSession, eventsCors, favicon, getSignUp, postSignUp, getSignIn, postSignIn, getIndex, notFound, postPageViewEvent)
	mainApp := ProvideApp(config, app, logLogger)
	return mainApp
}
