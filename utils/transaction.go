package utils

import (
	"context"
	"database/sql"

	"github.com/mladenovic-13/bank-api/internal/database"
	"github.com/mladenovic-13/bank-api/models"
)

func DepositTransaction(
	ctx context.Context,
	db *sql.DB,
	queries *database.Queries,
	accountParams database.UpdateAccountBalanceParams,
	transactionParams database.CreateTransactionParams,
) (*models.Account, *models.Transaction, error) {
	tx, err := db.Begin()

	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback()

	qtx := queries.WithTx(tx)

	account, err := qtx.UpdateAccountBalance(ctx, accountParams)

	if err != nil {
		return nil, nil, err

	}

	transaction, err := qtx.CreateTransaction(ctx, transactionParams)

	if err != nil {
		return nil, nil, err

	}

	return models.ToAccount(account),
		models.ToTransaction(transaction),
		tx.Commit()
}
