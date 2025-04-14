package usecases

import (
	"context"

	"github.com/mohammadne/bank-teller/inernal/entities"
	"github.com/mohammadne/bank-teller/inernal/repository"
)

type Users interface {
	ListUsers(ctx context.Context) []entities.User
}

func NewUsers(b repository.Bank) Users {
	return &users{bank: b}
}

type users struct {
	bank repository.Bank
}

// ListUsers implements Users.
func (u *users) ListUsers(ctx context.Context) []entities.User {
	return u.bank.ListUsers(ctx)
}
