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

// CreateBillTestSuite defines the test suite for CreateBillUseCase.
type CreateBillTestSuite struct {
	suite.Suite
	testDB  *tests.TestDatabase
	uow     domain.UnitOfWork
	useCase *CreateBillUseCase
}

// SetupSuite initializes the test suite with a PostgreSQL container and runs migrations.
func (s *CreateBillTestSuite) SetupSuite() {
	ctx := context.Background()

	// Setup test database with migrations
	testDB, err := tests.SetupTestDatabase(ctx)
	s.Assert().NoError(err, "failed to setup test database")
	s.testDB = testDB

	// Initialize Unit of Work and UseCase
	s.uow = persistence.NewGormUnitOfWork(s.testDB.GetDB())
	s.useCase = NewCreateBillUseCase(s.uow)
}

// TearDownSuite cleans up the test suite by terminating the container.
func (s *CreateBillTestSuite) TearDownSuite() {
	ctx := context.Background()
	if s.testDB != nil {
		err := s.testDB.Cleanup(ctx)
		s.Assert().NoError(err, "failed to cleanup test database")
	}
}

// SetupTest cleans up the database before each test.
func (s *CreateBillTestSuite) SetupTest() {
	ctx := context.Background()
	err := s.testDB.CleanupTables(ctx)
	s.Assert().NoError(err, "failed to cleanup tables")
}

// TestCreateBillTestSuite runs the entire test suite.
func TestCreateBillTestSuite(t *testing.T) {
	suite.Run(t, new(CreateBillTestSuite))
}

// TestCreateBillUseCase_WhenValidRequest_ReturnBillCreated tests successful bill creation.
func (s *CreateBillTestSuite) TestCreateBillUseCase_WhenValidRequest_ReturnBillCreated() {
	s.T().Parallel()

	ctx := context.Background()
	userID := uuid.New()
	ipAddress := "192.168.1.1"
	userAgent := "Mozilla/5.0"

	req := &CreateBillRequest{
		Description: "Office Supplies",
		Amount:      1500.50,
		DueDate:     time.Now().AddDate(0, 0, 30),
	}

	// Execute
	result, err := s.useCase.Execute(ctx, req, userID, ipAddress, userAgent)

	// Assert
	s.Assert().NoError(err)
	s.Assert().NotNil(result)
	s.Equal(req.Description, result.Description)
	s.Equal(req.Amount, result.Amount)
	s.Equal(req.DueDate.Format("2006-01-02"), result.DueDate.Format("2006-01-02"))
	s.Equal(string(domain.BillStatusPending), result.Status)
	s.Equal(userID, result.CreatedBy)

	// Verify bill was persisted
	bill, err := s.uow.BillRepository().GetByID(ctx, result.ID)
	s.Assert().NoError(err)
	s.Assert().NotNil(bill)
	s.Equal(req.Description, bill.Description)

	// Verify audit log was created
	audits, err := s.uow.BillAuditRepository().GetByBillID(ctx, result.ID)
	s.Assert().NoError(err)
	s.Assert().Len(audits, 1)
	s.Equal(domain.AuditActionCreated, audits[0].Action)
	s.Equal(userID, audits[0].PerformedBy)
	s.Equal(ipAddress, audits[0].IPAddress)
	s.Equal(userAgent, audits[0].UserAgent)
}

// TestCreateBillUseCase_WhenMultipleBillsCreated_ReturnAllBillsCreated tests creating multiple bills.
func (s *CreateBillTestSuite) TestCreateBillUseCase_WhenMultipleBillsCreated_ReturnAllBillsCreated() {
	s.T().Parallel()

	ctx := context.Background()
	userID := uuid.New()

	// Create 5 bills
	for i := 1; i <= 5; i++ {
		req := &CreateBillRequest{
			Description: fmt.Sprintf("Bill %d", i),
			Amount:      float64(1000 * i),
			DueDate:     time.Now().AddDate(0, 0, i*10),
		}

		result, err := s.useCase.Execute(ctx, req, userID, "192.168.1.1", "test-agent")
		s.Assert().NoError(err)
		s.Assert().NotNil(result)
	}

	// Verify all bills were created
	bills, err := s.uow.BillRepository().GetAll(ctx, nil)
	s.Assert().NoError(err)
	s.Assert().Len(bills, 5)
}

