package database

import "github.com/jmoiron/sqlx"

func CloseDBTransaction(resp *sqlx.Rows, err error) error {
	if err != nil {
		return err
	}

	if err := resp.Close(); err != nil {
		return err
	}

	return err
}

func CloseDBTransactionReturningID(resp *sqlx.Rows, err error) (int64, error) {
	if err != nil {
		return 0, err
	}

	var id int64
	resp.Next()

	if err := resp.Scan(&id); err != nil {
		return 0, err
	}

	if err := resp.Close(); err != nil {
		return 0, err
	}

	return id, err
}
