//go:build wireinject
// +build wireinject

package ingestion

import (
	"github.com/google/wire"
	"github.com/prismelabs/analytics/pkg/clickhouse"
	"github.com/prismelabs/analytics/pkg/handlers"
	"github.com/prismelabs/analytics/pkg/middlewares"
	"github.com/prismelabs/analytics/pkg/services/eventstore"
	"github.com/prismelabs/analytics/pkg/services/ipgeolocator"
	"github.com/prismelabs/analytics/pkg/services/originregistry"
	"github.com/prismelabs/analytics/pkg/services/saltmanager"
	"github.com/prismelabs/analytics/pkg/services/teardown"
	"github.com/prismelabs/analytics/pkg/services/uaparser"
	"github.com/prismelabs/analytics/pkg/wired"
)

func Initialize(logger wired.BootstrapLogger) wired.App {
	wire.Build(
		ProvideFiber,
		clickhouse.ProvideCh,
		clickhouse.ProvideEmbeddedSourceDriver,
		eventstore.ProvideClickhouseService,
		handlers.ProvideHealthCheck,
		handlers.ProvidePostEventsCustom,
		handlers.ProvidePostEventsPageViews,
		ipgeolocator.ProvideMmdbService,
		middlewares.ProvideAccessLog,
		middlewares.ProvideErrorHandler,
		middlewares.ProvideEventsCors,
		middlewares.ProvideEventsRateLimiter,
		middlewares.ProvideMetrics,
		middlewares.ProvideNonRegisteredOriginFilter,
		middlewares.ProvideRequestId,
		middlewares.ProvideStatic,
		originregistry.ProvideEnvVarService,
		saltmanager.ProvideService,
		teardown.ProvideService,
		uaparser.ProvideService,
		wired.ProvideApp,
		wired.ProvideClickhouseConfig,
		wired.ProvideLogger,
		wired.ProvideMinimalFiber,
		wired.ProvideMinimalFiberConfig,
		wired.ProvidePromHttpLogger,
		wired.ProvidePrometheusRegistry,
		wired.ProvideServerConfig,
		wired.ProvideSetup,
	)
	return wired.App{}
}
