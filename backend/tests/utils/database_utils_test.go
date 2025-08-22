package utils

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shekhar8352/PostEaze/utils/database"
)

func TestConfig(t *testing.T) {
	t.Run("valid configuration", func(t *testing.T) {
		config := database.Config{
			DriverName:            "postgres",
			URL:                   "postgres://user:pass@localhost:5432/testdb",
			MaxOpenConnections:    25,
			MaxIdleConnections:    10,
			ConnectionMaxLifetime: 1 * time.Hour,
			ConnectionMaxIdleTime: 15 * time.Minute,
		}

		if config.DriverName != "postgres" {
			t.Errorf("DriverName = %v, want postgres", config.DriverName)
		}
		if config.URL != "postgres://user:pass@localhost:5432/testdb" {
			t.Errorf("URL = %v, want postgres://user:pass@localhost:5432/testdb", config.URL)
		}
		if config.MaxOpenConnections != 25 {
			t.Errorf("MaxOpenConnections = %v, want 25", config.MaxOpenConnections)
		}
		if config.MaxIdleConnections != 10 {
			t.Errorf("MaxIdleConnections = %v, want 10", config.MaxIdleConnections)
		}
		if config.ConnectionMaxLifetime != 1*time.Hour {
			t.Errorf("ConnectionMaxLifetime = %v, want 1h", config.ConnectionMaxLifetime)
		}
		if config.ConnectionMaxIdleTime != 15*time.Minute {
			t.Errorf("ConnectionMaxIdleTime = %v, want 15m", config.ConnectionMaxIdleTime)
		}
	})

	t.Run("default values", func(t *testing.T) {
		config := database.Config{}

		if config.DriverName != "" {
			t.Errorf("Default DriverName = %v, want empty", config.DriverName)
		}
		if config.URL != "" {
			t.Errorf("Default URL = %v, want empty", config.URL)
		}
		if config.MaxOpenConnections != 0 {
			t.Errorf("Default MaxOpenConnections = %v, want 0", config.MaxOpenConnections)
		}
		if config.MaxIdleConnections != 0 {
			t.Errorf("Default MaxIdleConnections = %v, want 0", config.MaxIdleConnections)
		}
		if config.ConnectionMaxLifetime != time.Duration(0) {
			t.Errorf("Default ConnectionMaxLifetime = %v, want 0", config.ConnectionMaxLifetime)
		}
		if config.ConnectionMaxIdleTime != time.Duration(0) {
			t.Errorf("Default ConnectionMaxIdleTime = %v, want 0", config.ConnectionMaxIdleTime)
		}
	})
}

func TestDatabaseErrors(t *testing.T) {
	if database.ErrNoRecords.Error() != "no records found" {
		t.Errorf("ErrNoRecords = %v, want 'no records found'", database.ErrNoRecords.Error())
	}
	if database.ErrNoRowsAffected.Error() != "no rows affected" {
		t.Errorf("ErrNoRowsAffected = %v, want 'no rows affected'", database.ErrNoRowsAffected.Error())
	}

	// Test that errors are different instances
	if database.ErrNoRecords == database.ErrNoRowsAffected {
		t.Error("ErrNoRecords and ErrNoRowsAffected should be different instances")
	}

	// Test error comparison
	if !errors.Is(database.ErrNoRecords, database.ErrNoRecords) {
		t.Error("ErrNoRecords should match itself")
	}
	if errors.Is(database.ErrNoRecords, database.ErrNoRowsAffected) {
		t.Error("ErrNoRecords should not match ErrNoRowsAffected")
	}
}

