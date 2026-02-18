package null

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

type Int struct {
	sql.NullInt64
}

func NewInt(i int64, valid bool) Int {
	return Int{
		NullInt64: sql.NullInt64{
			Int64: i,
			Valid: valid,
		},
	}
}

func IntFrom(i int64) Int {
	return NewInt(i, true)
}

func IntFromPtr(i *int64) Int {
	if i == nil {
		return NewInt(0, false)
	}
	return NewInt(*i, true)
}

func (i *Int) UnmarshalJSON(data []byte) error {
	var err error
	var v any
	err = json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	switch v.(type) {
	case float64:

		err = json.Unmarshal(data, &i.Int64)
	case map[string]any:
		err = json.Unmarshal(data, &i.NullInt64)
	case nil:
		i.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Int", reflect.TypeOf(v).Name())
	}
	i.Valid = err == nil
	return err
}

func (i *Int) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == null {
		i.Valid = false
		return nil
	}
	var err error
	i.Int64, err = strconv.ParseInt(string(text), 10, 64)
	i.Valid = err == nil
	return err
}

func (i *Int) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return []byte(null), nil
	}
	return []byte(strconv.FormatInt(i.Int64, 10)), nil
}

func (i *Int) MarshalText() ([]byte, error) {
	if !i.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatInt(i.Int64, 10)), nil
}

func (i *Int) SetValid(n int64) {
	i.Int64 = n
	i.Valid = true
}

func (i *Int) Ptr() *int64 {
	if !i.Valid {
		return nil
	}
	return &i.Int64
}

func (i *Int) IsZero() bool {
	return !i.Valid
}

func (i *Int) String() string {
	if i.Ptr() == nil {
		return ""
	}
	return strconv.FormatInt(i.Int64, 10)
}
