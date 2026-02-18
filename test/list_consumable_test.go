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

func TestShouldListConsumable(t *testing.T) {
	CleanViper()

	t.Run("should list consumables with pagination and filter", func(t *testing.T) {
		defer CleanAll()

		for _, seed := range []consumableports.SaveConsumable{
			{
				Name:        null.StringFrom("Pano Microfibra - LIST"),
				Price:       null.FloatFrom(12.5),
				Quantity:    null.IntFrom(5),
				Description: null.StringFrom("pano"),
			},
			{
				Name:        null.StringFrom("Óleo 10W40 - LIST"),
				Price:       null.FloatFrom(59.9),
				Quantity:    null.IntFrom(3),
				Description: null.StringFrom("óleo"),
			},
		} {
			resp, err := RestClient().SetBody(&seed).Post("/v1/consumable")
			assert.Equal(t, 201, resp.StatusCode())
			require.NoError(t, err)
		}

		resp, err := RestClient().Get("/v1/consumable?page=0&limit=10")
		assert.Equal(t, 200, resp.StatusCode())
		require.NoError(t, err)

		var out struct {
			Data       []map[string]any `json:"data"`
			Pagination struct {
				Page  int64 `json:"page"`
				Limit int64 `json:"limit"`
				Total int64 `json:"total"`
			} `json:"pagination"`
		}
		require.NoError(t, json.Unmarshal(resp.Body(), &out))
		assert.GreaterOrEqual(t, len(out.Data), 2)
		assert.EqualValues(t, 0, out.Pagination.Page)
		assert.EqualValues(t, 10, out.Pagination.Limit)
		assert.GreaterOrEqual(t, out.Pagination.Total, int64(2))

		resp2, err := RestClient().Get("/v1/consumable?name=" + url.QueryEscape("LIST"))
		assert.Equal(t, 200, resp2.StatusCode())
		require.NoError(t, err)

		var out2 struct {
			Data []map[string]any `json:"data"`
		}
		require.NoError(t, json.Unmarshal(resp2.Body(), &out2))

		assert.GreaterOrEqual(t, len(out2.Data), 2)
	})
}
