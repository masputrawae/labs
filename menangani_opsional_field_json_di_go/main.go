package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Optional[T any] struct {
	Value T
	Set   bool
	Null  bool
}

func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	o.Set = true

	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		o.Null = true
		var zero T
		o.Value = zero
		return nil
	}

	o.Null = false
	return json.Unmarshal(data, &o.Value)
}

type User struct {
	Name     Optional[string] `json:"name"`
	Age      Optional[int]    `json:"age"`
	IsActive Optional[bool]   `json:"is_active"`
}

var dummy = []string{
	`{ "name": "", "age": 0, "is_active": false }`,
	`{ "age": 20 }`,
	`{ }`,
}

func main() {
	for _, d := range dummy {
		var data User
		_ = json.Unmarshal([]byte(d), &data)
		j, _ := json.MarshalIndent(data, "", "  ")
		fmt.Println(string(j))
	}
}
