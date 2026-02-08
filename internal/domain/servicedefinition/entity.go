package servicedefinition

import (
	"github.com/Logstay/tech-challange-wherehouse-service/pkg/null"
)

type ServiceDefinition interface {
	Name() null.String
	Description() null.String
	Price() null.Float
	CreatedAt() null.DateTime
	UpdatedAt() null.DateTime
}

type serviceDefinition struct {
	name        null.String
	description null.String
	price       null.Float
	createdAt   null.DateTime
	updatedAt   null.DateTime
}

func (c *serviceDefinition) Price() null.Float {
	return c.price
}

func (c *serviceDefinition) CreatedAt() null.DateTime {
	return c.createdAt
}

func (c *serviceDefinition) UpdatedAt() null.DateTime {
	return c.updatedAt
}

func (c *serviceDefinition) Name() null.String {
	return c.name
}

func (c *serviceDefinition) Description() null.String {
	return c.description
}
