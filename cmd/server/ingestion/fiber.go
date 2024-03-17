package ingestion

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prismelabs/analytics/pkg/handlers"
	"github.com/prismelabs/analytics/pkg/middlewares"
	"github.com/prismelabs/analytics/pkg/wired"
)

// ProvideFiber is a wire provider for fiber.App.
func ProvideFiber(
	eventsCorsMiddleware middlewares.EventsCors,
	eventsRateLimiterMiddleware middlewares.EventsRateLimiter,
	minimalFiber wired.MinimalFiber,
	postPageViewEventHandler handlers.PostEventsPageview,
	postCustomEventHandler handlers.PostEventsCustom,
) *fiber.App {
	app := (*fiber.App)(minimalFiber)

	// Public endpoints.
	app.Use("/api/v1/events/*",
		fiber.Handler(eventsCorsMiddleware),
		fiber.Handler(eventsRateLimiterMiddleware),
	)
	app.Post("/api/v1/events/pageviews", fiber.Handler(postPageViewEventHandler))
	app.Post("/api/v1/events/custom/:name", fiber.Handler(postCustomEventHandler))

	return app
}