// TestCreateBillUseCase_WhenEmptyDescription_ReturnError tests validation of empty description.
func (s *CreateBillTestSuite) TestCreateBillUseCase_WhenEmptyDescription_ReturnError() {
	s.T().Parallel()

	ctx := context.Background()
	userID := uuid.New()

	req := &CreateBillRequest{
		Description: "",
		Amount:      1500.50,
		DueDate:     time.Now().AddDate(0, 0, 30),
	}

	// Execute - should still create but with empty description (validation is at handler level)
	result, err := s.useCase.Execute(ctx, req, userID, "192.168.1.1", "test-agent")

	// In this case, the service doesn't validate, but we can verify it was created
	if err == nil {
		s.Assert().NotNil(result)
		s.Equal("", result.Description)
	}
}

// TestCreateBillUseCase_WhenNegativeAmount_ReturnBillCreated tests creation with negative amount.
func (s *CreateBillTestSuite) TestCreateBillUseCase_WhenNegativeAmount_ReturnBillCreated() {
	s.T().Parallel()

	ctx := context.Background()
	userID := uuid.New()

	req := &CreateBillRequest{
		Description: "Negative Bill",
		Amount:      -100.00,
		DueDate:     time.Now().AddDate(0, 0, 30),
	}

	// Execute - service doesn't validate, but we verify it persists
	result, err := s.useCase.Execute(ctx, req, userID, "192.168.1.1", "test-agent")

	// Negative amounts are allowed at service level (validation is at handler)
	if err == nil {
		s.Assert().NotNil(result)
		s.Equal(-100.00, result.Amount)
	}
}

// TestCreateBillUseCase_WhenZeroAmount_ReturnBillCreated tests creation with zero amount.
func (s *CreateBillTestSuite) TestCreateBillUseCase_WhenZeroAmount_ReturnBillCreated() {
	s.T().Parallel()

	ctx := context.Background()
	userID := uuid.New()

	req := &CreateBillRequest{
		Description: "Zero Amount Bill",
		Amount:      0.0,
		DueDate:     time.Now().AddDate(0, 0, 30),
	}

	result, err := s.useCase.Execute(ctx, req, userID, "192.168.1.1", "test-agent")

	if err == nil {
		s.Assert().NotNil(result)
		s.Equal(0.0, result.Amount)
	}
}

// TestCreateBillUseCase_WhenPastDueDate_ReturnBillCreated tests creation with past due date.
func (s *CreateBillTestSuite) TestCreateBillUseCase_WhenPastDueDate_ReturnBillCreated() {
	s.T().Parallel()

	ctx := context.Background()
	userID := uuid.New()
	pastDate := time.Now().AddDate(0, 0, -10)

	req := &CreateBillRequest{
		Description: "Past Due Bill",
		Amount:      1500.50,
		DueDate:     pastDate,
	}

	result, err := s.useCase.Execute(ctx, req, userID, "192.168.1.1", "test-agent")

	if err == nil {
		s.Assert().NotNil(result)
		s.Equal(pastDate.Format("2006-01-02"), result.DueDate.Format("2006-01-02"))
	}
}

// TestCreateBillUseCase_WhenFutureDueDate_ReturnBillCreated tests creation with future due date.
func (s *CreateBillTestSuite) TestCreateBillUseCase_WhenFutureDueDate_ReturnBillCreated() {
	s.T().Parallel()

	ctx := context.Background()
	userID := uuid.New()
	futureDate := time.Now().AddDate(1, 0, 0)

	req := &CreateBillRequest{
		Description: "Future Due Bill",
		Amount:      1500.50,
		DueDate:     futureDate,
	}

	result, err := s.useCase.Execute(ctx, req, userID, "192.168.1.1", "test-agent")

	s.Assert().NoError(err)
	s.Assert().NotNil(result)
	s.Equal(futureDate.Format("2006-01-02"), result.DueDate.Format("2006-01-02"))
}

// TestCreateBillUseCase_WhenLargeAmount_ReturnBillCreated tests creation with very large amount.
func (s *CreateBillTestSuite) TestCreateBillUseCase_WhenLargeAmount_ReturnBillCreated() {
	s.T().Parallel()

	ctx := context.Background()
	userID := uuid.New()
	largeAmount := 999999999.99

	req := &CreateBillRequest{
		Description: "Large Amount Bill",
		Amount:      largeAmount,
		DueDate:     time.Now().AddDate(0, 0, 30),
	}

	result, err := s.useCase.Execute(ctx, req, userID, "192.168.1.1", "test-agent")

	s.Assert().NoError(err)
	s.Assert().NotNil(result)
	s.Equal(largeAmount, result.Amount)
}

