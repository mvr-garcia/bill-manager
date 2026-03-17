package domain

// UnitOfWork defines the interface for managing transactions across multiple repositories.
// It ensures that operations on multiple repositories are executed atomically.
type UnitOfWork interface {
	// Execute encapsulates transactional logic.
	Execute(fn func(UnitOfWork) error) error
	// BillRepository returns the bill repository for the current transaction.
	BillRepository() BillRepository
	// BillAuditRepository returns the bill audit repository for the current transaction.
	BillAuditRepository() BillAuditRepository
}
