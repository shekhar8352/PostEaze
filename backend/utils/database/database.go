package database

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Config struct {
	DriverName            string        `json:"driverName"`
	URL                   string        `json:"url"`
	MaxOpenConnections    int           `json:"maxOpenConnections"`
	MaxIdleConnections    int           `json:"maxIdleConnections"`
	ConnectionMaxLifetime time.Duration `json:"connectionMaxLifetime"`
	ConnectionMaxIdleTime time.Duration `json:"connectionMaxIdleTime"`
}

// Database is the set of methods available for the database connection.
type Database interface {
	QueryRaw(ctx context.Context, entity RawEntity, code int) error
	QueryMultiRaw(ctx context.Context, entity RawEntity, code int) ([]RawEntity, error)
	ExecRaws(ctx context.Context, source string, execs ...RawExec) error
	ExecRawsConsistent(ctx context.Context, source string, execs ...RawExec) error
}

// common errors
var (
	ErrNoRecords      = errors.New("no records found")
	ErrNoRowsAffected = errors.New("no rows affected")
)

var db *dbClient

// Init is used to initialise the database client.
func Init(ctx context.Context, config Config) error {
	conn, err := sql.Open(config.DriverName, config.URL)
	if err != nil {
		return err
	}
	err = conn.PingContext(ctx)
	if err != nil {
		_ = conn.Close()
		return err
	}
	conn.SetMaxOpenConns(config.MaxOpenConnections)
	conn.SetMaxIdleConns(config.MaxIdleConnections)
	conn.SetConnMaxLifetime(config.ConnectionMaxLifetime)
	conn.SetConnMaxIdleTime(config.ConnectionMaxIdleTime)
	db = &dbClient{conn}
	return nil
}

// Close is used to close the database instance.
func Close() error {
	return db.Close()
}

// Ping is used to check the connectivity to the database instance.
func Ping(ctx context.Context) error {
	return db.PingContext(ctx)
}

// Get is used to get the database instance.
func Get() Database {
	return db
}

// GetTx is used to get the transactional database instance.
func GetTx(ctx context.Context, options *sql.TxOptions) (Database, error) {
	tx, err := db.BeginTx(ctx, options)
	if err != nil {
		return nil, err
	}
	return &dbTxClient{tx}, nil
}

// CommitTx is used to commit the transaction.
func CommitTx(db Database) error {
	if tx, ok := db.(*dbTxClient); ok {
		return tx.Commit()
	}
	return nil
}

// RollbackTx is used to roll back the transaction.
func RollbackTx(db Database) {
	if tx, ok := db.(*dbTxClient); ok {
		_ = tx.Rollback()
	}
}
