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
	var resp *CreateBillResponse

	err := uc.uow.Execute(func(uow domain.UnitOfWork) error {
		// Create the bill entity
		bill := domain.NewBill(req.Description, req.Amount, req.DueDate, userID)

		// Persist the bill
		billRepo := uow.BillRepository()
		if err := billRepo.Create(ctx, bill); err != nil {
			return err
		}

		// Create and persist the audit log
		audit := domain.NewBillAudit(bill.ID, domain.AuditActionCreated, userID, ipAddress, userAgent)
		auditRepo := uow.BillAuditRepository()
		if err := auditRepo.Create(ctx, audit); err != nil {
			return err
		}

		resp = &CreateBillResponse{
			ID:          bill.ID,
			Description: bill.Description,
			Amount:      bill.Amount,
			DueDate:     bill.DueDate,
			Status:      string(bill.Status),
			CreatedBy:   bill.CreatedBy,
			CreatedAt:   bill.CreatedAt,
		}
		return nil
	})

	return resp, err
}
