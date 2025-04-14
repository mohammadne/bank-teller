package entities

import (
	"errors"
	"time"
)

type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusConfirmed TransactionStatus = "confirmed"
	TransactionStatusCanceled  TransactionStatus = "canceled"
)

var (
	ErrInvalidStatus = errors.New("ErrInvalidStatus")
)

func ToTransactionStatus(status string) (TransactionStatus, error) {
	switch TransactionStatus(status) {
	case TransactionStatusPending:
		return TransactionStatusPending, nil
	case TransactionStatusCanceled:
		return TransactionStatusCanceled, nil
	case TransactionStatusConfirmed:
		return TransactionStatusConfirmed, nil
	}

	return "", ErrInvalidStatus
}

type Transaction struct {
	ID        string            `json:"id"`
	Status    TransactionStatus `json:"status"`
	From      Sheba             `json:"fromShebaNumber"`
	To        Sheba             `json:"ToShebaNumber"`
	Amount    int               `json:"price"`
	CreatedAt time.Time         `json:"createdAt"`
}

func (t *Transaction) CheckImmutable(n *Transaction) bool {
	return t.ID == n.ID && t.From == n.From && t.To == n.To && t.Amount == n.Amount && t.CreatedAt == n.CreatedAt
}

type Pool []Transaction
