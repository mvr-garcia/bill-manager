package application

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/mvr-garcia/bill-manager/internal/domain"
)

// ApproveBillRequest represents the input for approving a bill.
type ApproveBillRequest struct {
	BillID uuid.UUID
}

// ApproveBillResponse represents the output of approving a bill.
type ApproveBillResponse struct {
	ID         uuid.UUID
	Status     string
	ApprovedBy uuid.UUID
	UpdatedAt  time.Time
}

// ApproveBillUseCase handles the approval of a bill with audit logging.
type ApproveBillUseCase struct {
	uow domain.UnitOfWork
}

// NewApproveBillUseCase creates a new instance of ApproveBillUseCase.
func NewApproveBillUseCase(uow domain.UnitOfWork) *ApproveBillUseCase {
	return &ApproveBillUseCase{uow: uow}
}

// Execute approves a bill and logs the approval in the audit table.
func (uc *ApproveBillUseCase) Execute(ctx context.Context, req *ApproveBillRequest, approverID uuid.UUID, ipAddress, userAgent string) (*ApproveBillResponse, error) {
	var resp *ApproveBillResponse

	err := uc.uow.Execute(func(uow domain.UnitOfWork) error {
		// Retrieve the bill
		billRepo := uow.BillRepository()
		bill, err := billRepo.GetByID(ctx, req.BillID)
		if err != nil {
			return err
		}

		if bill == nil {
			return errors.New("bill not found")
		}

		// Check if the bill is already approved
		if bill.Status != domain.BillStatusPending {
			return errors.New("only pending bills can be approved")
		}

		// Approve the bill
		bill.Approve(approverID)

		// Update the bill
		if err := billRepo.Update(ctx, bill); err != nil {
			return err
		}

		// Create and persist the audit log
		audit := domain.NewBillAudit(bill.ID, domain.AuditActionApproved, approverID, ipAddress, userAgent)
		auditRepo := uow.BillAuditRepository()
		if err := auditRepo.Create(ctx, audit); err != nil {
			return err
		}

		resp = &ApproveBillResponse{
			ID:         bill.ID,
			Status:     string(bill.Status),
			ApprovedBy: approverID,
			UpdatedAt:  bill.UpdatedAt,
		}
		return nil
	})

	return resp, err
}
