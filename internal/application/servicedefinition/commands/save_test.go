package command

import (
	"errors"

	"github.com/Logstay/tech-challange-execution-service/config"
	ports "github.com/Logstay/tech-challange-execution-service/internal/ports/servicedefinition"
	"github.com/Logstay/tech-challange-execution-service/pkg/null"
	"github.com/Logstay/tech-challange-execution-service/pkg/tracing"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"testing"
)

func TestSaveServiceDefinitionHandler(t *testing.T) {
	repositoryMock := &ports.RepositoryMock{}

	tr, err := tracing.NewTracingProvider(&config.Config{})
	require.NoError(t, err)

	service := &SaveServiceDefinitionHandler{
		Repository: repositoryMock,
		Tracer:     tr.NewTracer("test-tracing"),
	}

	t.Run("should save serviceDefinition", func(t *testing.T) {
		cmd := ports.SaveServiceDefinition{
			Price: null.FloatFrom(200),
			Name:  null.StringFrom("Filtro de Óleo"),
		}

		repositoryMock.On("Exist", mock.Anything, mock.Anything).Times(1).Return(false, nil)

		repositoryMock.On("Save", mock.Anything, mock.Anything).Times(1).Return(nil)

		err := service.Handle(t.Context(), cmd)

		require.NoError(t, err)
		repositoryMock.AssertExpectations(t)
	})

	t.Run("should not save serviceDefinition when save fails", func(t *testing.T) {
		cmd := ports.SaveServiceDefinition{
			Price: null.FloatFrom(50),
			Name:  null.StringFrom("Óleo 5W30"),
		}

		repositoryMock.On("Exist", mock.Anything, mock.Anything).Times(1).Return(false, nil)

		repositoryMock.On("Save", mock.Anything, mock.Anything).Times(1).Return(errors.New("error"))

		err := service.Handle(t.Context(), cmd)

		require.Error(t, err)
		repositoryMock.AssertExpectations(t)
	})

	t.Run("should not save serviceDefinition when exist fails", func(t *testing.T) {
		cmd := ports.SaveServiceDefinition{
			Price: null.FloatFrom(50),
			Name:  null.StringFrom("Óleo 5W30"),
		}

		repositoryMock.On("Exist", mock.Anything, mock.Anything).Times(1).Return(false, errors.New("error"))

		err := service.Handle(t.Context(), cmd)

		require.Error(t, err)
		repositoryMock.AssertExpectations(t)
	})

	t.Run("should not save serviceDefinition when name is empty", func(t *testing.T) {
		cmd := ports.SaveServiceDefinition{
			Price: null.FloatFrom(50),
		}

		repositoryMock.On("Exist", mock.Anything, mock.Anything).Times(1).Return(false, nil)

		err := service.Handle(t.Context(), cmd)

		require.EqualError(t, err, "nome do serviço é obrigatorio")
		repositoryMock.AssertExpectations(t)
	})

	t.Run("should not save serviceDefinition when price is empty", func(t *testing.T) {
		cmd := ports.SaveServiceDefinition{
			Name: null.StringFrom("Óleo 5W30"),
		}

		repositoryMock.On("Exist", mock.Anything, mock.Anything).Times(1).Return(false, nil)

		err := service.Handle(t.Context(), cmd)

		require.EqualError(t, err, "preço do serviço é obrigatorio")
		repositoryMock.AssertExpectations(t)
	})

}
