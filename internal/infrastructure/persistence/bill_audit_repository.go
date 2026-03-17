package persistence

import (
	"context"

	"github.com/google/uuid"
	"github.com/mvr-garcia/bill-manager/internal/domain"
	"gorm.io/gorm"
)

// GormBillAuditRepository implements the BillAuditRepository interface using GORM.
type GormBillAuditRepository struct {
	db *gorm.DB
}

// NewGormBillAuditRepository creates a new instance of GormBillAuditRepository.
func NewGormBillAuditRepository(db *gorm.DB) *GormBillAuditRepository {
	return &GormBillAuditRepository{db: db}
}

// Create persists a new audit log entry in the database.
func (r *GormBillAuditRepository) Create(ctx context.Context, audit *domain.BillAudit) error {
	model := BillAuditDomainToModel(audit)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	return nil
}

// GetByBillID retrieves all audit logs for a specific bill from the database.
func (r *GormBillAuditRepository) GetByBillID(ctx context.Context, billID uuid.UUID) ([]*domain.BillAudit, error) {
	var models []BillAuditModel
	if err := r.db.WithContext(ctx).Where("bill_id = ?", billID).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}

	audits := make([]*domain.BillAudit, len(models))
	for i, model := range models {
		audits[i] = BillAuditModelToDomain(&model)
	}
	return audits, nil
}

// GetAll retrieves all audit logs from the database.
func (r *GormBillAuditRepository) GetAll(ctx context.Context) ([]*domain.BillAudit, error) {
	var models []BillAuditModel
	if err := r.db.WithContext(ctx).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}

	audits := make([]*domain.BillAudit, len(models))
	for i, model := range models {
		audits[i] = BillAuditModelToDomain(&model)
	}
	return audits, nil
}
