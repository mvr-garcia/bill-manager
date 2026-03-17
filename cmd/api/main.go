package main

import (
	"context"
	"log"

	"github.com/mvr-garcia/bill-manager/internal/config"
	"github.com/mvr-garcia/bill-manager/internal/presentation"
	"go.uber.org/fx"
)

func main() {
	// Create the FX app with all modules
	app := fx.New(
		fx.Provide(config.LoadConfig),
		presentation.Module,
		presentation.ApplicationModule,
		presentation.PresentationModule,
		presentation.ServerModule,
	)

	// Start the application
	if err := app.Start(context.Background()); err != nil {
		log.Fatalf("Application error: %v\n", err)
	}
}
