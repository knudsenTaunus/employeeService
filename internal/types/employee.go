package types

import (
	"encoding/json"
	"io"
	"strings"
	"time"
)

// Employee is the struct used for response
type Employee struct {
	ID        int `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Salary    int `json:"salary"`
	Birthday JsonDate `json:"birthday"`
	EmployeeNumber int `json:"employee_number"`
	EntryDate JsonDate `json:"entry_date"`
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

func (e *Employee) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(e)
}

func (e *Employee) ToJSON(r io.Writer) error {
	d := json.NewEncoder(r)
	return d.Encode(e)
}