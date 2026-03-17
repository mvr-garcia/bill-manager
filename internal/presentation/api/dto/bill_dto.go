package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateBillRequest represents the DTO for creating a bill.
type CreateBillRequest struct {
	Description string    `json:"description" validate:"required,min=1,max=500"`
	Amount      float64   `json:"amount" validate:"required,gt=0"`
	DueDate     time.Time `json:"due_date" validate:"required"`
}

// CreateBillResponse represents the DTO for the create bill response.
type CreateBillResponse struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	DueDate     time.Time `json:"due_date"`
	Status      string    `json:"status"`
	CreatedBy   uuid.UUID `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

// ApproveBillRequest represents the DTO for approving a bill.
type ApproveBillRequest struct {
	BillID uuid.UUID `json:"bill_id" validate:"required"`
}

// ApproveBillResponse represents the DTO for the approve bill response.
type ApproveBillResponse struct {
	ID         uuid.UUID `json:"id"`
	Status     string    `json:"status"`
	ApprovedBy uuid.UUID `json:"approved_by"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// GetBillResponse represents the DTO for retrieving a bill.
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

// ListBillsResponse represents the DTO for listing bills.
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

// BillAuditResponse represents the DTO for a bill audit log entry.
type BillAuditResponse struct {
	ID         uuid.UUID `json:"id"`
	BillID     uuid.UUID `json:"bill_id"`
	Action     string    `json:"action"`
	PerformedBy uuid.UUID `json:"performed_by"`
	IPAddress  string    `json:"ip_address"`
	UserAgent  string    `json:"user_agent"`
	CreatedAt  time.Time `json:"created_at"`
}
