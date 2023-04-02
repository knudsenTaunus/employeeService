package integration

import (
	"github.com/knudsenTaunus/employeeService/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) GetAll() ([]model.User, error) {
	args := m.Called()
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockDatabase) GetPaginatedAndFiltered(page, pageSize int, filter string) ([]model.User, error) {
	args := m.Called(page, pageSize, filter)
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockDatabase) Get(id string) (model.User, error) {
	args := m.Called(id)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockDatabase) Create(user model.User) (model.User, error) {
	args := m.Called(user)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockDatabase) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDatabase) Update(user model.User) (model.User, error) {
	args := m.Called(user)
	return args.Get(0).(model.User), args.Error(1)
}
