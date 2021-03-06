package types

import (
	"encoding/json"
	"errors"
	"io"
	"regexp"
	"strings"
	"time"
)

type HandlerEmployee struct {
	EmployeeNumber int      `json:"employee_number"`
	FirstName      string   `json:"first_name"`
	LastName       string   `json:"last_name"`
	Salary         int      `json:"salary"`
	Birthday       JsonDate `json:"birthday"`
	EntryDate      JsonDate `json:"entry_date"`
}

func (e *HandlerEmployee) ToStorageEmployee() (*StorageEmployee) {
	return &StorageEmployee{
		FirstName:      e.FirstName,
		LastName:       e.LastName,
		Salary:         e.Salary,
		Birthday:       e.Birthday.Time,
		EmployeeNumber: e.EmployeeNumber,
		EntryDate:      e.EntryDate.Time,
	}
}

type JsonDate struct {
	Time time.Time
}

func (jd *JsonDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("02.01.2006", s)
	if err != nil {
		return err
	}
	jd.Time = t
	return nil
}

func (jd *JsonDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(jd.Time.Format("01.02.2006"))
}

func (e *HandlerEmployee) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(e)
}

func (e *HandlerEmployee) ToJSON(r io.Writer) error {
	d := json.NewEncoder(r)
	return d.Encode(e)
}

func (e *HandlerEmployee) Validate() error {
	//validate Birthday and EntryDate
	dateRegex := regexp.MustCompile("[0-9]{4}-[0-9]{2}-[0-9]{2}")
	result := dateRegex.FindString(e.EntryDate.Time.String())
	if result == "" {
		return errors.New("No valid date")
	}
	return nil
}