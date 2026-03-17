package application

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mvr-garcia/bill-manager/internal/domain"
	"github.com/mvr-garcia/bill-manager/internal/infrastructure/persistence"
	"github.com/mvr-garcia/bill-manager/internal/tests"
	"github.com/stretchr/testify/suite"
)

// ApproveBillTestSuite defines the test suite for ApproveBillUseCase.
type ApproveBillTestSuite struct {
	suite.Suite
	testDB             *tests.TestDatabase
	uow                domain.UnitOfWork
	createBillUseCase  *CreateBillUseCase
	approveBillUseCase *ApproveBillUseCase
}

// SetupSuite initializes the test suite with a PostgreSQL container and runs migrations.
func (s *ApproveBillTestSuite) SetupSuite() {
	ctx := context.Background()

	// Setup test database with migrations
	testDB, err := tests.SetupTestDatabase(ctx)
	s.Assert().NoError(err, "failed to setup test database")
	s.testDB = testDB

	// Initialize Unit of Work and UseCases
	s.uow = persistence.NewGormUnitOfWork(s.testDB.GetDB())
	s.createBillUseCase = NewCreateBillUseCase(s.uow)
	s.approveBillUseCase = NewApproveBillUseCase(s.uow)
}

// TearDownSuite cleans up the test suite by terminating the container.
func (s *ApproveBillTestSuite) TearDownSuite() {
	ctx := context.Background()
	if s.testDB != nil {
		err := s.testDB.Cleanup(ctx)
		s.Assert().NoError(err, "failed to cleanup test database")
	}
}

// SetupTest cleans up the database before each test.
func (s *ApproveBillTestSuite) SetupTest() {
	ctx := context.Background()
	err := s.testDB.CleanupTables(ctx)
	s.Assert().NoError(err, "failed to cleanup tables")
}

// TestApproveBillTestSuite runs the entire test suite.
func TestApproveBillTestSuite(t *testing.T) {
	suite.Run(t, new(ApproveBillTestSuite))
}

// createBillForTesting is a helper function to create a bill for testing approval.
func (s *ApproveBillTestSuite) createBillForTesting(ctx context.Context) *CreateBillResponse {
	userID := uuid.New()
	req := &CreateBillRequest{
		Description: "Test Bill",
		Amount:      1500.50,
		DueDate:     time.Now().AddDate(0, 0, 30),
	}

	result, err := s.createBillUseCase.Execute(ctx, req, userID, "192.168.1.1", "test-agent")
	s.Assert().NoError(err)
	s.Assert().NotNil(result)

	return result
}

// TestApproveBillUseCase_WhenValidRequest_ReturnBillApproved tests successful bill approval.
func (s *ApproveBillTestSuite) TestApproveBillUseCase_WhenValidRequest_ReturnBillApproved() {
	s.T().Parallel()

	ctx := context.Background()
	bill := s.createBillForTesting(ctx)
	approverID := uuid.New()
	ipAddress := "192.168.1.100"
	userAgent := "Mozilla/5.0"

	req := &ApproveBillRequest{
		BillID: bill.ID,
	}

	// Execute
	result, err := s.approveBillUseCase.Execute(ctx, req, approverID, ipAddress, userAgent)

	// Assert
	s.Assert().NoError(err)
	s.Assert().NotNil(result)
	s.Equal(bill.ID, result.ID)
	s.Equal(string(domain.BillStatusApproved), result.Status)
	s.Equal(approverID, result.ApprovedBy)

	// Verify bill was updated in database
	updatedBill, err := s.uow.BillRepository().GetByID(ctx, bill.ID)
	s.Assert().NoError(err)
	s.Assert().NotNil(updatedBill)
	s.Equal(domain.BillStatusApproved, updatedBill.Status)
	s.Equal(approverID, *updatedBill.ApprovedBy)

	// Verify audit log was created
	audits, err := s.uow.BillAuditRepository().GetByBillID(ctx, bill.ID)
	s.Assert().NoError(err)
	s.Assert().Len(audits, 2) // One for creation, one for approval
	s.Equal(domain.AuditActionApproved, audits[0].Action)
	s.Equal(approverID, audits[0].PerformedBy)
	s.Equal(ipAddress, audits[0].IPAddress)
	s.Equal(userAgent, audits[0].UserAgent)
}