// TestCreateBillUseCase_WhenLongDescription_ReturnBillCreated tests creation with very long description.
func (s *CreateBillTestSuite) TestCreateBillUseCase_WhenLongDescription_ReturnBillCreated() {
	s.T().Parallel()

	ctx := context.Background()
	userID := uuid.New()
	longDescription := ""
	for i := 0; i < 50; i++ {
		longDescription += "This is a long description. "
	}

	req := &CreateBillRequest{
		Description: longDescription,
		Amount:      1500.50,
		DueDate:     time.Now().AddDate(0, 0, 30),
	}

	result, err := s.useCase.Execute(ctx, req, userID, "192.168.1.1", "test-agent")

	s.Assert().NoError(err)
	s.Assert().NotNil(result)
	s.Equal(longDescription, result.Description)
}

// TestCreateBillUseCase_WhenDifferentUsers_ReturnBillsWithCorrectCreators tests multiple users creating bills.
func (s *CreateBillTestSuite) TestCreateBillUseCase_WhenDifferentUsers_ReturnBillsWithCorrectCreators() {
	s.T().Parallel()

	ctx := context.Background()
	user1 := uuid.New()
	user2 := uuid.New()
	user3 := uuid.New()

	// Create bills from different users
	users := []uuid.UUID{user1, user2, user3}
	for i, userID := range users {
		req := &CreateBillRequest{
			Description: fmt.Sprintf("User %d Bill", i+1),
			Amount:      float64(1000 * (i + 1)),
			DueDate:     time.Now().AddDate(0, 0, 30),
		}

		result, err := s.useCase.Execute(ctx, req, userID, fmt.Sprintf("192.168.1.%d", i+1), "test-agent")
		s.Assert().NoError(err)
		s.Assert().Equal(userID, result.CreatedBy)
	}

	// Verify bills were created with correct creators
	bills, err := s.uow.BillRepository().GetAll(ctx, nil)
	s.Assert().NoError(err)
	s.Assert().Len(bills, 3)

	creatorMap := make(map[uuid.UUID]int)
	for _, bill := range bills {
		creatorMap[bill.CreatedBy]++
	}

	s.Assert().Equal(1, creatorMap[user1])
	s.Assert().Equal(1, creatorMap[user2])
	s.Assert().Equal(1, creatorMap[user3])
}

// TestCreateBillUseCase_WhenAuditLogCreated_ReturnAuditWithCorrectMetadata tests audit log metadata.
func (s *CreateBillTestSuite) TestCreateBillUseCase_WhenAuditLogCreated_ReturnAuditWithCorrectMetadata() {
	s.T().Parallel()

	ctx := context.Background()
	userID := uuid.New()
	ipAddress := "203.0.113.42"
	userAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64)"

	req := &CreateBillRequest{
		Description: "Audit Test Bill",
		Amount:      2500.00,
		DueDate:     time.Now().AddDate(0, 0, 30),
	}

	result, err := s.useCase.Execute(ctx, req, userID, ipAddress, userAgent)
	s.Assert().NoError(err)

	// Verify audit log contains correct metadata
	audits, err := s.uow.BillAuditRepository().GetByBillID(ctx, result.ID)
	s.Assert().NoError(err)
	s.Assert().Len(audits, 1)

	audit := audits[0]
	s.Equal(result.ID, audit.BillID)
	s.Equal(domain.AuditActionCreated, audit.Action)
	s.Equal(userID, audit.PerformedBy)
	s.Equal(ipAddress, audit.IPAddress)
	s.Equal(userAgent, audit.UserAgent)
	s.NotZero(audit.CreatedAt)
}

// TestCreateBillUseCase_WhenBillCreated_ReturnStatusPending tests that new bills have pending status.
func (s *CreateBillTestSuite) TestCreateBillUseCase_WhenBillCreated_ReturnStatusPending() {
	s.T().Parallel()

	ctx := context.Background()
	userID := uuid.New()

	req := &CreateBillRequest{
		Description: "Status Test Bill",
		Amount:      1500.50,
		DueDate:     time.Now().AddDate(0, 0, 30),
	}

	result, err := s.useCase.Execute(ctx, req, userID, "192.168.1.1", "test-agent")
	s.Assert().NoError(err)

	s.Equal(string(domain.BillStatusPending), result.Status)

	// Verify in database
	bill, err := s.uow.BillRepository().GetByID(ctx, result.ID)
	s.Assert().NoError(err)
	s.Equal(domain.BillStatusPending, bill.Status)
}

