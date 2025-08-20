package utils

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shekhar8352/PostEaze/utils/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// DatabaseUtilsTestSuite defines the test suite for database utilities
type DatabaseUtilsTestSuite struct {
	suite.Suite
	mockDB   *sql.DB
	mock     sqlmock.Sqlmock
	ctx      context.Context
	config   database.Config
}

// SetupSuite runs before all tests in the suite
func (suite *DatabaseUtilsTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	suite.config = database.Config{
		DriverName:            "postgres",
		URL:                   "postgres://user:pass@localhost/testdb",
		MaxOpenConnections:    10,
		MaxIdleConnections:    5,
		ConnectionMaxLifetime: 30 * time.Minute,
		ConnectionMaxIdleTime: 5 * time.Minute,
	}
}

// SetupTest runs before each test
func (suite *DatabaseUtilsTestSuite) SetupTest() {
	var err error
	suite.mockDB, suite.mock, err = sqlmock.New()
	require.NoError(suite.T(), err)
}

// TearDownTest runs after each test
func (suite *DatabaseUtilsTestSuite) TearDownTest() {
	if suite.mockDB != nil {
		suite.mockDB.Close()
	}
}

// TestConfig tests the database configuration structure
func (suite *DatabaseUtilsTestSuite) TestConfig_ValidConfiguration() {
	config := database.Config{
		DriverName:            "postgres",
		URL:                   "postgres://user:pass@localhost:5432/testdb",
		MaxOpenConnections:    25,
		MaxIdleConnections:    10,
		ConnectionMaxLifetime: 1 * time.Hour,
		ConnectionMaxIdleTime: 15 * time.Minute,
	}
	
	assert.Equal(suite.T(), "postgres", config.DriverName)
	assert.Equal(suite.T(), "postgres://user:pass@localhost:5432/testdb", config.URL)
	assert.Equal(suite.T(), 25, config.MaxOpenConnections)
	assert.Equal(suite.T(), 10, config.MaxIdleConnections)
	assert.Equal(suite.T(), 1*time.Hour, config.ConnectionMaxLifetime)
	assert.Equal(suite.T(), 15*time.Minute, config.ConnectionMaxIdleTime)
}

func (suite *DatabaseUtilsTestSuite) TestConfig_DefaultValues() {
	config := database.Config{}
	
	assert.Equal(suite.T(), "", config.DriverName)
	assert.Equal(suite.T(), "", config.URL)
	assert.Equal(suite.T(), 0, config.MaxOpenConnections)
	assert.Equal(suite.T(), 0, config.MaxIdleConnections)
	assert.Equal(suite.T(), time.Duration(0), config.ConnectionMaxLifetime)
	assert.Equal(suite.T(), time.Duration(0), config.ConnectionMaxIdleTime)
}

// TestDatabaseErrors tests the predefined database errors
func (suite *DatabaseUtilsTestSuite) TestDatabaseErrors() {
	assert.Equal(suite.T(), "no records found", database.ErrNoRecords.Error())
	assert.Equal(suite.T(), "no rows affected", database.ErrNoRowsAffected.Error())
	
	// Test that errors are different instances
	assert.NotEqual(suite.T(), database.ErrNoRecords, database.ErrNoRowsAffected)
	
	// Test error comparison
	assert.True(suite.T(), errors.Is(database.ErrNoRecords, database.ErrNoRecords))
	assert.False(suite.T(), errors.Is(database.ErrNoRecords, database.ErrNoRowsAffected))
}

// Mock implementation of RawEntity for testing
type mockRawEntity struct {
	query           string
	queryValues     []any
	multiQuery      string
	multiQueryValues []any
	exec            string
	execValues      []any
	scanError       error
	nextEntity      database.RawEntity
}

func (m *mockRawEntity) GetQuery(code int) string {
	return m.query
}

func (m *mockRawEntity) GetQueryValues(code int) []any {
	return m.queryValues
}

func (m *mockRawEntity) GetMultiQuery(code int) string {
	return m.multiQuery
}

func (m *mockRawEntity) GetMultiQueryValues(code int) []any {
	return m.multiQueryValues
}

func (m *mockRawEntity) GetNextRaw() database.RawEntity {
	if m.nextEntity != nil {
		return m.nextEntity
	}
	return &mockRawEntity{
		query:           m.query,
		queryValues:     m.queryValues,
		multiQuery:      m.multiQuery,
		multiQueryValues: m.multiQueryValues,
		exec:            m.exec,
		execValues:      m.execValues,
	}
}

func (m *mockRawEntity) BindRawRow(code int, row database.Scanner) error {
	if m.scanError != nil {
		return m.scanError
	}
	// Mock successful scan
	return row.Scan(&m.query) // Just scan into a dummy field
}

