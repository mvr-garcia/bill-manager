package domain

import (
	"context"

	"github.com/google/uuid"
)

// BillRepository defines the interface for bill persistence operations.
type BillRepository interface {
	// Create persists a new bill in the repository.
	Create(ctx context.Context, bill *Bill) error
	// GetByID retrieves a bill by its ID.
	GetByID(ctx context.Context, id uuid.UUID) (*Bill, error)
	// GetAll retrieves all bills with optional filtering.
	GetAll(ctx context.Context, status *BillStatus) ([]*Bill, error)
	// Update updates an existing bill.
	Update(ctx context.Context, bill *Bill) error
	// Delete removes a bill from the repository.
	Delete(ctx context.Context, id uuid.UUID) error
}

// BillAuditRepository defines the interface for bill audit log persistence operations.
type BillAuditRepository interface {
	// Create persists a new audit log entry.
	Create(ctx context.Context, audit *BillAudit) error
	// GetByBillID retrieves all audit logs for a specific bill.
	GetByBillID(ctx context.Context, billID uuid.UUID) ([]*BillAudit, error)
	// GetAll retrieves all audit logs.
	GetAll(ctx context.Context) ([]*BillAudit, error)
}
