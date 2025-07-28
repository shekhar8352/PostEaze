package database

import (
	"context"
	"database/sql"
	"errors"
)

type dbClient struct {
	*sql.DB
}

func (d *dbClient) QueryRaw(ctx context.Context, entity RawEntity, code int) error {
	row := d.QueryRowContext(ctx, entity.GetQuery(code), entity.GetQueryValues(code)...)
	err := entity.BindRawRow(code, row)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNoRecords
	}
	return err
}

func (d *dbClient) QueryMultiRaw(ctx context.Context, entity RawEntity, code int) ([]RawEntity, error) {
	rows, err := d.QueryContext(ctx, entity.GetMultiQuery(code), entity.GetMultiQueryValues(code)...)
	return handleQueryMultiRawResponse(rows, err, entity, code)
}

func (d *dbClient) ExecRaws(ctx context.Context, source string, execs ...RawExec) error {
	tx, err := d.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer rollback(tx)
	for _, exec := range execs {
		_, err = tx.ExecContext(ctx, exec.Entity.GetExec(exec.Code), exec.Entity.GetExecValues(exec.Code, source)...)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (d *dbClient) ExecRawsConsistent(ctx context.Context, source string, execs ...RawExec) error {
	tx, err := d.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer rollback(tx)
	for _, exec := range execs {
		result, err := tx.ExecContext(ctx, exec.Entity.GetExec(exec.Code), exec.Entity.GetExecValues(exec.Code, source)...)
		if err != nil {
			return err
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return ErrNoRowsAffected
		}
	}
	return tx.Commit()
}

func rollback(tx *sql.Tx) {
	if tx != nil {
		_ = tx.Rollback()
	}
}

func closeRows(rows *sql.Rows) {
	if rows != nil {
		_ = rows.Close()
	}
}

func handleQueryMultiRawResponse(rows *sql.Rows, err error, entity RawEntity, code int) ([]RawEntity, error) {
	defer closeRows(rows)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNoRecords
	}
	if err != nil {
		return nil, err
	}
	result := make([]RawEntity, 0)
	for rows.Next() {
		err = entity.BindRawRow(code, rows)
		if err != nil {
			return nil, err
		}
		result = append(result, entity)
		entity = entity.GetNextRaw()
	}
	if len(result) == 0 {
		return nil, ErrNoRecords
	}
	return result, nil
}
