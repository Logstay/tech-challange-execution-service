package inventory

import (
	"github.com/Logstay/tech-challange-wherehouse-service/pkg/null"
)

type Inventory interface {
	Quantity() null.Int
	ConsumableUID() null.String
	CreatedAt() null.DateTime
	UpdatedAt() null.DateTime
}

type inventory struct {
	quantity      null.Int
	consumableUID null.String
	createdAt     null.DateTime
	updatedAt     null.DateTime
}

func (i *inventory) CreatedAt() null.DateTime {
	return i.createdAt
}

func (i *inventory) UpdatedAt() null.DateTime {
	return i.updatedAt
}

func (i *inventory) Quantity() null.Int {
	return i.quantity
}

func (i *inventory) ConsumableUID() null.String {
	return i.consumableUID
}
