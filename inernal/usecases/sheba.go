package usecases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mohammadne/bank-teller/inernal/entities"
	"github.com/mohammadne/bank-teller/inernal/repository"
)

type Sheba interface {
	ListTransactions(_ context.Context, status string) (entities.Pool, error)
	Transfer(ctx context.Context, from, to entities.Sheba, amount int) (*entities.Transaction, error)
	MoveTransaction(ctx context.Context, transactionID string, status entities.TransactionStatus) (*entities.Transaction, error)
}

func NewSheba(b repository.Bank) Sheba {
	return &sheba{bank: b}
}

type sheba struct {
	bank repository.Bank
}

// ListTransactions implements Sheba.
func (s *sheba) ListTransactions(ctx context.Context, statusRaw string) (entities.Pool, error) {
	status, err := entities.ToTransactionStatus(statusRaw)
	if err != nil {
		return nil, err
	}

	return s.bank.ListTransactions(ctx, status), nil
}

var (
	ErrTransferRetrieveFromUser = errors.New("ErrTransferRetrieveFromUser")
	ErrTransferNotEnoughBalance = errors.New("ErrTransferNotEnoughBalance")
	ErrTransferRetrieveToUser   = errors.New("ErrTransferRetrieveToUser")
)

// Transfer implements Sheba.
func (s *sheba) Transfer(ctx context.Context, from, to entities.Sheba, amount int) (*entities.Transaction, error) {
	fromUser, err := s.bank.RetrieveUserBySheba(ctx, from)
	if err != nil {
		return nil, errors.Join(ErrTransferRetrieveFromUser, err)
	}

	if fromUser.Balance < amount {
		return nil, ErrTransferNotEnoughBalance
	}

	_, err = s.bank.RetrieveUserBySheba(ctx, to)
	if err != nil {
		// just to check validity of destination
		return nil, errors.Join(ErrTransferRetrieveToUser, err)
	}

	fromUser.Balance -= amount
	err = s.bank.UpdateUser(ctx, fromUser)
	if err != nil {
		return nil, err
	}

	createdAt := time.Now()
	transaction := entities.Transaction{
		ID:        fmt.Sprintf("%d-%s", createdAt.Unix(), uuid.New().String()),
		Status:    entities.TransactionStatusPending,
		From:      from,
		To:        to,
		Amount:    amount,
		CreatedAt: createdAt,
	}
	s.bank.CreateTransaction(ctx, transaction)

	return &transaction, nil
}

var (
	ErrInvalidMoveTransactionStatus1 = errors.New("ErrInvalidMovetransactionStatus1")
	ErrInvalidMoveTransactionStatus2 = errors.New("ErrInvalidMovetransactionStatus2")
	ErrRetrieveTransaction           = errors.New("ErrRetrieveTransaction")
)

// MoveTransaction implements Sheba.
func (s *sheba) MoveTransaction(ctx context.Context, transactionID string, status entities.TransactionStatus) (*entities.Transaction, error) {
	if status != entities.TransactionStatusConfirmed && status != entities.TransactionStatusCanceled {
		return nil, ErrInvalidMoveTransactionStatus1
	}

	transaction, err := s.bank.RetrieveTransactionByID(ctx, transactionID)
	if err != nil {
		return nil, errors.Join(ErrRetrieveTransaction, err)
	}
	if transaction.Status != entities.TransactionStatusPending {
		return nil, ErrInvalidMoveTransactionStatus2
	}

	transaction.Status = status
	err = s.bank.UpdateTransaction(ctx, transaction)
	if err != nil {
		return nil, err
	}

	var user entities.User
	switch status {
	case entities.TransactionStatusConfirmed:
		user, err = s.bank.RetrieveUserBySheba(ctx, transaction.To)
	case entities.TransactionStatusCanceled:
		user, err = s.bank.RetrieveUserBySheba(ctx, transaction.From)
	}

	if err != nil {
		return nil, err
		// TODO: rollback transaction to pending
	}

	err = s.bank.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
		// TODO: rollback transaction to pending
	}

	return &transaction, nil
}
