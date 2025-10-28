package server

import (
	"github.com/opg-sirius-supervision-management-information/internal/api"
	"github.com/opg-sirius-supervision-management-information/internal/model"
	"io"
)

type mockTemplate struct {
	executed         bool
	executedTemplate bool
	lastVars         interface{}
	lastW            io.Writer
	error            error
}

func (m *mockTemplate) Execute(w io.Writer, vars any) error {
	m.executed = true
	m.lastVars = vars
	m.lastW = w
	return m.error
}

func (m *mockTemplate) ExecuteTemplate(w io.Writer, name string, vars any) error {
	m.executedTemplate = true
	m.lastVars = vars
	m.lastW = w
	return m.error
}

type mockApiClient struct {
	Error error
	User  model.User
}

func (m mockApiClient) GetCurrentUserDetails(context api.Context) (model.User, error) {
	return m.User, m.Error
}
