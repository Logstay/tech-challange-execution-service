package command

import (
	"errors"

	"github.com/Logstay/tech-challange-wherehouse-service/config"
	ports "github.com/Logstay/tech-challange-wherehouse-service/internal/ports/consumable"
	"github.com/Logstay/tech-challange-wherehouse-service/pkg/null"
	"github.com/Logstay/tech-challange-wherehouse-service/pkg/tracing"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"testing"
)

func TestDeleteConsumableHandler(t *testing.T) {
	repositoryMock := &ports.RepositoryMock{}

	tr, err := tracing.NewTracingProvider(&config.Config{})
	require.NoError(t, err)

	service := &DeleteConsumableHandler{
		Repository: repositoryMock,
		Tracer:     tr.NewTracer("test-tracing"),
	}

	uid := "6ecf221f-f933-4efb-9d57-3dd9f62ee737"

	t.Run("should delete consumable", func(t *testing.T) {

		repositoryMock.On("Get", mock.Anything, uid).Times(1).Return(ports.GetConsumable{UID: null.StringFrom(uid)}, nil)

		repositoryMock.On("Delete", mock.Anything, uid).Times(1).Return(nil)

		err := service.Handle(t.Context(), uid)

		require.NoError(t, err)
		repositoryMock.AssertExpectations(t)
	})

	t.Run("should not delete consumable when delete fails", func(t *testing.T) {

		repositoryMock.On("Get", mock.Anything, uid).Times(1).Return(ports.GetConsumable{UID: null.StringFrom(uid)}, nil)

		repositoryMock.On("Delete", mock.Anything, uid).Times(1).Return(errors.New("error"))

		err := service.Handle(t.Context(), uid)

		require.Error(t, err)
		repositoryMock.AssertExpectations(t)
	})

	t.Run("should not delete consumable when get fails", func(t *testing.T) {

		repositoryMock.On("Get", mock.Anything, uid).Times(1).Return(ports.GetConsumable{}, errors.New("error"))

		err := service.Handle(t.Context(), uid)

		require.Error(t, err)
		repositoryMock.AssertExpectations(t)
	})

}