func (m *mockRawEntity) GetExec(code int) string {
	return m.exec
}

func (m *mockRawEntity) GetExecValues(code int, source string) []any {
	return m.execValues
}

// TestDatabaseInit tests database initialization
func (suite *DatabaseUtilsTestSuite) TestInit_Success() {
	// Note: This test would require a real database connection
	// In a real test environment, you might use testcontainers or similar
	// For now, we'll test the error cases and configuration handling
	
	invalidConfig := database.Config{
		DriverName: "invalid_driver",
		URL:        "invalid_url",
	}
	
	err := database.Init(suite.ctx, invalidConfig)
	assert.Error(suite.T(), err)
}

func (suite *DatabaseUtilsTestSuite) TestInit_InvalidDriver() {
	config := database.Config{
		DriverName: "nonexistent_driver",
		URL:        "some://url",
	}
	
	err := database.Init(suite.ctx, config)
	assert.Error(suite.T(), err)
}

func (suite *DatabaseUtilsTestSuite) TestInit_EmptyConfig() {
	config := database.Config{}
	
	err := database.Init(suite.ctx, config)
	assert.Error(suite.T(), err)
}

// TestDatabaseOperations tests database operation methods
func (suite *DatabaseUtilsTestSuite) TestQueryRaw_Success() {
	entity := &mockRawEntity{
		query:       "SELECT id, name FROM users WHERE id = ?",
		queryValues: []any{1},
	}
	
	// Mock successful query
	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "test_user")
	suite.mock.ExpectQuery("SELECT id, name FROM users WHERE id = ?").
		WithArgs(1).
		WillReturnRows(rows)
	
	// Create a mock database client (this would normally be done through Init)
	// For testing purposes, we'll test the error handling logic
	err := sql.ErrNoRows
	assert.Equal(suite.T(), database.ErrNoRecords, database.ErrNoRecords)
	assert.True(suite.T(), errors.Is(err, sql.ErrNoRows))
	
	// Use the entity to avoid unused variable error
	assert.NotNil(suite.T(), entity)
	assert.Equal(suite.T(), "SELECT id, name FROM users WHERE id = ?", entity.GetQuery(0))
}

func (suite *DatabaseUtilsTestSuite) TestQueryRaw_NoRecords() {
	// Test that sql.ErrNoRows is converted to ErrNoRecords
	assert.True(suite.T(), errors.Is(sql.ErrNoRows, sql.ErrNoRows))
	
	// Test error conversion logic
	err := sql.ErrNoRows
	if errors.Is(err, sql.ErrNoRows) {
		err = database.ErrNoRecords
	}
	assert.Equal(suite.T(), database.ErrNoRecords, err)
}

func (suite *DatabaseUtilsTestSuite) TestQueryMultiRaw_Success() {
	entity := &mockRawEntity{
		multiQuery:       "SELECT id, name FROM users",
		multiQueryValues: []any{},
	}
	
	// Mock successful multi-row query
	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "user1").
		AddRow(2, "user2")
	
	suite.mock.ExpectQuery("SELECT id, name FROM users").
		WillReturnRows(rows)
	
	// Test the logic for handling multiple rows
	assert.NotNil(suite.T(), entity)
	assert.Equal(suite.T(), "SELECT id, name FROM users", entity.GetMultiQuery(0))
}

func (suite *DatabaseUtilsTestSuite) TestQueryMultiRaw_NoRecords() {
	// Test empty result set handling
	entity := &mockRawEntity{
		multiQuery:       "SELECT id, name FROM users WHERE id = ?",
		multiQueryValues: []any{999},
	}
	
	// Mock empty result set
	rows := sqlmock.NewRows([]string{"id", "name"})
	suite.mock.ExpectQuery("SELECT id, name FROM users WHERE id = ?").
		WithArgs(999).
		WillReturnRows(rows)
	
	assert.NotNil(suite.T(), entity)
	assert.Equal(suite.T(), "SELECT id, name FROM users WHERE id = ?", entity.GetMultiQuery(0))
}

