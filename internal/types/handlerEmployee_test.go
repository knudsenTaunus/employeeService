package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestEmployee_Validate(t *testing.T) {
	testEmployee := &HandlerEmployee{
		FirstName:      "Max",
		LastName:       "Mustermann",
		Salary:         80000,
		Birthday:       JsonDate{Time: time.Date(1984,02,24,0,0,0,0,time.UTC)},
		EmployeeNumber: 5,
		EntryDate:      JsonDate{Time: time.Date(2020,05,0,0,0,0,0,time.UTC)},
	}
	assert.NoError(t, testEmployee.Validate())
}
