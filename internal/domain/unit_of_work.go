package domain

import "context"

// UnitOfWork defines the interface for managing transactions across multiple repositories.
// It ensures that operations on multiple repositories are executed atomically.
type UnitOfWork interface {
	// Begin starts a new transaction.
	Begin(ctx context.Context) error
	// Commit commits the current transaction.
	Commit(ctx context.Context) error
	// Rollback rolls back the current transaction.
	Rollback(ctx context.Context) error
	// BillRepository returns the bill repository for the current transaction.
	BillRepository() BillRepository
	// BillAuditRepository returns the bill audit repository for the current transaction.
	BillAuditRepository() BillAuditRepository
}
