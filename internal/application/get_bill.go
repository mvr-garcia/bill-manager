package application

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/mvr-garcia/bill-manager/internal/domain"
)

// GetBillResponse represents the output of retrieving a bill.
type GetBillResponse struct {
	ID          uuid.UUID  `json:"id"`
	Description string     `json:"description"`
	Amount      float64    `json:"amount"`
	DueDate     time.Time  `json:"due_date"`
	Status      string     `json:"status"`
	CreatedBy   uuid.UUID  `json:"created_by"`
	ApprovedBy  *uuid.UUID `json:"approved_by,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// GetBillUseCase handles retrieving a bill by ID.
type GetBillUseCase struct {
	uow domain.UnitOfWork
}

// NewGetBillUseCase creates a new instance of GetBillUseCase.
func NewGetBillUseCase(uow domain.UnitOfWork) *GetBillUseCase {
	return &GetBillUseCase{uow: uow}
}

// Execute retrieves a bill by ID.
func (uc *GetBillUseCase) Execute(ctx context.Context, billID uuid.UUID) (*GetBillResponse, error) {
	billRepo := uc.uow.BillRepository()
	bill, err := billRepo.GetByID(ctx, billID)
	if err != nil {
		return nil, err
	}

	if bill == nil {
		return nil, errors.New("bill not found")
	}

	return &GetBillResponse{
		ID:          bill.ID,
		Description: bill.Description,
		Amount:      bill.Amount,
		DueDate:     bill.DueDate,
		Status:      string(bill.Status),
		CreatedBy:   bill.CreatedBy,
		ApprovedBy:  bill.ApprovedBy,
		CreatedAt:   bill.CreatedAt,
		UpdatedAt:   bill.UpdatedAt,
	}, nil
}
