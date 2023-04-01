package handler

import (
	"bytes"
	"github.com/knudsenTaunus/employeeService/internal/model"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

type TestDatabase struct {
	mock.Mock
}

func (tdb *TestDatabase) FindAllEmployees() ([]model.Employee, error) {
	return nil, nil
}
func (tdb *TestDatabase) FindAllEmployeesLimit(limit string) ([]model.Employee, error) {
	args := tdb.Called(limit)
	return args.Get(0).([]model.Employee), args.Error(0)
}
func (tdb *TestDatabase) FindEmployee(id string) (model.Employee, error) {
	args := tdb.Called(id)
	return args.Get(0).(model.Employee), args.Error(0)
}

func (tdb *TestDatabase) AddEmployee(employee model.Employee) error {
	args := tdb.Called(employee)
	return args.Error(0)

}
func (tdb *TestDatabase) RemoveEmployee(id string) error {
	args := tdb.Called(id)
	return args.Error(0)
}

func TestHandler_Add(t *testing.T) {
	requestBody := []byte(`{"first_name":"test","last_name":"buh","salary":50000,"birthday":"14.07.1988","entry_date":"14.07.2003","employee_number":6}`)
	req, err := http.NewRequest("POST", "/employee", bytes.NewBuffer(requestBody))
	assert.NoError(t, err)

	testEmployee := model.Employee{
		EmployeeNumber: 6,
		FirstName:      "test",
		LastName:       "buh",
		Salary:         50000,
		Birthday:       model.JsonDate{Time: time.Date(1988, 7, 14, 0, 0, 0, 0, time.UTC)},
		EntryDate:      model.JsonDate{Time: time.Date(2003, 7, 14, 0, 0, 0, 0, time.UTC)},
	}

	tdb := &TestDatabase{}

	tdb.Mock.On("AddEmployee", testEmployee).Return(nil)

	rr := httptest.NewRecorder()
	handler := NewEmployee(tdb, zerolog.New(os.Stdout))

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusCreated, rr.Code)
	tdb.AssertCalled(t, "AddEmployee", testEmployee)
}