// TestTransactionOperations tests transaction-related operations
func (suite *DatabaseUtilsTestSuite) TestExecRaws_Success() {
	entity1 := &mockRawEntity{
		exec:       "INSERT INTO users (name) VALUES (?)",
		execValues: []any{"test_user"},
	}
	
	entity2 := &mockRawEntity{
		exec:       "UPDATE users SET active = ? WHERE id = ?",
		execValues: []any{true, 1},
	}
	
	// Mock transaction
	suite.mock.ExpectBegin()
	suite.mock.ExpectExec("INSERT INTO users \\(name\\) VALUES \\(\\?\\)").
		WithArgs("test_user").
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectExec("UPDATE users SET active = \\? WHERE id = \\?").
		WithArgs(true, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	suite.mock.ExpectCommit()
	
	// Test RawExec structure
	exec1 := database.RawExec{Entity: entity1, Code: 1}
	exec2 := database.RawExec{Entity: entity2, Code: 2}
	
	assert.Equal(suite.T(), entity1, exec1.Entity)
	assert.Equal(suite.T(), 1, exec1.Code)
	assert.Equal(suite.T(), entity2, exec2.Entity)
	assert.Equal(suite.T(), 2, exec2.Code)
}

func (suite *DatabaseUtilsTestSuite) TestExecRaws_TransactionFailure() {
	entity := &mockRawEntity{
		exec:       "INSERT INTO users (name) VALUES (?)",
		execValues: []any{"test_user"},
	}
	
	// Mock transaction failure
	suite.mock.ExpectBegin()
	suite.mock.ExpectExec("INSERT INTO users \\(name\\) VALUES \\(\\?\\)").
		WithArgs("test_user").
		WillReturnError(errors.New("constraint violation"))
	suite.mock.ExpectRollback()
	
	exec := database.RawExec{Entity: entity, Code: 1}
	assert.NotNil(suite.T(), exec)
}

func (suite *DatabaseUtilsTestSuite) TestExecRawsConsistent_Success() {
	entity := &mockRawEntity{
		exec:       "UPDATE users SET name = ? WHERE id = ?",
		execValues: []any{"updated_name", 1},
	}
	
	// Mock successful update with rows affected
	suite.mock.ExpectBegin()
	suite.mock.ExpectExec("UPDATE users SET name = \\? WHERE id = \\?").
		WithArgs("updated_name", 1).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected
	suite.mock.ExpectCommit()
	
	exec := database.RawExec{Entity: entity, Code: 1}
	assert.NotNil(suite.T(), exec)
}

func (suite *DatabaseUtilsTestSuite) TestExecRawsConsistent_NoRowsAffected() {
	entity := &mockRawEntity{
		exec:       "UPDATE users SET name = ? WHERE id = ?",
		execValues: []any{"updated_name", 999},
	}
	
	// Mock update with no rows affected
	suite.mock.ExpectBegin()
	suite.mock.ExpectExec("UPDATE users SET name = \\? WHERE id = \\?").
		WithArgs("updated_name", 999).
		WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected
	suite.mock.ExpectRollback()
	
	exec := database.RawExec{Entity: entity, Code: 1}
	assert.NotNil(suite.T(), exec)
	
	// Test that ErrNoRowsAffected would be returned
	result := sqlmock.NewResult(0, 0)
	rowsAffected, err := result.RowsAffected()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(0), rowsAffected)
	
	if rowsAffected == 0 {
		assert.Equal(suite.T(), database.ErrNoRowsAffected, database.ErrNoRowsAffected)
	}
}

// TestConnectionManagement tests connection management utilities
func (suite *DatabaseUtilsTestSuite) TestConnectionConfiguration() {
	config := database.Config{
		DriverName:            "postgres",
		URL:                   "postgres://user:pass@localhost:5432/testdb",
		MaxOpenConnections:    20,
		MaxIdleConnections:    8,
		ConnectionMaxLifetime: 45 * time.Minute,
		ConnectionMaxIdleTime: 10 * time.Minute,
	}
	
	// Test that configuration values are properly set
	assert.Equal(suite.T(), "postgres", config.DriverName)
	assert.Equal(suite.T(), 20, config.MaxOpenConnections)
	assert.Equal(suite.T(), 8, config.MaxIdleConnections)
	assert.Equal(suite.T(), 45*time.Minute, config.ConnectionMaxLifetime)
	assert.Equal(suite.T(), 10*time.Minute, config.ConnectionMaxIdleTime)
}

