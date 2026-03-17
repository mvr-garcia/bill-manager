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

	// Retrieve the bill
	billRepo := uc.uow.BillRepository()
	bill, err := billRepo.GetByID(ctx, req.BillID)
	if err != nil {
		_ = uc.uow.Rollback(ctx)
		return nil, err
	}

	if bill == nil {
		_ = uc.uow.Rollback(ctx)
		return nil, errors.New("bill not found")
	}

	// Check if the bill is already approved
	if bill.Status != domain.BillStatusPending {
		_ = uc.uow.Rollback(ctx)
		return nil, errors.New("only pending bills can be approved")
	}

	// Approve the bill
	bill.Approve(approverID)

	// Update the bill
	if err := billRepo.Update(ctx, bill); err != nil {
		_ = uc.uow.Rollback(ctx)
		return nil, err
	}

	// Create and persist the audit log
	audit := domain.NewBillAudit(bill.ID, domain.AuditActionApproved, approverID, ipAddress, userAgent)
	auditRepo := uc.uow.BillAuditRepository()
	if err := auditRepo.Create(ctx, audit); err != nil {
		_ = uc.uow.Rollback(ctx)
		return nil, err
	}

	// Commit the transaction
	if err := uc.uow.Commit(ctx); err != nil {
		return nil, err
	}

	return &ApproveBillResponse{
		ID:         bill.ID,
		Status:     string(bill.Status),
		ApprovedBy: approverID,
		UpdatedAt:  bill.UpdatedAt,
	}, nil
}
