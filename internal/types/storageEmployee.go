package types

import (
	"encoding/json"
	"io"
	"time"
)

type StorageEmployee struct {
	ID             int       `json:"id"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Salary         int       `json:"salary"`
	Birthday       time.Time `json:"birthday"`
	EmployeeNumber int       `json:"employee_number"`
	EntryDate      time.Time `json:"entry_date"`
}

func (e StorageEmployee) ToHandlerEmployee() HandlerEmployee {
	return HandlerEmployee{
		EmployeeNumber: e.EmployeeNumber,
		FirstName:      e.FirstName,
		LastName:       e.LastName,
		Salary:         e.Salary,
		Birthday:       JsonDate{Time: e.Birthday},
		EntryDate:      JsonDate{Time: e.EntryDate},
	}
}

func (e StorageEmployee) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(e)
}

func (e StorageEmployee) ToJSON(r io.Writer) error {
	d := json.NewEncoder(r)
	return d.Encode(e)
}

type StorageEmployees []*StorageEmployee

func (e StorageEmployees) ToHandlerEmployees() []HandlerEmployee {
	result := make([]HandlerEmployee, 0)
	for _, employee := range e {
		result = append(result, employee.ToHandlerEmployee())
	}
	return result
}