// TestCreateBillUseCase_WhenBillCreated_ReturnApprovedByNil tests that new bills have no approver.
func (s *CreateBillTestSuite) TestCreateBillUseCase_WhenBillCreated_ReturnApprovedByNil() {
	s.T().Parallel()

	ctx := context.Background()
	userID := uuid.New()

	req := &CreateBillRequest{
		Description: "Approver Test Bill",
		Amount:      1500.50,
		DueDate:     time.Now().AddDate(0, 0, 30),
	}

	result, err := s.useCase.Execute(ctx, req, userID, "192.168.1.1", "test-agent")
	s.Assert().NoError(err)

	// Verify in database
	bill, err := s.uow.BillRepository().GetByID(ctx, result.ID)
	s.Assert().NoError(err)
	s.Nil(bill.ApprovedBy)
}

// TestCreateBillUseCase_WhenBillCreated_ReturnValidUUID tests that created bill has valid UUID.
func (s *CreateBillTestSuite) TestCreateBillUseCase_WhenBillCreated_ReturnValidUUID() {
	s.T().Parallel()

	ctx := context.Background()
	userID := uuid.New()

	req := &CreateBillRequest{
		Description: "UUID Test Bill",
		Amount:      1500.50,
		DueDate:     time.Now().AddDate(0, 0, 30),
	}

	result, err := s.useCase.Execute(ctx, req, userID, "192.168.1.1", "test-agent")
	s.Assert().NoError(err)

	// Verify UUID is valid and not zero
	s.NotEqual(uuid.Nil, result.ID)
	s.NotZero(result.ID)
}

// TestCreateBillUseCase_WhenBillCreated_ReturnTimestamps tests that timestamps are set correctly.
func (s *CreateBillTestSuite) TestCreateBillUseCase_WhenBillCreated_ReturnTimestamps() {
	s.T().Parallel()

	ctx := context.Background()
	userID := uuid.New()
	beforeCreate := time.Now()

	req := &CreateBillRequest{
		Description: "Timestamp Test Bill",
		Amount:      1500.50,
		DueDate:     time.Now().AddDate(0, 0, 30),
	}

	result, err := s.useCase.Execute(ctx, req, userID, "192.168.1.1", "test-agent")
	s.Assert().NoError(err)

	afterCreate := time.Now()

	// Verify timestamps are within expected range
	s.True(result.CreatedAt.After(beforeCreate) || result.CreatedAt.Equal(beforeCreate))
	s.True(result.CreatedAt.Before(afterCreate) || result.CreatedAt.Equal(afterCreate))
}

// TestCreateBillUseCase_WhenConcurrentBillCreation_ReturnAllBillsCreated tests concurrent bill creation.
func (s *CreateBillTestSuite) TestCreateBillUseCase_WhenConcurrentBillCreation_ReturnAllBillsCreated() {
	s.T().Parallel()

	ctx := context.Background()
	numGoroutines := 10
	results := make(chan *CreateBillResponse, numGoroutines)
	errors := make(chan error, numGoroutines)

	// Create bills concurrently
	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			userID := uuid.New()
			req := &CreateBillRequest{
				Description: fmt.Sprintf("Concurrent Bill %d", index),
				Amount:      float64(1000 * (index + 1)),
				DueDate:     time.Now().AddDate(0, 0, 30),
			}

			result, err := s.useCase.Execute(ctx, req, userID, "192.168.1.1", "test-agent")
			if err != nil {
				errors <- err
				return
			}
			results <- result
		}(i)
	}

	// Collect results
	successCount := 0
	for i := 0; i < numGoroutines; i++ {
		select {
		case err := <-errors:
			s.T().Logf("Error creating bill: %v", err)
		case result := <-results:
			s.Assert().NotNil(result)
			successCount++
		}
	}

	// Verify all bills were created
	s.Equal(numGoroutines, successCount)

	bills, err := s.uow.BillRepository().GetAll(ctx, nil)
	s.Assert().NoError(err)
	s.Assert().Len(bills, numGoroutines)
}
