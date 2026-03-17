package domain

import (
	"time"

	"github.com/google/uuid"
)

// AuditAction represents the type of action performed on a bill.
type AuditAction string

const (
	// AuditActionCreated indicates a bill was created.
	AuditActionCreated AuditAction = "created"
	// AuditActionApproved indicates a bill was approved.
	AuditActionApproved AuditAction = "approved"
	// AuditActionRejected indicates a bill was rejected.
	AuditActionRejected AuditAction = "rejected"
	// AuditActionPaid indicates a bill was marked as paid.
	AuditActionPaid AuditAction = "paid"
)

// BillAudit represents an audit log entry for a bill operation.
type BillAudit struct {
	ID        uuid.UUID
	BillID    uuid.UUID
	Action    AuditAction
	PerformedBy uuid.UUID
	IPAddress string
	UserAgent string
	CreatedAt time.Time
}

// NewBillAudit creates a new BillAudit entity.
func NewBillAudit(billID uuid.UUID, action AuditAction, performedBy uuid.UUID, ipAddress, userAgent string) *BillAudit {
	return &BillAudit{
		ID:         uuid.New(),
		BillID:     billID,
		Action:     action,
		PerformedBy: performedBy,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		CreatedAt: time.Now(),
	}
}
