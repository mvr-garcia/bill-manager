package persistence

import (
	"context"

	"github.com/mvr-garcia/bill-manager/internal/domain"
	"gorm.io/gorm"
)

// gormUnitOfWork implements the UnitOfWork interface using GORM.
type gormUnitOfWork struct {
	db            *gorm.DB
	tx            *gorm.DB
	billRepo      domain.BillRepository
	billAuditRepo domain.BillAuditRepository
}

// NewGormUnitOfWork creates a new GORM-based UnitOfWork.
func NewGormUnitOfWork(db *gorm.DB) domain.UnitOfWork {
	return &gormUnitOfWork{
		db: db,
	}
}

// Begin starts a new transaction.
func (u *gormUnitOfWork) Begin(ctx context.Context) error {
	u.tx = u.db.Begin()
	return u.tx.Error
}

// Commit commits the current transaction.
func (u *gormUnitOfWork) Commit(ctx context.Context) error {
	if u.tx == nil {
		return nil
	}
	err := u.tx.Commit().Error
	u.tx = nil
	return err
}

// Rollback rolls back the current transaction.
func (u *gormUnitOfWork) Rollback(ctx context.Context) error {
	if u.tx == nil {
		return nil
	}
	err := u.tx.Rollback().Error
	u.tx = nil
	return err
}

// BillRepository returns the bill repository for the current transaction.
func (u *gormUnitOfWork) BillRepository() domain.BillRepository {
	if u.tx != nil {
		if u.billRepo == nil {
			u.billRepo = NewGormBillRepository(u.tx)
		}
		return u.billRepo
	}
	return NewGormBillRepository(u.db)
}

// BillAuditRepository returns the bill audit repository for the current transaction.
func (u *gormUnitOfWork) BillAuditRepository() domain.BillAuditRepository {
	if u.tx != nil {
		if u.billAuditRepo == nil {
			u.billAuditRepo = NewGormBillAuditRepository(u.tx)
		}
		return u.billAuditRepo
	}
	return NewGormBillAuditRepository(u.db)
}
