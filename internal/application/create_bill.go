package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mvr-garcia/bill-manager/internal/domain"
)

// CreateBillRequest represents the input for creating a bill.
type CreateBillRequest struct {
	Description string
	Amount      float64
	DueDate     time.Time
}

// CreateBillResponse represents the output of creating a bill.
type CreateBillResponse struct {
	ID          uuid.UUID
	Description string
	Amount      float64
	DueDate     time.Time
	Status      string
	CreatedBy   uuid.UUID
	CreatedAt   time.Time
}

// CreateBillUseCase handles the creation of a new bill with audit logging.
type CreateBillUseCase struct {
	uow domain.UnitOfWork
}

// NewCreateBillUseCase creates a new instance of CreateBillUseCase.
func NewCreateBillUseCase(uow domain.UnitOfWork) *CreateBillUseCase {
	return &CreateBillUseCase{uow: uow}
}

// Execute creates a new bill and logs the creation in the audit table.
func (uc *CreateBillUseCase) Execute(ctx context.Context, req *CreateBillRequest, userID uuid.UUID, ipAddress, userAgent string) (*CreateBillResponse, error) {
	// Start a transaction
	if err := uc.uow.Begin(ctx); err != nil {
		return nil, err
	}
	defer func() {
		if r := recover(); r != nil {
			_ = uc.uow.Rollback(ctx)
			panic(r)
		}
	}()

	// Create the bill entity
	bill := domain.NewBill(req.Description, req.Amount, req.DueDate, userID)

	// Persist the bill
	billRepo := uc.uow.BillRepository()
	if err := billRepo.Create(ctx, bill); err != nil {
		_ = uc.uow.Rollback(ctx)
		return nil, err
	}

	// Create and persist the audit log
	audit := domain.NewBillAudit(bill.ID, domain.AuditActionCreated, userID, ipAddress, userAgent)
	auditRepo := uc.uow.BillAuditRepository()
	if err := auditRepo.Create(ctx, audit); err != nil {
		_ = uc.uow.Rollback(ctx)
		return nil, err
	}

	// Commit the transaction
	if err := uc.uow.Commit(ctx); err != nil {
		return nil, err
	}

	return &CreateBillResponse{
		ID:          bill.ID,
		Description: bill.Description,
		Amount:      bill.Amount,
		DueDate:     bill.DueDate,
		Status:      string(bill.Status),
		CreatedBy:   bill.CreatedBy,
		CreatedAt:   bill.CreatedAt,
	}, nil
}