// TestApproveBillUseCase_WhenBillNotFound_ReturnError tests approval of non-existent bill.
func (s *ApproveBillTestSuite) TestApproveBillUseCase_WhenBillNotFound_ReturnError() {
	s.T().Parallel()

	ctx := context.Background()
	nonExistentBillID := uuid.New()
	approverID := uuid.New()

	req := &ApproveBillRequest{
		BillID: nonExistentBillID,
	}

	// Execute
	result, err := s.approveBillUseCase.Execute(ctx, req, approverID, "192.168.1.1", "test-agent")

	// Assert
	s.Assert().Error(err)
	s.Assert().Nil(result)
	s.Contains(err.Error(), "bill not found")
}

// TestApproveBillUseCase_WhenBillAlreadyApproved_ReturnError tests approval of already approved bill.
func (s *ApproveBillTestSuite) TestApproveBillUseCase_WhenBillAlreadyApproved_ReturnError() {
	s.T().Parallel()

	ctx := context.Background()
	bill := s.createBillForTesting(ctx)
	approverID1 := uuid.New()
	approverID2 := uuid.New()

	// First approval
	req := &ApproveBillRequest{
		BillID: bill.ID,
	}

	result1, err := s.approveBillUseCase.Execute(ctx, req, approverID1, "192.168.1.100", "test-agent")
	s.Assert().NoError(err)
	s.Assert().NotNil(result1)

	// Second approval attempt
	result2, err := s.approveBillUseCase.Execute(ctx, req, approverID2, "192.168.1.101", "test-agent")

	// Assert
	s.Assert().Error(err)
	s.Assert().Nil(result2)
	s.Contains(err.Error(), "only pending bills can be approved")
}

// TestApproveBillUseCase_WhenMultipleBillsApproved_ReturnAllApproved tests approving multiple bills.
func (s *ApproveBillTestSuite) TestApproveBillUseCase_WhenMultipleBillsApproved_ReturnAllApproved() {
	s.T().Parallel()

	ctx := context.Background()
	approverID := uuid.New()

	// Create and approve 5 bills
	for i := 1; i <= 5; i++ {
		bill := s.createBillForTesting(ctx)
		req := &ApproveBillRequest{
			BillID: bill.ID,
		}

		result, err := s.approveBillUseCase.Execute(ctx, req, approverID, "192.168.1.1", "test-agent")
		s.Assert().NoError(err)
		s.Assert().NotNil(result)
		s.Equal(string(domain.BillStatusApproved), result.Status)
	}

	// Verify all bills are approved
	bills, err := s.uow.BillRepository().GetAll(ctx, nil)
	s.Assert().NoError(err)
	s.Assert().Len(bills, 5)

	for _, b := range bills {
		s.Equal(domain.BillStatusApproved, b.Status)
		s.NotNil(b.ApprovedBy)
		s.Equal(approverID, *b.ApprovedBy)
	}
}

// TestApproveBillUseCase_WhenDifferentApprovers_ReturnCorrectApprovers tests approval by different users.
func (s *ApproveBillTestSuite) TestApproveBillUseCase_WhenDifferentApprovers_ReturnCorrectApprovers() {
	s.T().Parallel()

	ctx := context.Background()
	approver1 := uuid.New()
	approver2 := uuid.New()
	approver3 := uuid.New()

	billIDs := make([]uuid.UUID, 3)
	approvers := []uuid.UUID{approver1, approver2, approver3}

	// Create bills
	for i := 0; i < 3; i++ {
		bill := s.createBillForTesting(ctx)
		billIDs[i] = bill.ID
	}

	// Approve with different users
	for i, billID := range billIDs {
		req := &ApproveBillRequest{
			BillID: billID,
		}

		result, err := s.approveBillUseCase.Execute(ctx, req, approvers[i], fmt.Sprintf("192.168.1.%d", i+1), "test-agent")
		s.Assert().NoError(err)
		s.Assert().Equal(approvers[i], result.ApprovedBy)
	}

	// Verify each bill has correct approver
	for i, billID := range billIDs {
		bill, err := s.uow.BillRepository().GetByID(ctx, billID)
		s.Assert().NoError(err)
		s.Assert().NotNil(bill.ApprovedBy)
		s.Equal(approvers[i], *bill.ApprovedBy)
	}
}

