package testutils

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/DATA-DOG/go-sqlmock"
)

// MockDB holds a mock database connection and mock controller
type MockDB struct {
	DB   *sql.DB
	Mock sqlmock.Sqlmock
}

// SetupMockDB creates a mock database connection for testing
func SetupMockDB() (*MockDB, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create mock database: %w", err)
	}

	return &MockDB{
		DB:   db,
		Mock: mock,
	}, nil
}

// Close closes the mock database connection
func (m *MockDB) Close() error {
	return m.DB.Close()
}

// ExpectPing sets up expectation for database ping
func (m *MockDB) ExpectPing() *sqlmock.ExpectedPing {
	return m.Mock.ExpectPing()
}

// ExpectQuery sets up expectation for database query
func (m *MockDB) ExpectQuery(query string) *sqlmock.ExpectedQuery {
	return m.Mock.ExpectQuery(query)
}

// ExpectExec sets up expectation for database exec
func (m *MockDB) ExpectExec(query string) *sqlmock.ExpectedExec {
	return m.Mock.ExpectExec(query)
}

// ExpectBegin sets up expectation for transaction begin
func (m *MockDB) ExpectBegin() *sqlmock.ExpectedBegin {
	return m.Mock.ExpectBegin()
}

// ExpectCommit sets up expectation for transaction commit
func (m *MockDB) ExpectCommit() *sqlmock.ExpectedCommit {
	return m.Mock.ExpectCommit()
}

// ExpectRollback sets up expectation for transaction rollback
func (m *MockDB) ExpectRollback() *sqlmock.ExpectedRollback {
	return m.Mock.ExpectRollback()
}

// ExpectationsWereMet checks if all expectations were met
func (m *MockDB) ExpectationsWereMet() error {
	return m.Mock.ExpectationsWereMet()
}

// SetupMockDBWithDatabase initializes the database package with mock database
func SetupMockDBWithDatabase(ctx context.Context) (*MockDB, func(), error) {
	mockDB, err := SetupMockDB()
	if err != nil {
		return nil, nil, err
	}

	// Set up basic ping expectation
	mockDB.ExpectPing()

	// We can't actually initialize the database package with a mock
	// So we'll return the mock for direct use in tests
	cleanup := func() {
		mockDB.Close()
	}

	return mockDB, cleanup, nil
}

// MockDatabaseTestSuite provides a test suite with mock database
type MockDatabaseTestSuite struct {
	MockDB  *MockDB
	Cleanup func()
	ctx     context.Context
}

// SetupSuite initializes the mock database test suite
func (s *MockDatabaseTestSuite) SetupSuite() {
	s.ctx = context.Background()
	
	mockDB, cleanup, err := SetupMockDBWithDatabase(s.ctx)
	if err != nil {
		panic(fmt.Sprintf("Failed to setup mock database: %v", err))
	}
	
	s.MockDB = mockDB
	s.Cleanup = cleanup
}

// TearDownSuite cleans up the mock database test suite
func (s *MockDatabaseTestSuite) TearDownSuite() {
	if s.Cleanup != nil {
		s.Cleanup()
	}
}

// SetupTest runs before each test
func (s *MockDatabaseTestSuite) SetupTest() {
	// Reset expectations if needed
}

// TearDownTest runs after each test
func (s *MockDatabaseTestSuite) TearDownTest() {
	// Verify all expectations were met
	if err := s.MockDB.ExpectationsWereMet(); err != nil {
		panic(fmt.Sprintf("Mock expectations were not met: %v", err))
	}
}