package inventory

import (
	httperr "github.com/Logstay/tech-challange-wherehouse-service/pkg/errors"
	"github.com/Logstay/tech-challange-wherehouse-service/pkg/null"
)

type InventoryFactory interface {
	Build() (Inventory, error)
	Quantity(quantity null.Int) (err error)
	ConsumableUID(consumableUID null.String)
	CreatedAt(createdAt null.DateTime)
	UpdatedAt(updatedAt null.DateTime)
}

type inventoryFactory struct {
	inventory *inventory
}

func NewInventoryFactory(id null.String) InventoryFactory {
	return &inventoryFactory{
		inventory: &inventory{},
	}
}

func (f *inventoryFactory) Build() (Inventory, error) {
	return f.inventory, nil
}

func (f *inventoryFactory) CreatedAt(createdAt null.DateTime) {
	if createdAt.IsZero() {
		return
	}

	f.inventory.createdAt = createdAt
}

func (f *inventoryFactory) UpdatedAt(updatedAt null.DateTime) {
	if updatedAt.IsZero() {
		return
	}

	f.inventory.updatedAt = updatedAt
}

func (f *inventoryFactory) Quantity(quantity null.Int) (err error) {
	if quantity.Int64 <= 0 {
		return httperr.InvalidPriceError("quantidade deve ser maior que 0")

	}

	f.inventory.quantity = quantity
	return
}

func (f *inventoryFactory) ConsumableUID(consumableUID null.String) {
	f.inventory.consumableUID = consumableUID
}
