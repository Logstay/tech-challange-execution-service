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

func TestUpdateServiceDefinitionHandler(t *testing.T) {
	repositoryMock := &ports.RepositoryMock{}

	tr, err := tracing.NewTracingProvider(&config.Config{})
	require.NoError(t, err)

	service := &UpdateServiceDefinitionHandler{
		Repository: repositoryMock,
		Tracer:     tr.NewTracer("test-tracing"),
	}

	t.Run("should update serviceDefinition", func(t *testing.T) {
		cmd := ports.UpdateServiceDefinition{
			UID:   null.StringFrom("6ecf221f-f933-4efb-9d57-3dd9f62ee737"),
			Price: null.FloatFrom(50),
			Name:  null.StringFrom("Óleo 5W30"),
		}

		repositoryMock.On("Get", mock.Anything, cmd.UID.String).Times(1).Return(ports.GetServiceDefinition{UID: cmd.UID}, nil)

		repositoryMock.On("Update", mock.Anything, cmd.UID.String, mock.Anything).Times(1).Return(nil)

		err := service.Handle(t.Context(), cmd)

		require.NoError(t, err)
		repositoryMock.AssertExpectations(t)
	})

	t.Run("should not update serviceDefinition when update fails", func(t *testing.T) {
		cmd := ports.UpdateServiceDefinition{
			UID:   null.StringFrom("6ecf221f-f933-4efb-9d57-3dd9f62ee737"),
			Price: null.FloatFrom(50),
			Name:  null.StringFrom("Óleo 5W30"),
		}

		repositoryMock.On("Get", mock.Anything, cmd.UID.String).Times(1).Return(ports.GetServiceDefinition{UID: cmd.UID}, nil)

		repositoryMock.On("Update", mock.Anything, cmd.UID.String, mock.Anything).Times(1).Return(errors.New("error"))

		err := service.Handle(t.Context(), cmd)

		require.Error(t, err)
		repositoryMock.AssertExpectations(t)
	})

	t.Run("should not update serviceDefinition when get fails", func(t *testing.T) {
		cmd := ports.UpdateServiceDefinition{
			UID:   null.StringFrom("6ecf221f-f933-4efb-9d57-3dd9f62ee737"),
			Price: null.FloatFrom(50),
			Name:  null.StringFrom("Óleo 5W30"),
		}

		repositoryMock.On("Get", mock.Anything, cmd.UID.String).Times(1).Return(ports.GetServiceDefinition{}, errors.New("error"))

		err := service.Handle(t.Context(), cmd)

		require.Error(t, err)
		repositoryMock.AssertExpectations(t)
	})

	t.Run("should not update serviceDefinition when name is empty", func(t *testing.T) {
		cmd := ports.UpdateServiceDefinition{
			UID:   null.StringFrom("6ecf221f-f933-4efb-9d57-3dd9f62ee737"),
			Price: null.FloatFrom(50),
		}

		repositoryMock.On("Get", mock.Anything, cmd.UID.String).Times(1).Return(ports.GetServiceDefinition{UID: cmd.UID}, nil)

		err := service.Handle(t.Context(), cmd)

		require.EqualError(t, err, "nome do serviço é obrigatorio")
		repositoryMock.AssertExpectations(t)
	})

	t.Run("should not update serviceDefinition when price is empty", func(t *testing.T) {
		cmd := ports.UpdateServiceDefinition{
			UID:  null.StringFrom("6ecf221f-f933-4efb-9d57-3dd9f62ee737"),
			Name: null.StringFrom("Óleo 5W30"),
		}

		repositoryMock.On("Get", mock.Anything, cmd.UID.String).Times(1).Return(ports.GetServiceDefinition{UID: cmd.UID}, nil)

		err := service.Handle(t.Context(), cmd)

		require.EqualError(t, err, "preço do serviço é obrigatorio")
		repositoryMock.AssertExpectations(t)
	})

}
