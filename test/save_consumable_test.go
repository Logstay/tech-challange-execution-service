package test

import (
	"testing"

	consumableports "github.com/Logstay/tech-challange-execution-service/internal/ports/consumable"
	"github.com/Logstay/tech-challange-execution-service/pkg/null"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveConsumableValidations(t *testing.T) {
	CleanViper()

	t.Run("should save consumable", func(t *testing.T) {
		defer CleanAll()

		body := consumableports.SaveConsumable{
			Name:        null.StringFrom("Óleo 5W30"),
			Price:       null.FloatFrom(89.9),
			Quantity:    null.IntFrom(10),
			Description: null.StringFrom("Óleo sintético 5W30 - 1L"),
		}

		resp, err := RestClient().SetBody(&body).Post("/v1/consumable")
		require.NoError(t, err)
		assert.Equal(t, 201, resp.StatusCode())
	})

	t.Run("should not save when name is missing", func(t *testing.T) {

		body := consumableports.SaveConsumable{
			Price:       null.FloatFrom(10),
			Quantity:    null.IntFrom(1),
			Description: null.StringFrom("sem nome"),
		}

		resp, err := RestClient().SetBody(&body).Post("/v1/consumable")
		require.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode())
	})

	t.Run("should not save when price is missing", func(t *testing.T) {

		body := consumableports.SaveConsumable{
			Name:        null.StringFrom("Pano Microfibra"),
			Quantity:    null.IntFrom(1),
			Description: null.StringFrom("sem preço"),
		}

		resp, err := RestClient().SetBody(&body).Post("/v1/consumable")
		require.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode())
	})

	t.Run("should not save when price is zero", func(t *testing.T) {

		body := consumableports.SaveConsumable{
			Name:        null.StringFrom("Detergente"),
			Price:       null.FloatFrom(0),
			Quantity:    null.IntFrom(1),
			Description: null.StringFrom("preço zero"),
		}

		resp, err := RestClient().SetBody(&body).Post("/v1/consumable")
		require.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode())
	})

	t.Run("should not save when price is negative", func(t *testing.T) {
		body := consumableports.SaveConsumable{
			Name:        null.StringFrom("Graxa"),
			Price:       null.FloatFrom(-5.5),
			Quantity:    null.IntFrom(1),
			Description: null.StringFrom("preço negativo"),
		}
		resp, err := RestClient().SetBody(&body).Post("/v1/consumable")
		require.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode())
	})

	t.Run("should not save when quantity is missing", func(t *testing.T) {
		body := consumableports.SaveConsumable{
			Name:        null.StringFrom("Esponja"),
			Price:       null.FloatFrom(3.5),
			Description: null.StringFrom("sem quantidade"),
		}
		resp, err := RestClient().SetBody(&body).Post("/v1/consumable")
		require.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode())
	})

	t.Run("should not save when quantity is zero", func(t *testing.T) {
		body := consumableports.SaveConsumable{
			Name:        null.StringFrom("Flanela"),
			Price:       null.FloatFrom(2.5),
			Quantity:    null.IntFrom(0),
			Description: null.StringFrom("quantidade zero"),
		}
		resp, err := RestClient().SetBody(&body).Post("/v1/consumable")
		require.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode())
	})

	t.Run("should not save when quantity is negative", func(t *testing.T) {
		body := consumableports.SaveConsumable{
			Name:        null.StringFrom("Luva"),
			Price:       null.FloatFrom(12.0),
			Quantity:    null.IntFrom(-3),
			Description: null.StringFrom("quantidade negativa"),
		}
		resp, err := RestClient().SetBody(&body).Post("/v1/consumable")
		require.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode())
	})

	t.Run("should not save when name is duplicated", func(t *testing.T) {
		defer CleanAll()

		name := "Aditivo Radiador - Duplicado"
		first := consumableports.SaveConsumable{
			Name:        null.StringFrom(name),
			Price:       null.FloatFrom(29.9),
			Quantity:    null.IntFrom(3),
			Description: null.StringFrom("primeiro cadastro"),
		}
		_, _ = RestClient().SetBody(&first).Post("/v1/consumable")

		resp, err := RestClient().SetBody(&first).Post("/v1/consumable")
		require.NoError(t, err)
		assert.Contains(t, []int{400, 409}, resp.StatusCode())
	})
}
