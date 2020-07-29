package types

import (
	"encoding/json"
	"errors"
	"io"
	"regexp"
	"time"
)

// StorageEmployee is the struct used for response
type StorageEmployee struct {
	ID             int      `json:"id"`
	FirstName      string   `json:"first_name"`
	LastName       string   `json:"last_name"`
	Salary         int      `json:"salary"`
	Birthday       time.Time `json:"birthday"`
	EmployeeNumber int      `json:"employee_number"`
	EntryDate      time.Time `json:"entry_date"`
}

func (e *StorageEmployee) ToHandlerEmployee() (*HandlerEmployee) {
	return &HandlerEmployee{
		EmployeeNumber: e.EmployeeNumber,
		FirstName:      e.FirstName,
		LastName:       e.LastName,
		Salary:         e.Salary,
		Birthday:       JsonDate{Time: e.Birthday},
		EntryDate:      JsonDate{Time: e.EntryDate},
	}
}

func (e *StorageEmployee) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(e)
}

func (e *StorageEmployee) ToJSON(r io.Writer) error {
	d := json.NewEncoder(r)
	return d.Encode(e)
}

func (e *StorageEmployee) Validate() error {
	//validate Birthday and EntryDate
	dateRegex := regexp.MustCompile("[0-9]{4}-[0-9]{2}-[0-9]{2}")
	result := dateRegex.FindString(e.EntryDate.String())
	if result == "" {
		return errors.New("No valid date")
	}
	return nil
}