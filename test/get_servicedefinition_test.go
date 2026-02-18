package test

import (
	"encoding/json"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldGetServiceDefinition(t *testing.T) {
	CleanViper()

	t.Run("should get service definition", func(t *testing.T) {
		defer CleanAll()

		_, _ = RestClient().SetBody(map[string]any{
			"name":        "Revisão 10k",
			"price":       299.9,
			"description": "Revisão de 10 mil km",
		}).Post("/v1/service-definition")

		lr, err := RestClient().Get("/v1/service-definition?name=" + url.QueryEscape("Revisão 10k"))
		require.NoError(t, err)
		assert.Equal(t, 200, lr.StatusCode())
		var list struct {
			Data []struct {
				UID string `json:"uid"`
			} `json:"data"`
		}
		require.NoError(t, json.Unmarshal(lr.Body(), &list))
		require.GreaterOrEqual(t, len(list.Data), 1)

		resp, err := RestClient().Get("/v1/service-definition/" + list.Data[0].UID)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode())
	})

	t.Run("should 400 for invalid uuid", func(t *testing.T) {
		resp, err := RestClient().Get("/v1/service-definition/not-a-uuid")
		require.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode())
	})
}
