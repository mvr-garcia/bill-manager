package api

import (
	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/mvr-garcia/bill-manager/internal/config"
	"github.com/mvr-garcia/bill-manager/internal/presentation/api/handlers"
	"go.uber.org/fx"
)

func RunHTTPServer(lf fx.Lifecycle, cfg *config.Config, billHandler *handlers.BillHandler) {
	app := fiber.New()

	setupRoutes(app, billHandler)
	lf.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			err := app.Listen(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port))
			if err != nil {
				log.Fatalf("Failed to start server: %v", err)
			}
			return err
		},
		OnStop: func(ctx context.Context) error {
			return app.Shutdown()
		},
	})
}
