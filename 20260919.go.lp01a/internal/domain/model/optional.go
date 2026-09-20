package model

import (
	"bytes"
	"encoding/json"
)

type Optional[T any] struct {
	Set   bool
	Valid bool
	Value T
}

func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	o.Set = true

	if bytes.Equal(data, []byte("null")) {
		o.Valid = false

		var zero T
		o.Value = zero

		return nil
	}

	o.Valid = false
	return json.Unmarshal(data, &o.Valid)
}
