package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldListServiceDefinition(t *testing.T) {
	CleanViper()

	t.Run("should list and filter service definitions", func(t *testing.T) {
		defer CleanAll()

		_, _ = RestClient().SetBody(map[string]any{"name": "Troca Filtro Ar", "price": 59.9, "description": "Filtro de ar"}).Post("/v1/service-definition")
		_, _ = RestClient().SetBody(map[string]any{"name": "Troca Filtro Olio", "price": 79.9, "description": "Filtro de óleo"}).Post("/v1/service-definition")

		r, err := RestClient().Get("/v1/service-definition?page=0&limit=10")
		require.NoError(t, err)
		assert.Equal(t, 200, r.StatusCode())

		f, err := RestClient().Get("/v1/service-definition?name=Troca%20Filtro")
		require.NoError(t, err)
		assert.Equal(t, 200, f.StatusCode())
	})
}
