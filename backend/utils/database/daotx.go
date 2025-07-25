package database

import (
	"context"
	"database/sql"
	"errors"
)

type dbTxClient struct {
	*sql.Tx
}

func (d *dbTxClient) QueryRaw(ctx context.Context, entity RawEntity, code int) error {
	row := d.QueryRowContext(ctx, entity.GetQuery(code), entity.GetQueryValues(code)...)
	err := entity.BindRawRow(code, row)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNoRecords
	}
	return err
}

func (d *dbTxClient) QueryMultiRaw(ctx context.Context, entity RawEntity, code int) ([]RawEntity, error) {
	rows, err := d.QueryContext(ctx, entity.GetMultiQuery(code), entity.GetMultiQueryValues(code)...)
	return handleQueryMultiRawResponse(rows, err, entity, code)
}

func (d *dbTxClient) ExecRaws(ctx context.Context, source string, execs ...RawExec) error {
	for _, exec := range execs {
		_, err := d.ExecContext(ctx, exec.Entity.GetExec(exec.Code), exec.Entity.GetExecValues(exec.Code, source)...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *dbTxClient) ExecRawsConsistent(ctx context.Context, source string, execs ...RawExec) error {
	for _, exec := range execs {
		result, err := d.ExecContext(ctx, exec.Entity.GetExec(exec.Code), exec.Entity.GetExecValues(exec.Code, source)...)
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
	return nil
}