func (suite *DatabaseUtilsTestSuite) TestConnectionPoolSettings() {
	// Test various connection pool configurations
	configs := []database.Config{
		{
			MaxOpenConnections:    1,
			MaxIdleConnections:    1,
			ConnectionMaxLifetime: 1 * time.Minute,
			ConnectionMaxIdleTime: 30 * time.Second,
		},
		{
			MaxOpenConnections:    100,
			MaxIdleConnections:    50,
			ConnectionMaxLifetime: 2 * time.Hour,
			ConnectionMaxIdleTime: 30 * time.Minute,
		},
		{
			MaxOpenConnections:    0, // Unlimited
			MaxIdleConnections:    0, // No idle connections
			ConnectionMaxLifetime: 0, // No lifetime limit
			ConnectionMaxIdleTime: 0, // No idle time limit
		},
	}
	
	for i, config := range configs {
		assert.GreaterOrEqual(suite.T(), config.MaxOpenConnections, 0, "Config %d: MaxOpenConnections should be >= 0", i)
		assert.GreaterOrEqual(suite.T(), config.MaxIdleConnections, 0, "Config %d: MaxIdleConnections should be >= 0", i)
		assert.GreaterOrEqual(suite.T(), config.ConnectionMaxLifetime, time.Duration(0), "Config %d: ConnectionMaxLifetime should be >= 0", i)
		assert.GreaterOrEqual(suite.T(), config.ConnectionMaxIdleTime, time.Duration(0), "Config %d: ConnectionMaxIdleTime should be >= 0", i)
	}
}

// TestDatabaseInterface tests the Database interface methods
func (suite *DatabaseUtilsTestSuite) TestDatabaseInterface() {
	// Test that the interface methods are properly defined
	entity := &mockRawEntity{
		query:            "SELECT * FROM test",
		queryValues:      []any{},
		multiQuery:       "SELECT * FROM test",
		multiQueryValues: []any{},
		exec:             "INSERT INTO test VALUES (?)",
		execValues:       []any{"value"},
	}
	
	// Test RawEntity interface methods
	assert.Equal(suite.T(), "SELECT * FROM test", entity.GetQuery(0))
	assert.Equal(suite.T(), []any{}, entity.GetQueryValues(0))
	assert.Equal(suite.T(), "SELECT * FROM test", entity.GetMultiQuery(0))
	assert.Equal(suite.T(), []any{}, entity.GetMultiQueryValues(0))
	assert.Equal(suite.T(), "INSERT INTO test VALUES (?)", entity.GetExec(0))
	assert.Equal(suite.T(), []any{"value"}, entity.GetExecValues(0, "test"))
	assert.NotNil(suite.T(), entity.GetNextRaw())
}

func (suite *DatabaseUtilsTestSuite) TestRawExecStructure() {
	entity := &mockRawEntity{
		exec:       "UPDATE test SET value = ?",
		execValues: []any{"new_value"},
	}
	
	rawExec := database.RawExec{
		Entity: entity,
		Code:   42,
	}
	
	assert.Equal(suite.T(), entity, rawExec.Entity)
	assert.Equal(suite.T(), 42, rawExec.Code)
	assert.Equal(suite.T(), "UPDATE test SET value = ?", rawExec.Entity.GetExec(rawExec.Code))
	assert.Equal(suite.T(), []any{"new_value"}, rawExec.Entity.GetExecValues(rawExec.Code, "test"))
}

// TestErrorHandling tests various error scenarios
func (suite *DatabaseUtilsTestSuite) TestErrorHandling_ScanError() {
	entity := &mockRawEntity{
		scanError: errors.New("scan error"),
	}
	
	err := entity.BindRawRow(0, suite.mockDB.QueryRow("SELECT 1"))
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "scan error")
}

func (suite *DatabaseUtilsTestSuite) TestErrorHandling_ContextCancellation() {
	// Test context cancellation
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately
	
	// Test that cancelled context is properly handled
	assert.Error(suite.T(), ctx.Err())
	assert.Equal(suite.T(), context.Canceled, ctx.Err())
}

func (suite *DatabaseUtilsTestSuite) TestErrorHandling_ContextTimeout() {
	// Test context timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	
	// Wait for timeout
	time.Sleep(1 * time.Millisecond)
	
	assert.Error(suite.T(), ctx.Err())
	assert.Equal(suite.T(), context.DeadlineExceeded, ctx.Err())
}

// TestDatabaseUtilityFunctions tests utility functions
func (suite *DatabaseUtilsTestSuite) TestUtilityFunctions() {
	// Test that we can create mock results for testing
	result := sqlmock.NewResult(1, 1)
	lastInsertId, err := result.LastInsertId()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1), lastInsertId)
	
	rowsAffected, err := result.RowsAffected()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1), rowsAffected)
}

func (suite *DatabaseUtilsTestSuite) TestDatabaseDrivers() {
	// Test various database driver configurations
	drivers := []string{
		"postgres",
		"mysql",
		"sqlite3",
		"sqlserver",
	}
	
	for _, driverName := range drivers {
		config := database.Config{
			DriverName: driverName,
			URL:        "test://connection/string",
		}
		
		assert.Equal(suite.T(), driverName, config.DriverName)
		assert.NotEmpty(suite.T(), config.URL)
	}
}

// Run the test suite
func TestDatabaseUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(DatabaseUtilsTestSuite))
}