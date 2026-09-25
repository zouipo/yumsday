package repository

import (
	"database/sql"
	"log/slog"
)

func RollbackTx(tx *sql.Tx, txDesc string) {
	if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
		slog.Error("failed to rollback transation", "error", err, "transaction", txDesc)
	}
}

func CloseRows(rows *sql.Rows) {
	if err := rows.Close(); err != nil {
		slog.Warn("failed to close rows", "error", err)
	}
}
