package server

import (
	"github.com/opg-sirius-supervision-management-information/internal/api"
	"github.com/opg-sirius-supervision-management-information/internal/model"
	"io"
	"net/http"
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

type mockRoute struct { //nolint:golint,unused
	client   ApiClient
	data     any
	executed bool
	lastW    io.Writer
	error
}

func (r *mockRoute) Client() ApiClient { //nolint:golint,unused
	return r.client
}

func (r *mockRoute) execute(w http.ResponseWriter, req *http.Request, data any) error { //nolint:golint,unused
	r.executed = true
	r.lastW = w
	r.data = data
	return r.error
}

type mockApiClient struct {
	User model.User
}

func (m mockApiClient) GetCurrentUserDetails(context api.Context) (model.User, error) {
	return m.User, nil
}
