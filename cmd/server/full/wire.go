//go:build wireinject
// +build wireinject

package full

import (
	"github.com/google/wire"
	"github.com/prismelabs/analytics/pkg/clickhouse"
	grafanaCli "github.com/prismelabs/analytics/pkg/grafana"
	"github.com/prismelabs/analytics/pkg/handlers"
	"github.com/prismelabs/analytics/pkg/middlewares"
	"github.com/prismelabs/analytics/pkg/services/eventstore"
	"github.com/prismelabs/analytics/pkg/services/grafana"
	"github.com/prismelabs/analytics/pkg/services/ipgeolocator"
	"github.com/prismelabs/analytics/pkg/services/sourceregistry"
	"github.com/prismelabs/analytics/pkg/services/uaparser"
	"github.com/prismelabs/analytics/pkg/wired"
)

func Initialize(logger wired.BootstrapLogger) wired.App {
	wire.Build(
		ProvideFiber,
		ProvideSetup,
		clickhouse.ProvideCh,
		eventstore.ProvideClickhouseService,
		grafana.ProvideService,
		grafanaCli.ProvideClient,
		handlers.ProvideHealthCheck,
		handlers.ProvidePostEventsPageViews,
		ipgeolocator.ProvideMmdbService,
		middlewares.ProvideAccessLog,
		middlewares.ProvideErrorHandler,
		middlewares.ProvideEventsCors,
		middlewares.ProvideEventsRateLimiter,
		middlewares.ProvideLogger,
		middlewares.ProvideRequestId,
		middlewares.ProvideStatic,
		sourceregistry.ProvideEnvVarService,
		uaparser.ProvideService,
		wired.ProvideApp,
		wired.ProvideClickhouseConfig,
		wired.ProvideGrafanaConfig,
		wired.ProvideLogger,
		wired.ProvideMinimalFiber,
		wired.ProvideServerConfig,
	)
	return wired.App{}
}