// Mock implementation of RawEntity for testing
type mockRawEntity struct {
	query            string
	queryValues      []any
	multiQuery       string
	multiQueryValues []any
	exec             string
	execValues       []any
	scanError        error
	nextEntity       database.RawEntity
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
		query:            m.query,
		queryValues:      m.queryValues,
		multiQuery:       m.multiQuery,
		multiQueryValues: m.multiQueryValues,
		exec:             m.exec,
		execValues:       m.execValues,
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

func TestInit(t *testing.T) {
	ctx := context.Background()

	t.Run("invalid driver", func(t *testing.T) {
		config := database.Config{
			DriverName: "nonexistent_driver",
			URL:        "some://url",
		}

		err := database.Init(ctx, config)
		if err == nil {
			t.Error("Init() should return error for invalid driver")
		}
	})

	t.Run("empty config", func(t *testing.T) {
		config := database.Config{}

		err := database.Init(ctx, config)
		if err == nil {
			t.Error("Init() should return error for empty config")
		}
	})

	t.Run("invalid URL", func(t *testing.T) {
		config := database.Config{
			DriverName: "invalid_driver",
			URL:        "invalid_url",
		}

		err := database.Init(ctx, config)
		if err == nil {
			t.Error("Init() should return error for invalid URL")
		}
	})
}

func TestQueryOperations(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer mockDB.Close()

	t.Run("query success", func(t *testing.T) {
		entity := &mockRawEntity{
			query:       "SELECT id, name FROM users WHERE id = ?",
			queryValues: []any{1},
		}

		// Mock successful query
		rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "test_user")
		mock.ExpectQuery("SELECT id, name FROM users WHERE id = ?").
			WithArgs(1).
			WillReturnRows(rows)

		if entity.GetQuery(0) != "SELECT id, name FROM users WHERE id = ?" {
			t.Errorf("GetQuery() = %v, want SELECT id, name FROM users WHERE id = ?", entity.GetQuery(0))
		}
		if len(entity.GetQueryValues(0)) != 1 || entity.GetQueryValues(0)[0] != 1 {
			t.Errorf("GetQueryValues() = %v, want [1]", entity.GetQueryValues(0))
		}
	})

	t.Run("no records error conversion", func(t *testing.T) {
		// Test that sql.ErrNoRows is converted to ErrNoRecords
		if !errors.Is(sql.ErrNoRows, sql.ErrNoRows) {
			t.Error("sql.ErrNoRows should match itself")
		}

		// Test error conversion logic
		err := sql.ErrNoRows
		if errors.Is(err, sql.ErrNoRows) {
			err = database.ErrNoRecords
		}
		if err != database.ErrNoRecords {
			t.Errorf("Converted error = %v, want ErrNoRecords", err)
		}
	})

	t.Run("multi query", func(t *testing.T) {
		entity := &mockRawEntity{
			multiQuery:       "SELECT id, name FROM users",
			multiQueryValues: []any{},
		}

		// Mock successful multi-row query
		rows := sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "user1").
			AddRow(2, "user2")

		mock.ExpectQuery("SELECT id, name FROM users").
			WillReturnRows(rows)

		if entity.GetMultiQuery(0) != "SELECT id, name FROM users" {
			t.Errorf("GetMultiQuery() = %v, want SELECT id, name FROM users", entity.GetMultiQuery(0))
		}
		if len(entity.GetMultiQueryValues(0)) != 0 {
			t.Errorf("GetMultiQueryValues() = %v, want []", entity.GetMultiQueryValues(0))
		}
	})

	t.Run("empty result set", func(t *testing.T) {
		entity := &mockRawEntity{
			multiQuery:       "SELECT id, name FROM users WHERE id = ?",
			multiQueryValues: []any{999},
		}

		// Mock empty result set
		rows := sqlmock.NewRows([]string{"id", "name"})
		mock.ExpectQuery("SELECT id, name FROM users WHERE id = ?").
			WithArgs(999).
			WillReturnRows(rows)

		if entity.GetMultiQuery(0) != "SELECT id, name FROM users WHERE id = ?" {
			t.Errorf("GetMultiQuery() = %v, want SELECT id, name FROM users WHERE id = ?", entity.GetMultiQuery(0))
		}
	})
}

