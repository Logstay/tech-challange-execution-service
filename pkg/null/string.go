package null

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
)

type String struct {
	sql.NullString
}

func StringFrom(s string) String {
	return NewString(s, true)
}

func StringFromPtr(s *string) String {
	if s == nil {
		return NewString("", false)
	}
	return NewString(*s, true)
}

func NewString(s string, valid bool) String {
	return String{
		NullString: sql.NullString{
			String: s,
			Valid:  valid,
		},
	}
}

func (s *String) UnmarshalJSON(data []byte) error {
	var err error
	var v any
	err = json.Unmarshal(data, &v)
	switch x := v.(type) {
	case string:
		s.String = x
	case map[string]any:
		err = json.Unmarshal(data, &s.NullString)
	case nil:
		s.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.String", reflect.TypeOf(v).Name())
	}
	s.Valid = (err == nil) && (s.String != "")
	return err
}

func (s *String) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(s.String)
}

func (s *String) UnmarshalText(text []byte) error {
	s.String = string(text)
	s.Valid = s.String != ""
	return nil
}

func (s *String) SetValid(v string) {
	s.String = v
	s.Valid = true
}

func (s *String) Ptr() *string {
	if !s.Valid {
		return nil
	}
	return &s.String
}

func (s *String) IsZero() bool {
	return !s.Valid
}
