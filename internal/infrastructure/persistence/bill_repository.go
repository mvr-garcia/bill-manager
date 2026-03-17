package persistence

import (
	"context"

	"github.com/google/uuid"
	"github.com/mvr-garcia/bill-manager/internal/domain"
	"gorm.io/gorm"
)

// GormBillRepository implements the BillRepository interface using GORM.
type GormBillRepository struct {
	db *gorm.DB
}

// NewGormBillRepository creates a new instance of GormBillRepository.
func NewGormBillRepository(db *gorm.DB) *GormBillRepository {
	return &GormBillRepository{db: db}
}

// Create persists a new bill in the database.
func (r *GormBillRepository) Create(ctx context.Context, bill *domain.Bill) error {
	model := BillDomainToModel(bill)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	return nil
}

// GetByID retrieves a bill by its ID from the database.
func (r *GormBillRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Bill, error) {
	var model BillModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return BillModelToDomain(&model), nil
}

// GetAll retrieves all bills from the database, optionally filtered by status.
func (r *GormBillRepository) GetAll(ctx context.Context, status *domain.BillStatus) ([]*domain.Bill, error) {
	var models []BillModel
	query := r.db.WithContext(ctx)

	if status != nil {
		query = query.Where("status = ?", string(*status))
	}

	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	bills := make([]*domain.Bill, len(models))
	for i, model := range models {
		bills[i] = BillModelToDomain(&model)
	}
	return bills, nil
}

// Update updates an existing bill in the database.
func (r *GormBillRepository) Update(ctx context.Context, bill *domain.Bill) error {
	model := BillDomainToModel(bill)
	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return err
	}
	return nil
}

// Delete removes a bill from the database.
func (r *GormBillRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&BillModel{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
