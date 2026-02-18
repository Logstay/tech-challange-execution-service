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

func TestShouldDeleteConsumable(t *testing.T) {
	CleanViper()

	t.Run("should delete consumable", func(t *testing.T) {
		defer CleanAll()

		body := consumableports.SaveConsumable{
			Name:        null.StringFrom("Pano - DELETE"),
			Price:       null.FloatFrom(9.9),
			Quantity:    null.IntFrom(1),
			Description: null.StringFrom("delete"),
		}
		sresp, err := RestClient().SetBody(&body).Post("/v1/consumable")
		assert.Equal(t, 201, sresp.StatusCode())
		require.NoError(t, err)

		lresp, err := RestClient().Get("/v1/consumable?name=" + url.QueryEscape("Pano - DELETE"))
		assert.Equal(t, 200, lresp.StatusCode())
		require.NoError(t, err)
		var list struct {
			Data []struct {
				UID string `json:"uid"`
			} `json:"data"`
		}
		require.NoError(t, json.Unmarshal(lresp.Body(), &list))
		require.GreaterOrEqual(t, len(list.Data), 1)
		uid := list.Data[0].UID

		dresp, err := RestClient().Delete("/v1/consumable/" + uid)
		assert.Equal(t, 204, dresp.StatusCode())
		require.NoError(t, err)

		gresp, err := RestClient().Get("/v1/consumable/" + uid)
		assert.Equal(t, 404, gresp.StatusCode())
		require.NoError(t, err)
	})

	t.Run("should validate uid", func(t *testing.T) {
		resp, err := RestClient().Delete("/v1/consumable/not-a-uuid")
		assert.Equal(t, 400, resp.StatusCode())
		require.NoError(t, err)
	})
}
