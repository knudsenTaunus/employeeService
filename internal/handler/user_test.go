package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/knudsenTaunus/employeeService/internal/model"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_Add(t *testing.T) {
	requestBody := []byte(`{"first_name": "Jan-Philipp","last_name": "Heinrich","nickname": "knudsenTaunus","password": "foo","email": "bar@barbara.ru","country": "UK"}`)
	req, err := http.NewRequest("POST", "/user", bytes.NewBuffer(requestBody))
	assert.NoError(t, err)

	testUser := model.User{
		FirstName: "Jan-Philipp",
		LastName:  "Heinrich",
		Nickname:  "knudsenTaunus",
		Password:  "foo",
		Email:     "bar@barbara.ru",
		Country:   "UK",
	}

	tdb := &TestDatabase{}

	tdb.Mock.On("Create", testUser).Return(model.User{
		ID:        "1234",
		FirstName: "Jan-Philipp",
		LastName:  "Heinrich",
		Nickname:  "knudsenTaunus",
		Password:  "foo",
		Email:     "bar@barbara.ru",
		Country:   "UK",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}, nil)

	rr := httptest.NewRecorder()
	testChan := make(chan model.User)
	testlogger := zerolog.New(os.Stdout)
	handler := NewUser(tdb, testChan, testlogger)

	go func() {
		result := <-testChan
		testlogger.Info().Msg(result.FirstName)
	}()

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	tdb.AssertCalled(t, "Create", testUser)
}

type TestDatabase struct {
	mock.Mock
}

func (tdb *TestDatabase) GetAll() ([]model.User, error) {
	return nil, nil
}

func (tdb *TestDatabase) GetPaginatedAndFiltered(page, pageSize int, filter string) ([]model.User, error) {
	args := tdb.Called(page, pageSize, filter)
	return args[0].([]model.User), args.Error(1)
}

func (tdb *TestDatabase) Get(id string) (model.User, error) {
	args := tdb.Called(id)
	return args[0].(model.User), args.Error(1)
}
func (tdb *TestDatabase) Create(user model.User) (model.User, error) {
	args := tdb.Called(user)
	return args.Get(0).(model.User), args.Error(1)
}
func (tdb *TestDatabase) Delete(id string) error {
	args := tdb.Called(id)
	return args.Error(0)
}
func (tdb *TestDatabase) Update(user model.User) (model.User, error) {
	args := tdb.Called(user)
	return args.Get(0).(model.User), args.Error(1)
}

func (tdb *TestDatabase) FindAllEmployees() ([]model.User, error) {
	return nil, nil
}
