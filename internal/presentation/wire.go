package presentation

import (
	"github.com/mvr-garcia/bill-manager/internal/application"
	"github.com/mvr-garcia/bill-manager/internal/config"
	"github.com/mvr-garcia/bill-manager/internal/domain"
	"github.com/mvr-garcia/bill-manager/internal/infrastructure/persistence"
	"github.com/mvr-garcia/bill-manager/internal/presentation/api"
	"github.com/mvr-garcia/bill-manager/internal/presentation/api/handlers"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// Module provides all infrastructure dependencies.
var Module = fx.Module(
	"infrastructure",
	fx.Provide(
		// Configuration
		config.LoadConfig,
		func(cfg *config.Config) config.DatabaseConfig {
			return cfg.Database
		},
		persistence.NewDatabase,

		// Persistence
		func(db *gorm.DB) domain.UnitOfWork {
			return persistence.NewGormUnitOfWork(db)
		},
	),
)

// ApplicationModule provides all application use cases.
var ApplicationModule = fx.Module(
	"application",
	fx.Provide(
		func(uow domain.UnitOfWork) *application.CreateBillUseCase {
			return application.NewCreateBillUseCase(uow)
		},
		func(uow domain.UnitOfWork) *application.ApproveBillUseCase {
			return application.NewApproveBillUseCase(uow)
		},
		func(uow domain.UnitOfWork) *application.GetBillUseCase {
			return application.NewGetBillUseCase(uow)
		},
		func(uow domain.UnitOfWork) *application.ListBillsUseCase {
			return application.NewListBillsUseCase(uow)
		},
		func(uow domain.UnitOfWork) *application.GetBillAuditsUseCase {
			return application.NewGetBillAuditsUseCase(uow)
		},
	),
)

// PresentationModule provides all presentation handlers.
var PresentationModule = fx.Module(
	"presentation",
	fx.Provide(
		func(
			createBillUseCase *application.CreateBillUseCase,
			approveBillUseCase *application.ApproveBillUseCase,
			getBillUseCase *application.GetBillUseCase,
			listBillsUseCase *application.ListBillsUseCase,
			getBillAuditsUseCase *application.GetBillAuditsUseCase,
		) *handlers.BillHandler {
			return handlers.NewBillHandler(
				createBillUseCase,
				approveBillUseCase,
				getBillUseCase,
				listBillsUseCase,
				getBillAuditsUseCase,
			)
		},
	),
)

// ServerModule provides the Fiber server.
var ServerModule = fx.Module(
	"server",
	fx.Invoke(
		func(lf fx.Lifecycle, cfg *config.Config, billHandler *handlers.BillHandler) {
			api.RunHTTPServer(lf, cfg, billHandler)
		},
	),
)
