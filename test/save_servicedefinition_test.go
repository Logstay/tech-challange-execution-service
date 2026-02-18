package test

import (
	"testing"

	ports "github.com/Logstay/tech-challange-execution-service/internal/ports/servicedefinition"
	"github.com/Logstay/tech-challange-execution-service/pkg/null"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldSaveServiceDefinition(t *testing.T) {
	CleanViper()

	t.Run("should save service definition", func(t *testing.T) {
		defer CleanAll()
		body := ports.SaveServiceDefinition{
			Name:        null.StringFrom("Troca de Óleo Full"),
			Price:       null.FloatFrom(199.9),
			Description: null.StringFrom("Inclui filtro e revisão"),
		}
		resp, err := RestClient().SetBody(&body).Post("/v1/service-definition")
		require.NoError(t, err)
		assert.Equal(t, 201, resp.StatusCode())
	})

	t.Run("should not save without name", func(t *testing.T) {
		body := ports.SaveServiceDefinition{
			Price:       null.FloatFrom(99.9),
			Description: null.StringFrom("Sem nome"),
		}
		resp, err := RestClient().SetBody(&body).Post("/v1/service-definition")
		require.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode())
	})

	t.Run("should not save with zero or negative price", func(t *testing.T) {
		body := ports.SaveServiceDefinition{
			Name:        null.StringFrom("Balanceamento"),
			Price:       null.FloatFrom(0),
			Description: null.StringFrom("valor zero"),
		}
		resp, err := RestClient().SetBody(&body).Post("/v1/service-definition")
		require.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode())

		body2 := ports.SaveServiceDefinition{
			Name:        null.StringFrom("Alinhamento"),
			Price:       null.FloatFrom(-10),
			Description: null.StringFrom("valor negativo"),
		}
		resp2, err := RestClient().SetBody(&body2).Post("/v1/service-definition")
		require.NoError(t, err)
		assert.Equal(t, 400, resp2.StatusCode())
	})
}
