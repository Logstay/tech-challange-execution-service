package null

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

const (
	null = "null"
)

type Bool struct {
	sql.NullBool
}

func NewBool(b bool, valid bool) Bool {
	return Bool{
		NullBool: sql.NullBool{
			Bool:  b,
			Valid: valid,
		},
	}
}

func BoolFrom(b bool) Bool {
	return NewBool(b, true)
}

func BoolFromPtr(b *bool) Bool {
	if b == nil {
		return NewBool(false, false)
	}
	return NewBool(*b, true)
}

func (b *Bool) UnmarshalJSON(data []byte) error {
	var err error
	var v any
	err = json.Unmarshal(data, &v)
	switch x := v.(type) {
	case bool:
		b.Bool = x
	case map[string]any:
		err = json.Unmarshal(data, &b.NullBool)
	case nil:
		b.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Bool", reflect.TypeOf(v).Name())
	}
	b.Valid = err == nil
	return err
}

func (b *Bool) UnmarshalText(text []byte) error {
	str := string(text)
	switch str {
	case "", null:
		b.Valid = false
		return nil
	case "true":
		b.Bool = true
	case "false":
		b.Bool = false
	default:
		b.Valid = false
		return errors.New("invalid input:" + str)
	}
	b.Valid = true
	return nil
}

func (b *Bool) MarshalJSON() ([]byte, error) {
	if !b.Valid {
		return []byte(null), nil
	}
	if !b.Bool {
		return []byte("false"), nil
	}
	return []byte("true"), nil
}

func (b *Bool) MarshalText() ([]byte, error) {
	if !b.Valid {
		return []byte{}, nil
	}
	if !b.Bool {
		return []byte("false"), nil
	}
	return []byte("true"), nil
}

func (b *Bool) SetValid(v bool) {
	b.Bool = v
	b.Valid = true
}

func (b *Bool) Ptr() *bool {
	if !b.Valid {
		return nil
	}
	return &b.Bool
}

func (b *Bool) IsZero() bool {
	return !b.Valid
}
