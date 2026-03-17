package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mvr-garcia/bill-manager/internal/domain"
)

// ListBillsRequest represents the input for listing bills.
type ListBillsRequest struct {
	Status *string
}

// ListBillsResponse represents a single bill in the list response.
type ListBillsResponse struct {
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

// ListBillsUseCase handles listing bills with optional filtering.
type ListBillsUseCase struct {
	uow domain.UnitOfWork
}

// NewListBillsUseCase creates a new instance of ListBillsUseCase.
func NewListBillsUseCase(uow domain.UnitOfWork) *ListBillsUseCase {
	return &ListBillsUseCase{uow: uow}
}

// Execute lists all bills, optionally filtered by status.
func (uc *ListBillsUseCase) Execute(ctx context.Context, req *ListBillsRequest) ([]*ListBillsResponse, error) {
	var status *domain.BillStatus
	if req.Status != nil {
		s := domain.BillStatus(*req.Status)
		status = &s
	}

	billRepo := uc.uow.BillRepository()
	bills, err := billRepo.GetAll(ctx, status)
	if err != nil {
		return nil, err
	}

	responses := make([]*ListBillsResponse, len(bills))
	for i, bill := range bills {
		responses[i] = &ListBillsResponse{
			ID:          bill.ID,
			Description: bill.Description,
			Amount:      bill.Amount,
			DueDate:     bill.DueDate,
			Status:      string(bill.Status),
			CreatedBy:   bill.CreatedBy,
			ApprovedBy:  bill.ApprovedBy,
			CreatedAt:   bill.CreatedAt,
			UpdatedAt:   bill.UpdatedAt,
		}
	}

	return responses, nil
}
