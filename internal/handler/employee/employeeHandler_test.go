package employee

import (
	"bytes"
	"github.com/knudsenTaunus/employeeService/internal/storage/development"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_Add(t *testing.T) {
	requestBody := []byte(`{"first_name":"test","last_name":"buh","salary":50000,"birthday":"14.07.1988","employee_number":6}`)
	req, err := http.NewRequest("POST", "/employee", bytes.NewBuffer(requestBody))
	assert.NoError(t, err)

	testDb, err := development.New()
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := New(testDb)

	handler.Add().ServeHTTP(rr, req)
	assert.Equal(t, 200, rr.Code)
}
