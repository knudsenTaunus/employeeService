package types

import (
	"encoding/json"
	"io"
)

// Employee is the struct used for response
type EmployeeCars struct {
	ID        int `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	NumberPlate    string `json:"number_plate"`
	Type string `json:"type"`
}

func (e *EmployeeCars) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(e)
}

func (e *EmployeeCars) ToJSON(r io.Writer) error {
	d := json.NewEncoder(r)
	return d.Encode(e)
}

