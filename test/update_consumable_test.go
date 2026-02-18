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

func TestShouldUpdateConsumable(t *testing.T) {
	CleanViper()

	t.Run("should update consumable", func(t *testing.T) {
		defer CleanAll()

		body := consumableports.SaveConsumable{
			Name:        null.StringFrom("Detergente - UPDATE"),
			Price:       null.FloatFrom(5.5),
			Quantity:    null.IntFrom(2),
			Description: null.StringFrom("antes"),
		}
		sresp, err := RestClient().SetBody(&body).Post("/v1/consumable")
		assert.Equal(t, 201, sresp.StatusCode())
		require.NoError(t, err)

		lresp, err := RestClient().Get("/v1/consumable?name=" + url.QueryEscape("Detergente - UPDATE"))
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

		up := map[string]any{
			"name":        "Detergente - UPDATE OK",
			"price":       6.0,
			"description": "depois",
		}
		uresp, err := RestClient().SetBody(up).Put("/v1/consumable/" + uid)
		assert.Equal(t, 204, uresp.StatusCode())
		require.NoError(t, err)

		gresp, err := RestClient().Get("/v1/consumable/" + uid)
		assert.Equal(t, 200, gresp.StatusCode())
		require.NoError(t, err)
		var got struct {
			Name string `json:"name"`
		}
		require.NoError(t, json.Unmarshal(gresp.Body(), &got))
		assert.Equal(t, "Detergente - UPDATE OK", got.Name)
	})

	t.Run("should validate uid and payload", func(t *testing.T) {

		resp, err := RestClient().SetBody(map[string]any{"name": "x", "price": 1.0}).Put("/v1/consumable/invalid-uuid")
		assert.Equal(t, 400, resp.StatusCode())
		require.NoError(t, err)
	})
}
