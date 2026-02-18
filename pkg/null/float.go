package null

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

type Float struct {
	sql.NullFloat64
}

func NewFloat(f float64, valid bool) Float {
	return Float{
		NullFloat64: sql.NullFloat64{
			Float64: f,
			Valid:   valid,
		},
	}
}

func FloatFrom(f float64) Float {
	return NewFloat(f, true)
}

func FloatFromPtr(f *float64) Float {
	if f == nil {
		return NewFloat(0, false)
	}
	return NewFloat(*f, true)
}

func (f *Float) ValueOrZero() float64 {
	if !f.Valid {
		return 0
	}
	return f.Float64
}

func (f *Float) UnmarshalJSON(data []byte) error {
	var err error
	var v any
	err = json.Unmarshal(data, &v)
	switch x := v.(type) {
	case float64:
		f.Float64 = float64(x)
	case map[string]any:
		err = json.Unmarshal(data, &f.NullFloat64)
	case nil:
		f.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Float", reflect.TypeOf(v).Name())
	}
	f.Valid = err == nil
	return err
}

func (f *Float) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == null {
		f.Valid = false
		return nil
	}
	var err error
	f.Float64, err = strconv.ParseFloat(string(text), 64)
	f.Valid = err == nil
	return err
}

func (f *Float) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return []byte(null), nil
	}
	return []byte(strconv.FormatFloat(f.Float64, 'f', -1, 64)), nil
}

func (f *Float) MarshalText() ([]byte, error) {
	if !f.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatFloat(f.Float64, 'f', -1, 64)), nil
}

func (f *Float) SetValid(n float64) {
	f.Float64 = n
	f.Valid = true
}

func (f *Float) Ptr() *float64 {
	if !f.Valid {
		return nil
	}
	return &f.Float64
}

func (f *Float) IsZero() bool {
	return !f.Valid
}
