package ingestion

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prismelabs/prismeanalytics/pkg/handlers"
	"github.com/prismelabs/prismeanalytics/pkg/middlewares"
	"github.com/prismelabs/prismeanalytics/pkg/wired"
)

// ProvideFiber is a wire provider for fiber.App.
func ProvideFiber(
	eventsCorsMiddleware middlewares.EventsCors,
	eventsRateLimiterMiddleware middlewares.EventsRateLimiter,
	minimalFiber wired.MinimalFiber,
	postPageViewEventHandler handlers.PostPageViewEvent,
) *fiber.App {
	app := (*fiber.App)(minimalFiber)

	app.Use("/api/v1/events/*",
		fiber.Handler(eventsCorsMiddleware),
		fiber.Handler(eventsRateLimiterMiddleware),
	)
	app.Post("/api/v1/events/pageviews", fiber.Handler(postPageViewEventHandler))

	return app
}
