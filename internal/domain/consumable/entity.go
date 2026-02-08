package consumable

import (
	"github.com/Logstay/tech-challange-wherehouse-service/pkg/null"
	gonanoid "github.com/matoous/go-nanoid"
)

type Consumable interface {
	Name() null.String
	Code() null.String
	Description() null.String
	Price() null.Float
	Quantity() null.Int
	CreatedAt() null.DateTime
	UpdatedAt() null.DateTime
}

type consumable struct {
	name        null.String
	description null.String
	code        null.String
	price       null.Float
	quantity    null.Int
	createdAt   null.DateTime
	updatedAt   null.DateTime
}

func (c *consumable) Price() null.Float {
	return c.price
}

func (c *consumable) CreatedAt() null.DateTime {
	return c.createdAt
}

func (c *consumable) UpdatedAt() null.DateTime {
	return c.updatedAt
}

func (c *consumable) Name() null.String {
	return c.name
}

func (c *consumable) Quantity() null.Int {
	return c.quantity
}

func (c *consumable) Description() null.String {
	return c.description
}

func (c *consumable) Code() null.String {
	return c.code
}

const consumableCodeLength = 18

func GenerateConsumableCode() null.String {
	return null.StringFrom(gonanoid.MustGenerate(
		"0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz",
		consumableCodeLength))
}
