package api

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mvr-garcia/bill-manager/internal/presentation/api/handlers"
	"github.com/mvr-garcia/bill-manager/internal/presentation/api/middleware"
)

// SetupRoutes configures all API routes.
func setupRoutes(app *fiber.App, billHandler *handlers.BillHandler) {
	// Health check endpoint (no auth required)
	app.Get("/healthz", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	// API v1 routes with JWT middleware
	api := app.Group("/api/v1", middleware.JWTMiddleware())

	// Bill routes
	bills := api.Group("/bills")
	bills.Post("", billHandler.CreateBill)
	bills.Get("", billHandler.ListBills)
	bills.Get("/:id", billHandler.GetBill)
	bills.Post("/:id/approve", billHandler.ApproveBill)
	bills.Get("/:id/audits", billHandler.GetBillAudits)
}
