package mocks

import (
	"context"
	"errors"
	"fmt"

	"github.com/shekhar8352/PostEaze/utils/database"
	"github.com/stretchr/testify/mock"
)

// MockDatabase is a mock implementation of the Database interface
type MockDatabase struct {
	mock.Mock
}

// QueryRaw mocks the QueryRaw method
func (m *MockDatabase) QueryRaw(ctx context.Context, entity database.RawEntity, code int) error {
	args := m.Called(ctx, entity, code)
	return args.Error(0)
}

// QueryMultiRaw mocks the QueryMultiRaw method
func (m *MockDatabase) QueryMultiRaw(ctx context.Context, entity database.RawEntity, code int) ([]database.RawEntity, error) {
	args := m.Called(ctx, entity, code)
	return args.Get(0).([]database.RawEntity), args.Error(1)
}

// ExecRaws mocks the ExecRaws method
func (m *MockDatabase) ExecRaws(ctx context.Context, source string, execs ...database.RawExec) error {
	args := m.Called(ctx, source, execs)
	return args.Error(0)
}

// ExecRawsConsistent mocks the ExecRawsConsistent method
func (m *MockDatabase) ExecRawsConsistent(ctx context.Context, source string, execs ...database.RawExec) error {
	args := m.Called(ctx, source, execs)
	return args.Error(0)
}

// NewMockDatabase creates a new mock database instance
func NewMockDatabase() *MockDatabase {
	return &MockDatabase{}
}

// MockDatabaseManager provides utilities for managing database mocks in tests
type MockDatabaseManager struct {
	mockDB *MockDatabase
	mockTx *MockDatabase
}

// NewMockDatabaseManager creates a new database manager for testing
func NewMockDatabaseManager() *MockDatabaseManager {
	return &MockDatabaseManager{
		mockDB: NewMockDatabase(),
		mockTx: NewMockDatabase(),
	}
}

// GetMockDB returns the mock database instance
func (m *MockDatabaseManager) GetMockDB() *MockDatabase {
	return m.mockDB
}

// GetMockTx returns the mock transaction instance
func (m *MockDatabaseManager) GetMockTx() *MockDatabase {
	return m.mockTx
}

// SetupSuccessfulQuery configures a mock to return successful query results
func (m *MockDatabaseManager) SetupSuccessfulQuery(entity database.RawEntity, code int) {
	m.mockDB.On("QueryRaw", mock.Anything, entity, code).Return(nil)
}

// SetupFailedQuery configures a mock to return query errors
func (m *MockDatabaseManager) SetupFailedQuery(entity database.RawEntity, code int, err error) {
	m.mockDB.On("QueryRaw", mock.Anything, entity, code).Return(err)
}

// SetupSuccessfulMultiQuery configures a mock to return successful multi-query results
func (m *MockDatabaseManager) SetupSuccessfulMultiQuery(entity database.RawEntity, code int, results []database.RawEntity) {
	m.mockDB.On("QueryMultiRaw", mock.Anything, entity, code).Return(results, nil)
}

// SetupFailedMultiQuery configures a mock to return multi-query errors
func (m *MockDatabaseManager) SetupFailedMultiQuery(entity database.RawEntity, code int, err error) {
	m.mockDB.On("QueryMultiRaw", mock.Anything, entity, code).Return([]database.RawEntity(nil), err)
}

// SetupSuccessfulExec configures a mock to return successful exec results
func (m *MockDatabaseManager) SetupSuccessfulExec(source string, execs ...database.RawExec) {
	m.mockDB.On("ExecRaws", mock.Anything, source, execs).Return(nil)
}

// SetupFailedExec configures a mock to return exec errors
func (m *MockDatabaseManager) SetupFailedExec(source string, err error, execs ...database.RawExec) {
	m.mockDB.On("ExecRaws", mock.Anything, source, execs).Return(err)
}

