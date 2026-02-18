package test

import (
	"encoding/json"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldUpdateServiceDefinition(t *testing.T) {
	CleanViper()

	t.Run("should update service definition", func(t *testing.T) {
		defer CleanAll()

		_, _ = RestClient().SetBody(map[string]any{"name": "Alinhamento 3D", "price": 149.9, "description": "alinhamento"}).Post("/v1/service-definition")

		lr, _ := RestClient().Get("/v1/service-definition?name=" + url.QueryEscape("Alinhamento 3D"))
		var list struct {
			Data []struct {
				UID string `json:"uid"`
			} `json:"data"`
		}
		_ = json.Unmarshal(lr.Body(), &list)
		uid := list.Data[0].UID

		resp, err := RestClient().SetBody(map[string]any{"name": "Alinhamento e Balanceamento", "price": 189.9, "description": "alinhamento + balanceamento"}).Put("/v1/service-definition/" + uid)
		require.NoError(t, err)
		assert.Equal(t, 204, resp.StatusCode())

		g, _ := RestClient().Get("/v1/service-definition/" + uid)
		assert.Equal(t, 200, g.StatusCode())
		assert.Contains(t, g.String(), "Alinhamento e Balanceamento")
	})

	t.Run("should 400 for invalid uuid", func(t *testing.T) {
		resp, err := RestClient().SetBody(map[string]any{"name": "X", "price": 1}).Put("/v1/service-definition/not-a-uuid")
		require.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode())
	})
}
