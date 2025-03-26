package models

import "github.com/mohammadne/snapp-food/inernal/entities"

type TransferRequest struct {
	Price           int            `json:"price"`
	FromShebaNumber entities.Sheba `json:"fromShebaNumber"`
	ToShebaNumber   entities.Sheba `json:"ToShebaNumber"`
	Note            string         `json:"note"`
}
