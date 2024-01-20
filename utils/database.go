package utils

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/mladenovic-13/bank-api/api"
	"github.com/mladenovic-13/bank-api/internal/database"
	"github.com/mladenovic-13/bank-api/models"
)

type DepostitResult struct {
	UpdatedAccount *models.Account     `json:"account"`
	Transaction    *models.Transaction `json:"transaction"`
}

func ExecDepositTransaction(
	ctx context.Context,
	db *sql.DB,
	queries *database.Queries,
	account *database.Account,
	depositRequest *api.DepositRequest,
) (*WithdrawResult, error) {
	tx, err := db.Begin()

	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	qtx := queries.WithTx(tx)

	balanceNumber, err := strconv.ParseFloat(account.Balance, 32)

	newBalance := balanceNumber + depositRequest.Amount

	if err != nil {
		return nil, err
	}

	updatedAccount, err := qtx.UpdateAccountBalance(ctx,
		database.UpdateAccountBalanceParams{
			ID:        account.ID,
			Balance:   fmt.Sprintf("%0.2f", newBalance),
			UpdatedAt: time.Now(),
		},
	)

	if err != nil {
		return nil, err
	}

	transaction, err := qtx.CreateTransaction(
		ctx,
		database.CreateTransactionParams{
			ID:              uuid.New(),
			SenderNumber:    account.Number,
			ReceiverNumber:  account.Number,
			Amount:          fmt.Sprintf("%0.2f", depositRequest.Amount),
			Currency:        depositRequest.Currency,
			TransactionType: database.TransactionTypeDEPOSIT,
			CreatedAt:       time.Now(),
		},
	)

	if err != nil {
		return nil, err
	}

	result := &WithdrawResult{
		UpdatedAccount: models.ToAccount(updatedAccount),
		Transaction:    models.ToTransaction(transaction),
	}

	return result, tx.Commit()
}

type WithdrawResult struct {
	UpdatedAccount *models.Account     `json:"account"`
	Transaction    *models.Transaction `json:"transaction"`
}

func ExecWithdrawTransaction(
	ctx context.Context,
	db *sql.DB,
	queries *database.Queries,
	account *database.Account,
	withdrawRequest *api.DepositRequest,
) (*WithdrawResult, error) {
	tx, err := db.Begin()

	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	qtx := queries.WithTx(tx)

	balanceNumber, err := strconv.ParseFloat(account.Balance, 32)

	newBalance := balanceNumber - withdrawRequest.Amount

	if err != nil {
		return nil, err
	}

	updatedAccount, err := qtx.UpdateAccountBalance(ctx,
		database.UpdateAccountBalanceParams{
			ID:        account.ID,
			Balance:   fmt.Sprintf("%0.2f", newBalance),
			UpdatedAt: time.Now(),
		},
	)

	if err != nil {
		return nil, err
	}

	transaction, err := qtx.CreateTransaction(
		ctx,
		database.CreateTransactionParams{
			ID:              uuid.New(),
			SenderNumber:    account.Number,
			ReceiverNumber:  account.Number,
			Amount:          fmt.Sprintf("%0.2f", withdrawRequest.Amount),
			Currency:        withdrawRequest.Currency,
			TransactionType: database.TransactionTypeWITHDRAWAL,
			CreatedAt:       time.Now(),
		},
	)

	if err != nil {
		return nil, err
	}

	result := &WithdrawResult{
		UpdatedAccount: models.ToAccount(updatedAccount),
		Transaction:    models.ToTransaction(transaction),
	}

	return result, tx.Commit()
}

type SendTransactionResult struct {
	SenderAccount   *models.Account     `json:"senderAccount"`
	ReceiverAccount *models.Account     `json:"receiverAccount"`
	Transaction     *models.Transaction `json:"transaction"`
}

func ExecSendTransaction(
	ctx context.Context,
	db *sql.DB,
	queries *database.Queries,
	senderAccount *database.Account,
	receiverAccount *database.Account,
	sendRequest *api.SendRequest,
) (*SendTransactionResult, error) {
	tx, err := db.Begin()

	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	qtx := queries.WithTx(tx)

	senderBalance, err := strconv.ParseFloat(senderAccount.Balance, 32)
	if err != nil {
		return nil, err
	}

	receiverBalance, err := strconv.ParseFloat(receiverAccount.Balance, 32)
	if err != nil {
		return nil, err
	}

	updatedSenderAccount, err := qtx.UpdateAccountBalance(
		ctx,
		database.UpdateAccountBalanceParams{
			ID:        senderAccount.ID,
			Balance:   fmt.Sprintf("%0.2f", senderBalance-sendRequest.Amount),
			UpdatedAt: time.Now(),
		},
	)

	if err != nil {
		return nil, err
	}

	updatedReceiverAccount, err := qtx.UpdateAccountBalance(
		ctx,
		database.UpdateAccountBalanceParams{
			ID:        receiverAccount.ID,
			Balance:   fmt.Sprintf("%0.2f", receiverBalance+sendRequest.Amount),
			UpdatedAt: time.Now(),
		},
	)

	if err != nil {
		return nil, err
	}

	transaction, err := qtx.CreateTransaction(
		ctx,
		database.CreateTransactionParams{
			ID:              uuid.New(),
			SenderNumber:    senderAccount.Number,
			ReceiverNumber:  receiverAccount.Number,
			Amount:          fmt.Sprintf("%0.2f", sendRequest.Amount),
			Currency:        sendRequest.Currency,
			TransactionType: database.TransactionTypeTRANSFER,
			CreatedAt:       time.Now(),
		},
	)

	if err != nil {
		return nil, err
	}

	result := &SendTransactionResult{
		SenderAccount:   models.ToAccount(updatedSenderAccount),
		ReceiverAccount: models.ToAccount(updatedReceiverAccount),
		Transaction:     models.ToTransaction(transaction),
	}

	log.Println(models.ToTransaction(transaction).ToString())

	return result, tx.Commit()
}
