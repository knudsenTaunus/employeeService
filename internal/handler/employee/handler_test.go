package employee

import (
	"bytes"
	"github.com/knudsenTaunus/employeeService/internal/config"
	"github.com/knudsenTaunus/employeeService/internal/store"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_Add(t *testing.T) {
	requestBody := []byte(`{"first_name":"test","last_name":"buh","salary":50000,"birthday":"14.07.1988","employee_number":6}`)
	req, err := http.NewRequest("POST", "/employee", bytes.NewBuffer(requestBody))
	assert.NoError(t, err)

	testConfig := &config.Config{
		Sqlitedatabase: struct {
			Path string `yaml:"path"`
		}{Path: "../development.db"},
		Environment: "development",
	}

	testDb, err := store.NewSQLite(testConfig)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := NewHandler(testDb)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, 200, rr.Code)
}