// TestApproveBillUseCase_WhenAuditLogCreated_ReturnAuditWithCorrectMetadata tests audit log metadata.
func (s *ApproveBillTestSuite) TestApproveBillUseCase_WhenAuditLogCreated_ReturnAuditWithCorrectMetadata() {
	s.T().Parallel()

	ctx := context.Background()
	bill := s.createBillForTesting(ctx)
	approverID := uuid.New()
	ipAddress := "203.0.113.100"
	userAgent := "Mozilla/5.0 (X11; Linux x86_64)"

	req := &ApproveBillRequest{
		BillID: bill.ID,
	}

	result, err := s.approveBillUseCase.Execute(ctx, req, approverID, ipAddress, userAgent)
	s.Assert().NoError(err)
	s.Assert().NotNil(result)

	// Verify audit log contains correct metadata
	audits, err := s.uow.BillAuditRepository().GetByBillID(ctx, bill.ID)
	s.Assert().NoError(err)
	s.Assert().Len(audits, 2)

	approvalAudit := audits[0]
	s.Equal(bill.ID, approvalAudit.BillID)
	s.Equal(domain.AuditActionApproved, approvalAudit.Action)
	s.Equal(approverID, approvalAudit.PerformedBy)
	s.Equal(ipAddress, approvalAudit.IPAddress)
	s.Equal(userAgent, approvalAudit.UserAgent)
	s.NotZero(approvalAudit.CreatedAt)
}

// TestApproveBillUseCase_WhenBillApproved_ReturnUpdatedTimestamp tests that updated_at is updated.
func (s *ApproveBillTestSuite) TestApproveBillUseCase_WhenBillApproved_ReturnUpdatedTimestamp() {
	s.T().Parallel()

	ctx := context.Background()
	bill := s.createBillForTesting(ctx)
	approverID := uuid.New()

	// Get original timestamps
	originalBill, err := s.uow.BillRepository().GetByID(ctx, bill.ID)
	s.Assert().NoError(err)
	originalUpdatedAt := originalBill.UpdatedAt

	// Wait a bit to ensure timestamp difference
	time.Sleep(100 * time.Millisecond)

	req := &ApproveBillRequest{
		BillID: bill.ID,
	}

	result, err := s.approveBillUseCase.Execute(ctx, req, approverID, "192.168.1.1", "test-agent")
	s.Assert().NoError(err)

	// Verify updated_at changed
	s.True(result.UpdatedAt.After(originalUpdatedAt))

	// Verify in database
	updatedBill, err := s.uow.BillRepository().GetByID(ctx, bill.ID)
	s.Assert().NoError(err)
	s.True(updatedBill.UpdatedAt.After(originalUpdatedAt))
}

// TestApproveBillUseCase_WhenBillApproved_ReturnCreatedByUnchanged tests that created_by is not changed.
func (s *ApproveBillTestSuite) TestApproveBillUseCase_WhenBillApproved_ReturnCreatedByUnchanged() {
	s.T().Parallel()

	ctx := context.Background()
	bill := s.createBillForTesting(ctx)
	originalCreatedBy := bill.CreatedBy
	approverID := uuid.New()

	req := &ApproveBillRequest{
		BillID: bill.ID,
	}

	result, err := s.approveBillUseCase.Execute(ctx, req, approverID, "192.168.1.1", "test-agent")
	s.Assert().NoError(err)
	s.Assert().NotNil(result)

	// Verify created_by is unchanged
	updatedBill, err := s.uow.BillRepository().GetByID(ctx, bill.ID)
	s.Assert().NoError(err)
	s.Equal(originalCreatedBy, updatedBill.CreatedBy)
}

// TestApproveBillUseCase_WhenConcurrentApprovals_ReturnOnlyFirstApproved tests concurrent approval attempts.
func (s *ApproveBillTestSuite) TestApproveBillUseCase_WhenConcurrentApprovals_ReturnOnlyFirstApproved() {
	s.T().Parallel()

	ctx := context.Background()
	bill := s.createBillForTesting(ctx)
	numApprovers := 5
	results := make(chan *ApproveBillResponse, numApprovers)
	errors := make(chan error, numApprovers)

	// Try to approve the same bill concurrently
	for i := 0; i < numApprovers; i++ {
		go func(index int) {
			approverID := uuid.New()
			req := &ApproveBillRequest{
				BillID: bill.ID,
			}

			result, err := s.approveBillUseCase.Execute(ctx, req, approverID, "192.168.1.1", "test-agent")
			if err != nil {
				errors <- err
				return
			}
			results <- result
		}(i)
	}

	// Collect results
	successCount := 0
	errorCount := 0
	for i := 0; i < numApprovers; i++ {
		select {
		case err := <-errors:
			errorCount++
			s.T().Logf("Expected error: %v", err)
		case result := <-results:
			s.Assert().NotNil(result)
			successCount++
		}
	}

	// Only one should succeed, others should fail
	s.Equal(1, successCount, "only one approval should succeed")
	s.Equal(numApprovers-1, errorCount, "other approvals should fail")

	// Verify only one audit log for approval
	audits, err := s.uow.BillAuditRepository().GetByBillID(ctx, bill.ID)
	s.Assert().NoError(err)
	approvalCount := 0
	for _, audit := range audits {
		if audit.Action == domain.AuditActionApproved {
			approvalCount++
		}
	}
	s.Equal(1, approvalCount, "only one approval audit should exist")
}