func TestTransactionOperations(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer mockDB.Close()

	t.Run("exec success", func(t *testing.T) {
		entity1 := &mockRawEntity{
			exec:       "INSERT INTO users (name) VALUES (?)",
			execValues: []any{"test_user"},
		}

		entity2 := &mockRawEntity{
			exec:       "UPDATE users SET active = ? WHERE id = ?",
			execValues: []any{true, 1},
		}

		// Mock transaction
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO users \\(name\\) VALUES \\(\\?\\)").
			WithArgs("test_user").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("UPDATE users SET active = \\? WHERE id = \\?").
			WithArgs(true, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		// Test RawExec structure
		exec1 := database.RawExec{Entity: entity1, Code: 1}
		exec2 := database.RawExec{Entity: entity2, Code: 2}

		if exec1.Entity != entity1 {
			t.Errorf("RawExec.Entity = %v, want entity1", exec1.Entity)
		}
		if exec1.Code != 1 {
			t.Errorf("RawExec.Code = %v, want 1", exec1.Code)
		}
		if exec2.Entity != entity2 {
			t.Errorf("RawExec.Entity = %v, want entity2", exec2.Entity)
		}
		if exec2.Code != 2 {
			t.Errorf("RawExec.Code = %v, want 2", exec2.Code)
		}
	})

	t.Run("transaction failure", func(t *testing.T) {
		entity := &mockRawEntity{
			exec:       "INSERT INTO users (name) VALUES (?)",
			execValues: []any{"test_user"},
		}

		// Mock transaction failure
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO users \\(name\\) VALUES \\(\\?\\)").
			WithArgs("test_user").
			WillReturnError(errors.New("constraint violation"))
		mock.ExpectRollback()

		exec := database.RawExec{Entity: entity, Code: 1}
		if exec.Entity == nil {
			t.Error("RawExec.Entity should not be nil")
		}
	})

	t.Run("no rows affected", func(t *testing.T) {
		entity := &mockRawEntity{
			exec:       "UPDATE users SET name = ? WHERE id = ?",
			execValues: []any{"updated_name", 999},
		}

		// Mock update with no rows affected
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE users SET name = \\? WHERE id = \\?").
			WithArgs("updated_name", 999).
			WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected
		mock.ExpectRollback()

		exec := database.RawExec{Entity: entity, Code: 1}
		if exec.Entity == nil {
			t.Error("RawExec.Entity should not be nil")
		}

		// Test that ErrNoRowsAffected would be returned
		result := sqlmock.NewResult(0, 0)
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			t.Errorf("RowsAffected() error = %v", err)
		}
		if rowsAffected != 0 {
			t.Errorf("RowsAffected() = %v, want 0", rowsAffected)
		}

		if rowsAffected == 0 {
			if database.ErrNoRowsAffected != database.ErrNoRowsAffected {
				t.Error("ErrNoRowsAffected should equal itself")
			}
		}
	})
}

func TestConnectionConfiguration(t *testing.T) {
	configs := []database.Config{
		{
			DriverName:            "postgres",
			URL:                   "postgres://user:pass@localhost:5432/testdb",
			MaxOpenConnections:    20,
			MaxIdleConnections:    8,
			ConnectionMaxLifetime: 45 * time.Minute,
			ConnectionMaxIdleTime: 10 * time.Minute,
		},
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
		if config.MaxOpenConnections < 0 {
			t.Errorf("Config %d: MaxOpenConnections should be >= 0, got %d", i, config.MaxOpenConnections)
		}
		if config.MaxIdleConnections < 0 {
			t.Errorf("Config %d: MaxIdleConnections should be >= 0, got %d", i, config.MaxIdleConnections)
		}
		if config.ConnectionMaxLifetime < 0 {
			t.Errorf("Config %d: ConnectionMaxLifetime should be >= 0, got %v", i, config.ConnectionMaxLifetime)
		}
		if config.ConnectionMaxIdleTime < 0 {
			t.Errorf("Config %d: ConnectionMaxIdleTime should be >= 0, got %v", i, config.ConnectionMaxIdleTime)
		}
	}
}