// SetupSuccessfulTransaction configures mocks for successful transaction operations
func (m *MockDatabaseManager) SetupSuccessfulTransaction() {
	m.mockTx.On("QueryRaw", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	m.mockTx.On("ExecRaws", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	m.mockTx.On("ExecRawsConsistent", mock.Anything, mock.Anything, mock.Anything).Return(nil)
}

// SetupFailedTransaction configures mocks for failed transaction operations
func (m *MockDatabaseManager) SetupFailedTransaction(err error) {
	m.mockTx.On("QueryRaw", mock.Anything, mock.Anything, mock.Anything).Return(err)
	m.mockTx.On("ExecRaws", mock.Anything, mock.Anything, mock.Anything).Return(err)
	m.mockTx.On("ExecRawsConsistent", mock.Anything, mock.Anything, mock.Anything).Return(err)
}

// VerifyQueryCalled verifies that QueryRaw was called with expected parameters
func (m *MockDatabaseManager) VerifyQueryCalled(entity database.RawEntity, code int) bool {
	for _, call := range m.mockDB.Calls {
		if call.Method == "QueryRaw" && len(call.Arguments) >= 3 {
			if call.Arguments[1] == entity && call.Arguments[2] == code {
				return true
			}
		}
	}
	return false
}

// VerifyExecCalled verifies that ExecRaws was called with expected parameters
func (m *MockDatabaseManager) VerifyExecCalled(source string, execs ...database.RawExec) bool {
	for _, call := range m.mockDB.Calls {
		if call.Method == "ExecRaws" && len(call.Arguments) >= 3 {
			if call.Arguments[1] == source {
				return true
			}
		}
	}
	return false
}

// AssertExpectations asserts that all expectations were met
func (m *MockDatabaseManager) AssertExpectations(t mock.TestingT) {
	m.mockDB.AssertExpectations(t)
	m.mockTx.AssertExpectations(t)
}

// Reset clears all expectations and call history
func (m *MockDatabaseManager) Reset() {
	m.mockDB.ExpectedCalls = nil
	m.mockDB.Calls = nil
	m.mockTx.ExpectedCalls = nil
	m.mockTx.Calls = nil
}

// Common database errors for testing
var (
	ErrMockNoRecords      = errors.New("mock: no records found")
	ErrMockNoRowsAffected = errors.New("mock: no rows affected")
	ErrMockConnection     = errors.New("mock: database connection failed")
	ErrMockTransaction    = errors.New("mock: transaction failed")
	ErrMockConstraint     = errors.New("mock: constraint violation")
)

// MockRepositoryHelper provides utilities for mocking repository operations
type MockRepositoryHelper struct {
	dbManager *MockDatabaseManager
}

// NewMockRepositoryHelper creates a new repository helper for testing
func NewMockRepositoryHelper() *MockRepositoryHelper {
	return &MockRepositoryHelper{
		dbManager: NewMockDatabaseManager(),
	}
}

// GetDatabaseManager returns the database manager
func (h *MockRepositoryHelper) GetDatabaseManager() *MockDatabaseManager {
	return h.dbManager
}

// SetupUserCreation configures mocks for successful user creation
func (h *MockRepositoryHelper) SetupUserCreation(userID string) {
	// Mock successful user creation
	h.dbManager.mockTx.On("QueryRaw", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
}

// SetupUserCreationFailure configures mocks for failed user creation
func (h *MockRepositoryHelper) SetupUserCreationFailure(err error) {
	h.dbManager.mockTx.On("QueryRaw", mock.Anything, mock.Anything, mock.Anything).Return(err).Once()
}

// SetupTeamCreation configures mocks for successful team creation
func (h *MockRepositoryHelper) SetupTeamCreation(teamID string) {
	h.dbManager.mockTx.On("QueryRaw", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
}

// SetupUserLookup configures mocks for user lookup operations
func (h *MockRepositoryHelper) SetupUserLookup(found bool, err error) {
	if found {
		h.dbManager.mockDB.On("QueryRaw", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	} else if err != nil {
		h.dbManager.mockDB.On("QueryRaw", mock.Anything, mock.Anything, mock.Anything).Return(err)
	} else {
		h.dbManager.mockDB.On("QueryRaw", mock.Anything, mock.Anything, mock.Anything).Return(database.ErrNoRecords)
	}
}

// SetupTokenOperations configures mocks for token-related operations
func (h *MockRepositoryHelper) SetupTokenOperations(success bool, err error) {
	if success {
		h.dbManager.mockDB.On("QueryRaw", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	} else {
		h.dbManager.mockDB.On("QueryRaw", mock.Anything, mock.Anything, mock.Anything).Return(err)
	}
}

// DatabaseMockBuilder provides a fluent interface for building database mocks
type DatabaseMockBuilder struct {
	mockDB *MockDatabase
}

// NewDatabaseMockBuilder creates a new database mock builder
func NewDatabaseMockBuilder() *DatabaseMockBuilder {
	return &DatabaseMockBuilder{
		mockDB: NewMockDatabase(),
	}
}

// ExpectQuery sets up an expectation for a QueryRaw call
func (b *DatabaseMockBuilder) ExpectQuery(entity database.RawEntity, code int) *DatabaseMockBuilder {
	b.mockDB.On("QueryRaw", mock.Anything, entity, code).Return(nil)
	return b
}

// ExpectQueryWithError sets up an expectation for a QueryRaw call that returns an error
func (b *DatabaseMockBuilder) ExpectQueryWithError(entity database.RawEntity, code int, err error) *DatabaseMockBuilder {
	b.mockDB.On("QueryRaw", mock.Anything, entity, code).Return(err)
	return b
}

// ExpectMultiQuery sets up an expectation for a QueryMultiRaw call
func (b *DatabaseMockBuilder) ExpectMultiQuery(entity database.RawEntity, code int, results []database.RawEntity) *DatabaseMockBuilder {
	b.mockDB.On("QueryMultiRaw", mock.Anything, entity, code).Return(results, nil)
	return b
}

// ExpectMultiQueryWithError sets up an expectation for a QueryMultiRaw call that returns an error
func (b *DatabaseMockBuilder) ExpectMultiQueryWithError(entity database.RawEntity, code int, err error) *DatabaseMockBuilder {
	b.mockDB.On("QueryMultiRaw", mock.Anything, entity, code).Return([]database.RawEntity(nil), err)
	return b
}

// ExpectExec sets up an expectation for an ExecRaws call
func (b *DatabaseMockBuilder) ExpectExec(source string, execs ...database.RawExec) *DatabaseMockBuilder {
	b.mockDB.On("ExecRaws", mock.Anything, source, execs).Return(nil)
	return b
}

// ExpectExecWithError sets up an expectation for an ExecRaws call that returns an error
func (b *DatabaseMockBuilder) ExpectExecWithError(source string, err error, execs ...database.RawExec) *DatabaseMockBuilder {
	b.mockDB.On("ExecRaws", mock.Anything, source, execs).Return(err)
	return b
}

// Build returns the configured mock database
func (b *DatabaseMockBuilder) Build() *MockDatabase {
	return b.mockDB
}

// CallVerifier provides utilities for verifying mock calls
type CallVerifier struct {
	mockDB *MockDatabase
}

// NewCallVerifier creates a new call verifier
func NewCallVerifier(mockDB *MockDatabase) *CallVerifier {
	return &CallVerifier{mockDB: mockDB}
}

// VerifyQueryCalledWith verifies that QueryRaw was called with specific arguments
func (v *CallVerifier) VerifyQueryCalledWith(ctx context.Context, entity database.RawEntity, code int) error {
	for _, call := range v.mockDB.Calls {
		if call.Method == "QueryRaw" && len(call.Arguments) >= 3 {
			if call.Arguments[1] == entity && call.Arguments[2] == code {
				return nil
			}
		}
	}
	return fmt.Errorf("QueryRaw was not called with expected arguments: entity=%v, code=%d", entity, code)
}

// VerifyExecCalledWith verifies that ExecRaws was called with specific arguments
func (v *CallVerifier) VerifyExecCalledWith(ctx context.Context, source string, execs ...database.RawExec) error {
	for _, call := range v.mockDB.Calls {
		if call.Method == "ExecRaws" && len(call.Arguments) >= 3 {
			if call.Arguments[1] == source {
				return nil
			}
		}
	}
	return fmt.Errorf("ExecRaws was not called with expected arguments: source=%s", source)
}

// GetCallCount returns the number of times a method was called
func (v *CallVerifier) GetCallCount(methodName string) int {
	count := 0
	for _, call := range v.mockDB.Calls {
		if call.Method == methodName {
			count++
		}
	}
	return count
}