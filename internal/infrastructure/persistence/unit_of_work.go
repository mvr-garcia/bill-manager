package persistence

import (
	"github.com/mvr-garcia/bill-manager/internal/domain"
	"gorm.io/gorm"
)

// gormUnitOfWork implements the UnitOfWork interface using GORM.
type gormUnitOfWork struct {
	db *gorm.DB
}

// NewGormUnitOfWork creates a new GORM-based UnitOfWork.
func NewGormUnitOfWork(db *gorm.DB) domain.UnitOfWork {
	return &gormUnitOfWork{db: db}
}

// Execute encapsula a lógica de transação do GORM.
func (u *gormUnitOfWork) Execute(fn func(domain.UnitOfWork) error) error {
	return u.db.Transaction(func(tx *gorm.DB) error {
		// Criamos uma nova instância de UoW que utiliza a transação 'tx'
		return fn(&gormUnitOfWork{db: tx})
	})
}

func (u *gormUnitOfWork) BillRepository() domain.BillRepository {
	return NewGormBillRepository(u.db)
}

func (u *gormUnitOfWork) BillAuditRepository() domain.BillAuditRepository {
	return NewGormBillAuditRepository(u.db)
}
