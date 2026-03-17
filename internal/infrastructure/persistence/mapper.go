package persistence

import (
	"github.com/mvr-garcia/bill-manager/internal/domain"
)

// BillModelToDomain converts a BillModel to a domain Bill entity.
func BillModelToDomain(model *BillModel) *domain.Bill {
	if model == nil {
		return nil
	}

	return &domain.Bill{
		ID:          model.ID,
		Description: model.Description,
		Amount:      model.Amount,
		DueDate:     model.DueDate,
		Status:      domain.BillStatus(model.Status),
		CreatedBy:   model.CreatedBy,
		ApprovedBy:  model.ApprovedBy,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}

// BillDomainToModel converts a domain Bill entity to a BillModel.
func BillDomainToModel(bill *domain.Bill) *BillModel {
	if bill == nil {
		return nil
	}

	return &BillModel{
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

// BillAuditModelToDomain converts a BillAuditModel to a domain BillAudit entity.
func BillAuditModelToDomain(model *BillAuditModel) *domain.BillAudit {
	if model == nil {
		return nil
	}

	return &domain.BillAudit{
		ID:          model.ID,
		BillID:      model.BillID,
		Action:      domain.AuditAction(model.Action),
		PerformedBy: model.PerformedBy,
		IPAddress:   model.IPAddress,
		UserAgent:   model.UserAgent,
		CreatedAt:   model.CreatedAt,
	}
}

// BillAuditDomainToModel converts a domain BillAudit entity to a BillAuditModel.
func BillAuditDomainToModel(audit *domain.BillAudit) *BillAuditModel {
	if audit == nil {
		return nil
	}

	return &BillAuditModel{
		ID:          audit.ID,
		BillID:      audit.BillID,
		Action:      string(audit.Action),
		PerformedBy: audit.PerformedBy,
		IPAddress:   audit.IPAddress,
		UserAgent:   audit.UserAgent,
		CreatedAt:   audit.CreatedAt,
	}
}
