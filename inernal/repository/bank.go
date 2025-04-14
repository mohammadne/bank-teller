package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/mohammadne/bank-teller/inernal/entities"
)

type Bank interface {
	// users
	CreateUser(context.Context, entities.User)
	RetrieveUserBySheba(_ context.Context, sheba entities.Sheba) (entities.User, error)
	RetrieveUserByID(_ context.Context, id int) (entities.User, error)
	ListUsers(context.Context) []entities.User
	UpdateUser(context.Context, entities.User) error

	// transactions
	CreateTransaction(context.Context, entities.Transaction)
	RetrieveTransactionByID(_ context.Context, id string) (entities.Transaction, error)
	ListTransactions(_ context.Context, status entities.TransactionStatus) entities.Pool
	UpdateTransaction(context.Context, entities.Transaction) error
}

func NewBank(initialUsers []entities.User) Bank {
	return &bank{
		users:        initialUsers,
		transactions: make(entities.Pool, 0, 100),
	}
}

type bank struct {
	users      []entities.User
	userLocker sync.RWMutex

	transactions      entities.Pool
	transactionLocker sync.RWMutex
}

func (b *bank) CreateUser(_ context.Context, user entities.User) {
	b.users = append(b.users, user)
}

var ErrUserNotFound = errors.New("ErrUserNotFound")

func (b *bank) RetrieveUserBySheba(_ context.Context, sheba entities.Sheba) (entities.User, error) {
	b.userLocker.RLock()
	defer b.userLocker.RUnlock()

	for _, user := range b.users {
		if user.Sheba == sheba {
			return user, nil
		}
	}
	return entities.User{}, ErrUserNotFound
}

func (b *bank) RetrieveUserByID(_ context.Context, id int) (entities.User, error) {
	b.userLocker.RLock()
	defer b.userLocker.RUnlock()

	for _, user := range b.users {
		if user.ID == id {
			return user, nil
		}
	}
	return entities.User{}, ErrUserNotFound
}

func (b *bank) ListUsers(_ context.Context) []entities.User {
	b.userLocker.RLock()
	defer b.userLocker.RUnlock()

	result := make([]entities.User, len(b.users))
	copy(result, b.users)

	return result
}

var (
	ErrInvalidUser      = errors.New("ErrInvalidUser")
	ErrInvalidUserSheba = errors.New("ErrInvalidUserSheba")
)

func (b *bank) UpdateUser(_ context.Context, updatedUser entities.User) error {
	if updatedUser.ID <= 0 {
		return ErrInvalidUser
	}

	for index, user := range b.users {
		if user.ID == updatedUser.ID {
			if !user.CheckImmutable(&updatedUser) {
				return ErrInvalidUserSheba
			}

			b.userLocker.Lock()
			b.users[index] = updatedUser
			b.userLocker.Unlock()
			break
		}
	}

	return nil
}

var (
	ErrSourceOrDestinationUsersNotFound = errors.New("error source or destination users not found")
	ErrNotEnoughBalance                 = errors.New("error not enough balance")
)

func (b *bank) CreateTransaction(_ context.Context, t entities.Transaction) {
	b.transactions = append(b.transactions, t)
}

var ErrTransactionNotFound = errors.New("ErrTransactionNotFound")

func (b *bank) RetrieveTransactionByID(_ context.Context, id string) (entities.Transaction, error) {
	b.transactionLocker.RLock()
	defer b.transactionLocker.RUnlock()

	for _, transaction := range b.transactions {
		if transaction.ID == id {
			return transaction, nil
		}
	}
	return entities.Transaction{}, ErrTransactionNotFound
}

func (b *bank) ListTransactions(_ context.Context, status entities.TransactionStatus) entities.Pool {
	b.transactionLocker.RLock()
	defer b.transactionLocker.RUnlock()

	result := make(entities.Pool, 0, len(b.transactions))
	for _, transaction := range b.transactions {
		if transaction.Status == status {
			result = append(result, transaction)
		}
	}
	return result
}

var (
	ErrInvalidTransaction = errors.New("ErrInvalidTransaction")
)

func (b *bank) UpdateTransaction(_ context.Context, updatedTransaction entities.Transaction) error {
	if len(updatedTransaction.ID) <= 0 {
		return ErrInvalidTransaction
	}

	for index, transaction := range b.transactions {
		if transaction.ID == updatedTransaction.ID {
			if !transaction.CheckImmutable(&updatedTransaction) {
				return ErrInvalidUserSheba
			}

			b.transactionLocker.Lock()
			b.transactions[index] = updatedTransaction
			b.transactionLocker.Unlock()
			break
		}
	}

	return nil
}
