package test

import (
	"encoding/json"
	"net/url"
	"testing"

	consumableports "github.com/Logstay/tech-challange-execution-service/internal/ports/consumable"
	"github.com/Logstay/tech-challange-execution-service/pkg/null"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldGetConsumable(t *testing.T) {
	CleanViper()

	t.Run("should get consumable", func(t *testing.T) {
		defer CleanAll()

		body := consumableports.SaveConsumable{
			Name:        null.StringFrom("Óleo 5W40 - GET"),
			Price:       null.FloatFrom(79.9),
			Quantity:    null.IntFrom(5),
			Description: null.StringFrom("para testes de get"),
		}
		sresp, err := RestClient().SetBody(&body).Post("/v1/consumable")
		assert.Equal(t, 201, sresp.StatusCode())
		require.NoError(t, err)

		lresp, err := RestClient().Get("/v1/consumable?name=" + url.QueryEscape("Óleo 5W40 - GET"))
		assert.Equal(t, 200, lresp.StatusCode())
		require.NoError(t, err)
		var list struct {
			Data []struct {
				UID string `json:"uid"`
			} `json:"data"`
		}
		require.NoError(t, json.Unmarshal(lresp.Body(), &list))
		require.GreaterOrEqual(t, len(list.Data), 1)

		gresp, err := RestClient().Get("/v1/consumable/" + list.Data[0].UID)
		assert.Equal(t, 200, gresp.StatusCode())
		require.NoError(t, err)
	})

	t.Run("should return 400 for invalid uuid", func(t *testing.T) {
		resp, err := RestClient().Get("/v1/consumable/not-a-uuid")
		assert.Equal(t, 400, resp.StatusCode())
		assert.Contains(t, resp.String(), "uid invalido")
		require.NoError(t, err)
	})

	t.Run("should return 404 when not found", func(t *testing.T) {
		resp, err := RestClient().Get("/v1/consumable/550e8400-e29b-41d4-a716-446655440000")
		assert.Equal(t, 404, resp.StatusCode())
		assert.Contains(t, resp.String(), "peça/insumo não encontrado")
		require.NoError(t, err)
	})
}
