package command

import (
	"errors"

	"github.com/Logstay/tech-challange-wherehouse-service/config"
	ports "github.com/Logstay/tech-challange-wherehouse-service/internal/ports/servicedefinition"
	"github.com/Logstay/tech-challange-wherehouse-service/pkg/null"
	"github.com/Logstay/tech-challange-wherehouse-service/pkg/tracing"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"testing"
)

func TestGetServiceDefinitionHandler(t *testing.T) {
	repositoryMock := &ports.RepositoryMock{}

	tr, err := tracing.NewTracingProvider(&config.Config{})
	require.NoError(t, err)

	service := &GetServiceDefinitionHandler{
		Repository: repositoryMock,
		Tracer:     tr.NewTracer("test-tracing"),
	}

	uid := "6ecf221f-f933-4efb-9d57-3dd9f62ee737"

	t.Run("should get serviceDefinition", func(t *testing.T) {

		repositoryMock.On("Get", mock.Anything, uid).Times(1).Return(ports.GetServiceDefinition{
			UID:         null.StringFrom(uid),
			Price:       null.FloatFrom(50),
			Name:        null.StringFrom("Óleo 5W30"),
			Description: null.StringFrom("Descrição do item"),
			UpdatedAt:   null.DateTimeFrom("1994-07-17 00:00:00"),
			CreatedAt:   null.DateTimeFrom("1994-07-17 00:00:00"),
		}, nil)

		result, err := service.Handle(t.Context(), uid)

		require.NotEmpty(t, result)
		require.NoError(t, err)
		repositoryMock.AssertExpectations(t)
	})

	t.Run("should not get serviceDefinition when Get fails", func(t *testing.T) {

		repositoryMock.On("Get", mock.Anything, uid).Times(1).Return(ports.GetServiceDefinition{},
			errors.New("error"))

		result, err := service.Handle(t.Context(), uid)

		require.Empty(t, result)
		require.Error(t, err)
		repositoryMock.AssertExpectations(t)
	})

}
