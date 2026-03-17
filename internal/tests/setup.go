package tests

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	migratepg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TestDatabase holds the database connection and container for cleanup.
type TestDatabase struct {
	DB        *gorm.DB
	Container testcontainers.Container
	DSN       string
}

// SetupTestDatabase creates a PostgreSQL container, runs migrations, and returns a connected database.
// This function handles:
// - Creating a PostgreSQL container
// - Waiting for the database to be ready
// - Running golang-migrate migrations
// - Validating migration files
// - Returning a ready-to-use GORM connection
func SetupTestDatabase(ctx context.Context) (*TestDatabase, error) {
	// Create PostgreSQL container
	container, err := createPostgresContainer(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create PostgreSQL container: %w", err)
	}

	// Get container connection details
	host, err := container.Host(ctx)
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, fmt.Errorf("failed to get container host: %w", err)
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, fmt.Errorf("failed to get container port: %w", err)
	}

	// Build DSN
	dsn := fmt.Sprintf(
		"host=%s port=%s user=test_user password=test_password dbname=test_db sslmode=disable",
		host,
		port.Port(),
	)

	// Connect to database
	db, err := gorm.Open(gormpostgres.Open(dsn), &gorm.Config{})
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Run migrations
	if err := runMigrations(db); err != nil {
		_ = container.Terminate(ctx)
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return &TestDatabase{
		DB:        db,
		Container: container,
		DSN:       dsn,
	}, nil
}

// createPostgresContainer creates and starts a PostgreSQL container.
func createPostgresContainer(ctx context.Context) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test_user",
			"POSTGRES_PASSWORD": "test_password",
			"POSTGRES_DB":       "test_db",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(30 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	return container, nil
}

// runMigrations executes all pending migrations using golang-migrate.
// This validates that:
// - Migration files exist and are valid
// - Database schema is correctly applied
// - All migrations complete successfully
func runMigrations(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	// Create migration driver
	driver, err := migratepg.WithInstance(sqlDB, &migratepg.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	// Create migration instance
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	// Run migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// Get current version
	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	if dirty {
		return fmt.Errorf("migration is in dirty state, version: %d", version)
	}

	return nil
}

// Cleanup terminates the container and closes the database connection.
func (td *TestDatabase) Cleanup(ctx context.Context) error {
	// Close database connection
	sqlDB, err := td.DB.DB()
	if err == nil {
		_ = sqlDB.Close()
	}

	// Terminate container
	if td.Container != nil {
		if err := td.Container.Terminate(ctx); err != nil {
			return fmt.Errorf("failed to terminate container: %w", err)
		}
	}

	return nil
}

// CleanupTables truncates all tables in the database.
// This is useful for cleaning up between tests while keeping the schema.
func (td *TestDatabase) CleanupTables(ctx context.Context) error {
	return td.DB.WithContext(ctx).Exec("TRUNCATE TABLE bill_audits CASCADE").Error
}

// GetDB returns the GORM database connection.
func (td *TestDatabase) GetDB() *gorm.DB {
	return td.DB
}

// GetDSN returns the database connection string.
func (td *TestDatabase) GetDSN() string {
	return td.DSN
}
