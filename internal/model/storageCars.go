package model

import (
	"encoding/json"
	"io"
)

type Car struct {
	ID             int    `json:"id"`
	Manufacturer   string `json:"manufacturer"`
	Type           string `json:"type"`
	NumberPlate    string `json:"number_plate"`
	EmployeeNumber int    `json:"employee_number"`
}

func (e *Car) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(e)
}

func (e *Car) ToJSON(r io.Writer) error {
	d := json.NewEncoder(r)
	return d.Encode(e)
}
