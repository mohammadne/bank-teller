package repository

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/mohammadne/bank-teller/inernal/entities"
)

type Bank interface {
	Transfer(ctx context.Context, from, to entities.Sheba, amount int) (*entities.Transaction, error)
	ListPendings(_ context.Context) (pool, error)
	ListConfirmed(_ context.Context) (pool, error)
	ListCanceled(_ context.Context) (pool, error)
	MoveTransaction(_ context.Context, transactionID string, status entities.TransactionStatus) (*entities.Transaction, error)
}

func NewBank(initialUsers []entities.User) Bank {
	return &bank{
		Users:     initialUsers,
		Pendings:  make(pool, 10),
		Confirmed: make(pool, 100),
		Canceled:  make(pool, 5),
	}
}

type bank struct {
	Users     []entities.User
	Locker    sync.RWMutex
	Pendings  pool
	Confirmed pool
	Canceled  pool
}

type pool []entities.Transaction

var (
	ErrSourceOrDestinationUsersNotFound = errors.New("error source or destination users not found")
	ErrNotEnoughBalance                 = errors.New("error not enough balance")
)

func (b *bank) Transfer(_ context.Context, from, to entities.Sheba, amount int) (*entities.Transaction, error) {
	var fromUser, toUser entities.User
	for _, user := range b.Users {
		if user.Sheba == from {
			fromUser = user
		} else if user.Sheba == to {
			toUser = user
		}
	}

	if fromUser.ID == 0 || toUser.ID == 0 {
		return nil, ErrSourceOrDestinationUsersNotFound
	}

	b.Locker.Lock()
	defer b.Locker.Unlock()

	if fromUser.Balance < amount {
		return nil, ErrNotEnoughBalance
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

	b.Pendings = append(b.Pendings, transaction)
	fromUser.Balance -= amount

	return &transaction, nil
}

func (b *bank) ListPendings(_ context.Context) (pool, error) {
	return b.Pendings, nil
}

func (b *bank) ListConfirmed(_ context.Context) (pool, error) {
	return b.Confirmed, nil
}

func (b *bank) ListCanceled(_ context.Context) (pool, error) {
	return b.Canceled, nil
}

var (
	ErrInvalidMoveTransactionStatus = errors.New("ErrInvalidMovetransactionStatus")
	ErrMoveTransactionNotFound      = errors.New("ErrMoveTransactionNotFound")
	ErrMoveTransactionUserNotFound  = errors.New("ErrMoveTransactionUserNotFound")
)

func (b *bank) MoveTransaction(_ context.Context, transactionID string, status entities.TransactionStatus) (*entities.Transaction, error) {
	if status == entities.TransactionStatusPending {
		return nil, ErrInvalidMoveTransactionStatus
	}

	var targetTransaction entities.Transaction
	var targetTransactionIndex = -1
	for index, transaction := range b.Pendings {
		if transaction.ID == transactionID {
			targetTransaction = transaction
			targetTransactionIndex = index
			break
		}
	}

	if targetTransactionIndex == -1 {
		return nil, ErrMoveTransactionNotFound
	}

	b.Locker.Lock()
	defer b.Locker.Unlock()

	var shebaToMatch entities.Sheba
	var addToPool func(entities.Transaction)

	if status == entities.TransactionStatusConfirmed {
		shebaToMatch = targetTransaction.To
		addToPool = func(t entities.Transaction) { b.Confirmed = append(b.Confirmed, t) }
	} else {
		shebaToMatch = targetTransaction.From
		addToPool = func(t entities.Transaction) { b.Canceled = append(b.Canceled, t) }
	}

	userIndex := -1
	for index, user := range b.Users {
		if user.Sheba == shebaToMatch {
			userIndex = index
			break
		}
	}
	if userIndex == -1 {
		return nil, ErrMoveTransactionUserNotFound
	}

	b.Users[userIndex].Balance += targetTransaction.Amount
	targetTransaction.Status = status
	b.Pendings = append(b.Pendings[:targetTransactionIndex], b.Pendings[targetTransactionIndex+1:]...)
	addToPool(targetTransaction)

	return &targetTransaction, nil
}
