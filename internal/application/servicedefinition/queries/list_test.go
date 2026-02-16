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

func TestListServiceDefinitionHandler(t *testing.T) {
	repositoryMock := &ports.RepositoryMock{}

	tr, err := tracing.NewTracingProvider(&config.Config{})
	require.NoError(t, err)

	service := &ListServiceDefinitionHandler{
		Repository: repositoryMock,
		Tracer:     tr.NewTracer("test-tracing"),
	}

	t.Run("should get serviceDefinition", func(t *testing.T) {
		filter := ports.Filter{}

		repositoryMock.On("List", mock.Anything, filter).Times(1).Return([]ports.GetServiceDefinition{
			{
				UID:       null.StringFrom("6ecf221f-f933-4efb-9d57-3dd9f62ee737"),
				Price:     null.FloatFrom(50),
				Name:      null.StringFrom("Óleo 5W30"),
				CreatedAt: null.DateTimeFrom("1994-07-17 00:00:00"),
			},
			{
				UID:       null.StringFrom("0518fb01-11a7-4ce4-a81e-5198c5c48d0a"),
				Price:     null.FloatFrom(200),
				Name:      null.StringFrom("Filtro de Óleo"),
				CreatedAt: null.DateTimeFrom("2000-07-17 00:00:00"),
			},
		}, int64(2), nil)

		result, total, err := service.Handle(t.Context(), filter)

		require.NotEmpty(t, result)
		require.NoError(t, err)
		require.Equal(t, int64(2), total)
		repositoryMock.AssertExpectations(t)
	})

	t.Run("should get serviceDefinition with filter", func(t *testing.T) {
		filter := ports.Filter{
			Name:  null.StringFrom("name1"),
			Page:  null.IntFrom(0),
			Limit: null.IntFrom(10),
		}

		repositoryMock.On("List", mock.Anything, filter).Times(1).Return([]ports.GetServiceDefinition{{
			UID:         null.StringFrom("6ecf221f-f933-4efb-9d57-3dd9f62ee737"),
			Price:       null.FloatFrom(50),
			Name:        null.StringFrom("Óleo 5W30"),
			Description: null.StringFrom("description1"),
			UpdatedAt:   null.DateTimeFrom("1994-07-17 00:00:00"),
			CreatedAt:   null.DateTimeFrom("1994-07-17 00:00:00"),
		}}, int64(1), nil)

		result, total, err := service.Handle(t.Context(), filter)

		require.NotEmpty(t, result)
		require.NoError(t, err)
		require.Equal(t, int64(1), total)
		repositoryMock.AssertExpectations(t)
	})

	t.Run("should not get serviceDefinition when List fails", func(t *testing.T) {
		filter := ports.Filter{
			Name:  null.StringFrom("name1"),
			Page:  null.IntFrom(0),
			Limit: null.IntFrom(10),
		}

		repositoryMock.On("List", mock.Anything, filter).Times(1).Return([]ports.GetServiceDefinition{}, int64(0), errors.New("error"))

		result, total, err := service.Handle(t.Context(), filter)

		require.Empty(t, result)
		require.Error(t, err)
		require.Equal(t, int64(0), total)
		repositoryMock.AssertExpectations(t)
	})

}
