package domain

import (
	"time"

	"github.com/google/uuid"
)

// BillStatus represents the status of a bill.
type BillStatus string

const (
	// BillStatusPending indicates the bill is pending approval.
	BillStatusPending BillStatus = "pending"
	// BillStatusApproved indicates the bill has been approved.
	BillStatusApproved BillStatus = "approved"
	// BillStatusRejected indicates the bill has been rejected.
	BillStatusRejected BillStatus = "rejected"
	// BillStatusPaid indicates the bill has been paid.
	BillStatusPaid BillStatus = "paid"
)

// Bill represents a payable bill in the system.
type Bill struct {
	ID        uuid.UUID
	Description string
	Amount    float64
	DueDate   time.Time
	Status    BillStatus
	CreatedBy uuid.UUID
	ApprovedBy *uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewBill creates a new Bill entity.
func NewBill(description string, amount float64, dueDate time.Time, createdBy uuid.UUID) *Bill {
	return &Bill{
		ID:        uuid.New(),
		Description: description,
		Amount:    amount,
		DueDate:   dueDate,
		Status:    BillStatusPending,
		CreatedBy: createdBy,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Approve sets the bill as approved by a specific user.
func (b *Bill) Approve(approvedBy uuid.UUID) {
	b.Status = BillStatusApproved
	b.ApprovedBy = &approvedBy
	b.UpdatedAt = time.Now()
}

// Reject sets the bill as rejected.
func (b *Bill) Reject() {
	b.Status = BillStatusRejected
	b.UpdatedAt = time.Now()
}

// MarkAsPaid sets the bill as paid.
func (b *Bill) MarkAsPaid() {
	b.Status = BillStatusPaid
	b.UpdatedAt = time.Now()
}
