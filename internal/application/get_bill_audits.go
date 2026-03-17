package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mvr-garcia/bill-manager/internal/domain"
)

// GetBillAuditsResponse represents a single audit log entry.
type GetBillAuditsResponse struct {
	ID          uuid.UUID `json:"id"`
	BillID      uuid.UUID `json:"bill_id"`
	Action      string    `json:"action"`
	PerformedBy uuid.UUID `json:"performed_by"`
	IPAddress   string    `json:"ip_address"`
	UserAgent   string    `json:"user_agent"`
	CreatedAt   time.Time `json:"created_at"`
}

// GetBillAuditsUseCase handles retrieving audit logs for a specific bill.
type GetBillAuditsUseCase struct {
	uow domain.UnitOfWork
}

// NewGetBillAuditsUseCase creates a new instance of GetBillAuditsUseCase.
func NewGetBillAuditsUseCase(uow domain.UnitOfWork) *GetBillAuditsUseCase {
	return &GetBillAuditsUseCase{uow: uow}
}

// Execute retrieves all audit logs for a specific bill.
func (uc *GetBillAuditsUseCase) Execute(ctx context.Context, billID uuid.UUID) ([]*GetBillAuditsResponse, error) {
	auditRepo := uc.uow.BillAuditRepository()
	audits, err := auditRepo.GetByBillID(ctx, billID)
	if err != nil {
		return nil, err
	}

	responses := make([]*GetBillAuditsResponse, len(audits))
	for i, audit := range audits {
		responses[i] = &GetBillAuditsResponse{
			ID:          audit.ID,
			BillID:      audit.BillID,
			Action:      string(audit.Action),
			PerformedBy: audit.PerformedBy,
			IPAddress:   audit.IPAddress,
			UserAgent:   audit.UserAgent,
			CreatedAt:   audit.CreatedAt,
		}
	}

	return responses, nil
}
