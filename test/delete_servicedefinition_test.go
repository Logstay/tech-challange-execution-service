package test

import (
	"encoding/json"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldDeleteServiceDefinition(t *testing.T) {
	CleanViper()

	t.Run("should delete service definition", func(t *testing.T) {
		defer CleanAll()

		_, _ = RestClient().SetBody(map[string]any{"name": "Flush Arrefecimento", "price": 249.9, "description": "flush"}).Post("/v1/service-definition")

		lr, _ := RestClient().Get("/v1/service-definition?name=" + url.QueryEscape("Flush Arrefecimento"))
		var list struct {
			Data []struct {
				UID string `json:"uid"`
			} `json:"data"`
		}
		_ = json.Unmarshal(lr.Body(), &list)
		uid := list.Data[0].UID

		resp, err := RestClient().Delete("/v1/service-definition/" + uid)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode())

		g, _ := RestClient().Get("/v1/service-definition/" + uid)
		assert.Equal(t, 404, g.StatusCode())
	})

	t.Run("should 400 for invalid uuid", func(t *testing.T) {
		resp, err := RestClient().Delete("/v1/service-definition/not-a-uuid")
		require.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode())
	})
}
