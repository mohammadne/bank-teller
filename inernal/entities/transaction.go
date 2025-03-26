package entities

import "time"

type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusConfirmed TransactionStatus = "confirmed"
	TransactionStatusCanceled  TransactionStatus = "canceled"
)

type Transaction struct {
	ID        string            `json:"id"`
	Status    TransactionStatus `json:"status"`
	From      Sheba             `json:"fromShebaNumber"`
	To        Sheba             `json:"ToShebaNumber"`
	Amount    int               `json:"price"`
	CreatedAt time.Time         `json:"createdAt"`
}
