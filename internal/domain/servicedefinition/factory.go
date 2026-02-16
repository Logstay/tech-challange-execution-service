package servicedefinition

import (
	httperr "github.com/Logstay/tech-challange-execution-service/pkg/errors"
	"github.com/Logstay/tech-challange-execution-service/pkg/null"
)

type ServiceDefinitionFactory interface {
	Build() ServiceDefinition
	CreatedAt(createdAt null.DateTime)
	UpdatedAt(updatedAt null.DateTime)
	Price(price null.Float) (err error)
	Name(name null.String) (err error)
	Description(description null.String) (err error)
}

type serviceDefinitionFactory struct {
	serviceDefinition *serviceDefinition
}

func NewServiceDefinitionFactory() ServiceDefinitionFactory {
	return &serviceDefinitionFactory{
		serviceDefinition: &serviceDefinition{},
	}
}

func (f *serviceDefinitionFactory) Build() ServiceDefinition {
	return f.serviceDefinition
}

func (f *serviceDefinitionFactory) CreatedAt(createdAt null.DateTime) {
	if createdAt.IsZero() {
		return
	}

	f.serviceDefinition.createdAt = createdAt
}

func (f *serviceDefinitionFactory) UpdatedAt(updatedAt null.DateTime) {
	if updatedAt.IsZero() {
		return
	}

	f.serviceDefinition.updatedAt = updatedAt
}

func (f *serviceDefinitionFactory) Price(price null.Float) (err error) {
	if price.Float64 <= 0 {
		return httperr.RequiredError("preço do serviço é obrigatorio")
	}

	f.serviceDefinition.price = price
	return
}

func (f *serviceDefinitionFactory) Name(name null.String) (err error) {
	if name.String == "" {
		return httperr.RequiredError("nome do serviço é obrigatorio")
	}

	f.serviceDefinition.name = name
	return
}

func (f *serviceDefinitionFactory) Description(description null.String) (err error) {
	f.serviceDefinition.description = description
	return
}