// TestApproveBillUseCase_WhenBillApprovedTwice_ReturnErrorOnSecondApproval tests idempotency.
func (s *ApproveBillTestSuite) TestApproveBillUseCase_WhenBillApprovedTwice_ReturnErrorOnSecondApproval() {
	s.T().Parallel()

	ctx := context.Background()
	bill := s.createBillForTesting(ctx)
	approverID := uuid.New()

	req := &ApproveBillRequest{
		BillID: bill.ID,
	}

	// First approval
	result1, err := s.approveBillUseCase.Execute(ctx, req, approverID, "192.168.1.1", "test-agent")
	s.Assert().NoError(err)
	s.Assert().NotNil(result1)

	// Second approval with same bill
	result2, err := s.approveBillUseCase.Execute(ctx, req, approverID, "192.168.1.1", "test-agent")
	s.Assert().Error(err)
	s.Assert().Nil(result2)
	s.Contains(err.Error(), "only pending bills can be approved")
}

// TestApproveBillUseCase_WhenBillApproved_ReturnCorrectAuditOrder tests audit log ordering.
func (s *ApproveBillTestSuite) TestApproveBillUseCase_WhenBillApproved_ReturnCorrectAuditOrder() {
	s.T().Parallel()

	ctx := context.Background()
	bill := s.createBillForTesting(ctx)
	approverID := uuid.New()

	req := &ApproveBillRequest{
		BillID: bill.ID,
	}

	result, err := s.approveBillUseCase.Execute(ctx, req, approverID, "192.168.1.1", "test-agent")
	s.Assert().NoError(err)
	s.Assert().NotNil(result)

	// Get audit logs (should be ordered by created_at DESC)
	audits, err := s.uow.BillAuditRepository().GetByBillID(ctx, bill.ID)
	s.Assert().NoError(err)
	s.Assert().Len(audits, 2)

	// First should be approval (most recent)
	s.Equal(domain.AuditActionApproved, audits[0].Action)
	// Second should be creation (oldest)
	s.Equal(domain.AuditActionCreated, audits[1].Action)

	// Verify timestamps are in correct order
	s.True(audits[0].CreatedAt.After(audits[1].CreatedAt))
}

// TestApproveBillUseCase_WhenBillApproved_ReturnBillDataPreserved tests that bill data is preserved.
func (s *ApproveBillTestSuite) TestApproveBillUseCase_WhenBillApproved_ReturnBillDataPreserved() {
	s.T().Parallel()

	ctx := context.Background()
	bill := s.createBillForTesting(ctx)
	approverID := uuid.New()

	// Get original bill data
	originalBill, err := s.uow.BillRepository().GetByID(ctx, bill.ID)
	s.Assert().NoError(err)

	req := &ApproveBillRequest{
		BillID: bill.ID,
	}

	result, err := s.approveBillUseCase.Execute(ctx, req, approverID, "192.168.1.1", "test-agent")
	s.Assert().NoError(err)
	s.Assert().NotNil(result)

	// Get updated bill
	updatedBill, err := s.uow.BillRepository().GetByID(ctx, bill.ID)
	s.Assert().NoError(err)

	// Verify data is preserved
	s.Equal(originalBill.Description, updatedBill.Description)
	s.Equal(originalBill.Amount, updatedBill.Amount)
	s.Equal(originalBill.DueDate, updatedBill.DueDate)
	s.Equal(originalBill.CreatedBy, updatedBill.CreatedBy)
	s.Equal(originalBill.CreatedAt, updatedBill.CreatedAt)

	// Verify only status and approver changed
	s.NotEqual(originalBill.Status, updatedBill.Status)
	s.Nil(originalBill.ApprovedBy)
	s.NotNil(updatedBill.ApprovedBy)
}
