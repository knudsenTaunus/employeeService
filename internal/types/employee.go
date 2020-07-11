package types

import (
	"encoding/json"
	"io"
)

// Employee is the struct used for response
type Employee struct {
	ID        int `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Salary    int `json:"salary"`
	Birthday string `json:"birthday"`
}

func (e *Employee) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(e)
}

func (e *Employee) ToJSON(r io.Writer) error {
	d := json.NewEncoder(r)
	return d.Encode(e)
}