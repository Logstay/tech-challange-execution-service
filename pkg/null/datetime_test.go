package null

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseDateTimeString(t *testing.T) {
	data1 := DateTimeFrom("17/04/2015 08:10:10")
	data2 := NewDateTime(time.Date(2015, time.April, 17, 8, 10, 10, 0, time.UTC))
	data3 := DateTimeFrom("")

	assert.Equal(t, "17/04/2015 08:10:10", data1.String())
	assert.Equal(t, "17/04/2015 08:10:10", data2.String())
	assert.Empty(t, data3.String())
}

func TestCriarDateTimeFrom(t *testing.T) {
	data1 := DateTimeFrom("17/04/2015 08:10:10")
	data2 := NewDateTime(time.Date(2015, time.April, 17, 8, 10, 10, 0, time.UTC))

	assert.NotNil(t, data1)
	assert.Equal(t, data1, data2)
	assert.True(t, data1.Valid)
}

func TestCriarDateTimeFromInvalidaSeInputVazio(t *testing.T) {
	data1 := DateTimeFrom("")
	data2 := NewDateTime(time.Date(2015, time.January, 17, 0, 0, 0, 0, time.UTC))

	assert.NotNil(t, data1)
	assert.NotEqual(t, data1, data2)
	assert.False(t, data1.Valid)
}

func TestCriarDateTimeNulo(t *testing.T) {
	data := DateTime{}

	var value *time.Time
	err := data.Scan(value)

	require.NoError(t, err)
	assert.False(t, data.Valid)
}

func TestLerDateTimeNulo(t *testing.T) {
	data := DateTime{}

	var value *time.Time
	err := data.Scan(value)

	resp, _ := data.Value()

	require.NoError(t, err)
	assert.Nil(t, resp)
}

func TestCriarDateTime(t *testing.T) {
	data := DateTime{}

	err := data.Scan(time.Now())

	require.NoError(t, err)
	assert.True(t, data.Valid)
}

func TestConverterDateTimeParaJson(t *testing.T) {
	data := DateTime{}

	now := time.Date(2015, time.April, 17, 8, 10, 10, 10, time.UTC)
	err := data.Scan(now)
	require.NoError(t, err)

	bytes, err := data.MarshalJSON()

	require.NoError(t, err)
	assert.Equal(t, `"17/04/2015 08:10:10"`, string(bytes))
}

func TestConverterJsonParaDateTime(t *testing.T) {
	now := time.Date(2015, time.April, 17, 8, 10, 10, 0, time.UTC)

	data := DateTime{}
	err := data.Scan(now)
	require.NoError(t, err)

	err = json.Unmarshal([]byte(`"17/04/2015 08:10:10"`), &data)

	require.NoError(t, err)
	assert.Equal(t, now.String(), data.DateTime.String())
}

func TestCompararDateTime(t *testing.T) {
	data1 := NewDateTime(time.Date(2015, time.April, 17, 0, 0, 0, 0, time.UTC))
	data2 := NewDateTime(time.Date(2015, time.April, 17, 0, 0, 0, 0, time.UTC))

	assert.True(t, data1.Equal(&data2))
}
