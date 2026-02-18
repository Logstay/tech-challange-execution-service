package consumable

import (
	httperr "github.com/Logstay/tech-challange-execution-service/pkg/errors"
	"github.com/Logstay/tech-challange-execution-service/pkg/null"
)

type ConsumableFactory interface {
	Build() Consumable
	CreatedAt(createdAt null.DateTime)
	UpdatedAt(updatedAt null.DateTime)
	Price(price null.Float) (err error)
	Name(name null.String) (err error)
	Description(description null.String) (err error)
	Quantity(quantity null.Int) (err error)
}

type consumableFactory struct {
	consumable *consumable
}

func NewConsumableFactory() ConsumableFactory {
	return &consumableFactory{
		consumable: &consumable{
			code: GenerateConsumableCode(),
		},
	}
}

func (f *consumableFactory) Build() Consumable {
	return f.consumable
}

func (f *consumableFactory) CreatedAt(createdAt null.DateTime) {
	if createdAt.IsZero() {
		return
	}

	f.consumable.createdAt = createdAt
}

func (f *consumableFactory) UpdatedAt(updatedAt null.DateTime) {
	if updatedAt.IsZero() {
		return
	}

	f.consumable.updatedAt = updatedAt
}

func (f *consumableFactory) Price(price null.Float) (err error) {
	if price.Float64 <= 0 {
		return httperr.RequiredError("preço da peça/insumo é obrigatorio")
	}

	f.consumable.price = price
	return
}

func (f *consumableFactory) Name(name null.String) (err error) {
	if name.String == "" {
		return httperr.RequiredError("nome da peça/insumo é obrigatorio")
	}

	f.consumable.name = name
	return
}

func (f *consumableFactory) Description(description null.String) (err error) {
	f.consumable.description = description
	return
}

func (f *consumableFactory) Quantity(quantity null.Int) (err error) {
	if quantity.Int64 <= 0 {
		return httperr.RequiredError("quantidade da peça/insumo é obrigatorio")
	}
	f.consumable.quantity = quantity
	return
}