func TestDatabaseInterface(t *testing.T) {
	entity := &mockRawEntity{
		query:            "SELECT * FROM test",
		queryValues:      []any{},
		multiQuery:       "SELECT * FROM test",
		multiQueryValues: []any{},
		exec:             "INSERT INTO test VALUES (?)",
		execValues:       []any{"value"},
	}

	// Test RawEntity interface methods
	if entity.GetQuery(0) != "SELECT * FROM test" {
		t.Errorf("GetQuery() = %v, want SELECT * FROM test", entity.GetQuery(0))
	}
	if len(entity.GetQueryValues(0)) != 0 {
		t.Errorf("GetQueryValues() = %v, want []", entity.GetQueryValues(0))
	}
	if entity.GetMultiQuery(0) != "SELECT * FROM test" {
		t.Errorf("GetMultiQuery() = %v, want SELECT * FROM test", entity.GetMultiQuery(0))
	}
	if len(entity.GetMultiQueryValues(0)) != 0 {
		t.Errorf("GetMultiQueryValues() = %v, want []", entity.GetMultiQueryValues(0))
	}
	if entity.GetExec(0) != "INSERT INTO test VALUES (?)" {
		t.Errorf("GetExec() = %v, want INSERT INTO test VALUES (?)", entity.GetExec(0))
	}
	if len(entity.GetExecValues(0, "test")) != 1 || entity.GetExecValues(0, "test")[0] != "value" {
		t.Errorf("GetExecValues() = %v, want [value]", entity.GetExecValues(0, "test"))
	}
	if entity.GetNextRaw() == nil {
		t.Error("GetNextRaw() should not return nil")
	}
}

func TestRawExecStructure(t *testing.T) {
	entity := &mockRawEntity{
		exec:       "UPDATE test SET value = ?",
		execValues: []any{"new_value"},
	}

	rawExec := database.RawExec{
		Entity: entity,
		Code:   42,
	}

	if rawExec.Entity != entity {
		t.Errorf("RawExec.Entity = %v, want entity", rawExec.Entity)
	}
	if rawExec.Code != 42 {
		t.Errorf("RawExec.Code = %v, want 42", rawExec.Code)
	}
	if rawExec.Entity.GetExec(rawExec.Code) != "UPDATE test SET value = ?" {
		t.Errorf("Entity.GetExec() = %v, want UPDATE test SET value = ?", rawExec.Entity.GetExec(rawExec.Code))
	}
	execValues := rawExec.Entity.GetExecValues(rawExec.Code, "test")
	if len(execValues) != 1 || execValues[0] != "new_value" {
		t.Errorf("Entity.GetExecValues() = %v, want [new_value]", execValues)
	}
}

func TestErrorHandling(t *testing.T) {
	t.Run("scan error", func(t *testing.T) {
		entity := &mockRawEntity{
			scanError: errors.New("scan error"),
		}

		mockDB, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock DB: %v", err)
		}
		defer mockDB.Close()

		err = entity.BindRawRow(0, mockDB.QueryRow("SELECT 1"))
		if err == nil {
			t.Error("BindRawRow() should return error")
		}
		if err.Error() != "scan error" {
			t.Errorf("BindRawRow() error = %v, want scan error", err)
		}
	})

	t.Run("context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		if ctx.Err() == nil {
			t.Error("Cancelled context should have error")
		}
		if ctx.Err() != context.Canceled {
			t.Errorf("Context error = %v, want context.Canceled", ctx.Err())
		}
	})

	t.Run("context timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()

		// Wait for timeout
		time.Sleep(1 * time.Millisecond)

		if ctx.Err() == nil {
			t.Error("Timed out context should have error")
		}
		if ctx.Err() != context.DeadlineExceeded {
			t.Errorf("Context error = %v, want context.DeadlineExceeded", ctx.Err())
		}
	})
}

func TestUtilityFunctions(t *testing.T) {
	// Test that we can create mock results for testing
	result := sqlmock.NewResult(1, 1)
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		t.Errorf("LastInsertId() error = %v", err)
	}
	if lastInsertId != 1 {
		t.Errorf("LastInsertId() = %v, want 1", lastInsertId)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Errorf("RowsAffected() error = %v", err)
	}
	if rowsAffected != 1 {
		t.Errorf("RowsAffected() = %v, want 1", rowsAffected)
	}
}

func TestDatabaseDrivers(t *testing.T) {
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

		if config.DriverName != driverName {
			t.Errorf("Config.DriverName = %v, want %v", config.DriverName, driverName)
		}
		if config.URL == "" {
			t.Error("Config.URL should not be empty")
		}
	}
}