package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/knudsenTaunus/employeeService/internal/server/rest"

	"github.com/stretchr/testify/assert"

	"github.com/rs/zerolog"

	"github.com/knudsenTaunus/employeeService/internal/handler"
	"github.com/knudsenTaunus/employeeService/internal/model"
)

func TestGetAllUsers(t *testing.T) {
	testdb := &MockDatabase{}
	testLogger := zerolog.New(os.Stdout)

	userChan := make(chan model.User)
	user := handler.NewUser(testdb, userChan, testLogger)
	testServer := httptest.NewServer(user)

	defer testServer.Close()
	r, err := http.NewRequest("GET", testServer.URL+"/user", nil)
	assert.NoError(t, err)

	testTimeString := time.Now().UTC().Format(time.RFC3339)
	fmt.Println(testTimeString)
	result, err := time.Parse(time.RFC3339, testTimeString)
	assert.NoError(t, err)
	fmt.Println(result.String())

	rr := httptest.NewRecorder()
	rest.MuxRouter(user).ServeHTTP(rr, r)
	testdb.Mock.On("GetAll").Return([]model.User{
		{
			ID:        "1234",
			FirstName: "test",
			LastName:  "user",
			Nickname:  "tester123",
			Password:  "foo",
			Email:     "foo@bar.de",
			Country:   "Germany",
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
	}, nil)

	resp, err := http.Get(testServer.URL + "/users")
	assert.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	var testResult []model.User
	err = json.Unmarshal(respBody, &testResult)
	assert.NoError(t, err)
	testdb.AssertCalled(t, "GetAll")
	assert.Equal(t, "tester123", testResult[0].Nickname)
}

func TestGetUser(t *testing.T) {
	testdb := &MockDatabase{}
	testLogger := zerolog.New(os.Stdout)

	userChan := make(chan model.User)
	user := handler.NewUser(testdb, userChan, testLogger)

	testServer := httptest.NewServer(user)
	testdb.Mock.On("Get", "1234").Return(model.User{
		ID:        "1234",
		FirstName: "test",
		LastName:  "user",
		Nickname:  "tester123",
		Password:  "foo",
		Email:     "foo@bar.de",
		Country:   "Germany",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil)

	r, err := http.NewRequest("GET", testServer.URL+"/users/1234", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	rest.MuxRouter(user).ServeHTTP(w, r)

	respBody, err := io.ReadAll(w.Body)
	if err != nil {
		log.Fatal(err)
	}

	var testResult model.User
	err = json.Unmarshal(respBody, &testResult)

	assert.NoError(t, err)
	testdb.AssertCalled(t, "Get", "1234")
	assert.Equal(t, "tester123", testResult.Nickname)
}

func TestUpdateUser(t *testing.T) {
	testdb := &MockDatabase{}
	testLogger := zerolog.New(os.Stdout)

	userChan := make(chan model.User)
	user := handler.NewUser(testdb, userChan, testLogger)

	testServer := httptest.NewServer(user)

	updateUser := model.User{
		ID:        "1234",
		FirstName: "test",
		LastName:  "updated",
		Nickname:  "tester123",
		Password:  "foo",
		Email:     "foo@bar.de",
		Country:   "Germany",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	testdb.Mock.On("Update", mock.AnythingOfType("model.User")).Return(updateUser, nil)

	jsonBody := []byte(`{"first_name": "test","last_name": "updated","nickname": "tester123","password": "foo",
    "email": "foo@bar.de",
    "country": "Germany"}`)
	bodyReader := bytes.NewReader(jsonBody)

	r, err := http.NewRequest(http.MethodPatch, testServer.URL+"/users/1234", bodyReader)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	rest.MuxRouter(user).ServeHTTP(w, r)

	respBody, err := io.ReadAll(w.Body)
	if err != nil {
		log.Fatal(err)
	}

	var testResult model.User
	err = json.Unmarshal(respBody, &testResult)

	assert.NoError(t, err)
	testdb.AssertCalled(t, "Update", mock.AnythingOfType("model.User"))
	assert.Equal(t, "updated", testResult.LastName)
}
