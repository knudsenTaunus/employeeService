package model

import (
	"encoding/json"
	"io"
)

type Car struct {
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	NumberPlate    string `json:"number_plate"`
	Type           string `json:"type"`
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

type StorageCar struct {
	Manufacturer   string `json:"manufacturer"`
	Type           string `json:"type"`
	NumberPlate    string `json:"number_plate"`
	EmployeeNumber int    `json:"employee_number"`
}

func (e *StorageCar) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(e)
}

func (e *StorageCar) ToJSON(r io.Writer) error {
	d := json.NewEncoder(r)
	return d.Encode(e)
}
