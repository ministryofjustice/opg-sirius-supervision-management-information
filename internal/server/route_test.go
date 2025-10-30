package server

import (
	"context"
	"errors"
	"github.com/opg-sirius-supervision-management-information/internal/api"
	"github.com/opg-sirius-supervision-management-information/internal/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestRoute_htmxRequest_withPermissions(t *testing.T) {
	tests := []struct {
		userRoles []string
	}{
		{userRoles: []string{"Case Manager"}},
		{userRoles: []string{"Reporting User", "Case Manager"}},
		{userRoles: []string{"System Admin", "Reporting User", "Case Manager"}},
		{userRoles: []string{""}},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i), func(t *testing.T) {
			client := mockApiClient{
				User: model.User{
					Id:      123,
					Name:    "Test User",
					Deleted: false,
					Email:   "test@user.com",
					Roles:   test.userRoles,
				},
			}
			template := &mockTemplate{}

			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodGet, "", nil)
			r.Header.Add("HX-Request", "true")

			data := PageData{
				Data: "test data",
			}

			sut := route{client: client, tmpl: template, partial: "test"}

			err := sut.execute(w, r, data)

			assert.Nil(t, err)
			assert.True(t, template.executedTemplate)
			assert.False(t, template.executed)
			assert.Equal(t, data, template.lastVars)
		})
	}

}

func TestRoute_fullPage_with_reportingUserPermissions(t *testing.T) {
	tests := []struct {
		userRoles   []string
		expectError bool
	}{
		{userRoles: []string{"Case Manager"}, expectError: true},
		{userRoles: []string{"Reporting User"}, expectError: false},
		{userRoles: []string{"Reporting User", "Case Manager"}, expectError: false},
		{userRoles: []string{"System Admin", "Reporting User", "Case Manager"}, expectError: false},
		{userRoles: []string{"System Admin", "Case Manager", "Reporting User"}, expectError: false},
		{userRoles: []string{""}, expectError: true},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i), func(t *testing.T) {
			client := mockApiClient{
				User: model.User{
					Id:        123,
					Name:      "Test User",
					Email:     "test@user.com",
					Firstname: "Test",
					Surname:   "User",
					Roles:     test.userRoles,
				},
			}
			template := &mockTemplate{}

			w := httptest.NewRecorder()
			ctx := api.Context{
				Context: context.Background(),
			}
			r, _ := http.NewRequestWithContext(ctx, http.MethodGet, "", nil)

			data := PageData{
				Data: "test data",
			}

			sut := route{client: client, tmpl: template, partial: "test"}
			err := sut.execute(w, r, data.Data)

			if test.expectError {
				assert.Equal(t, errors.New("not reporting user"), err)
				assert.False(t, template.executed)
				assert.False(t, template.executedTemplate)
				assert.Nil(t, template.lastVars)
			} else {
				assert.Nil(t, err)
				assert.True(t, template.executed)
				assert.False(t, template.executedTemplate)
				assert.Equal(t, data, template.lastVars)
			}
		})
	}
}

func TestRoute_error(t *testing.T) {
	client := mockApiClient{
		User: model.User{
			Id:    123,
			Name:  "Test User",
			Roles: []string{"Reporting User"},
		},
		Error: errors.New("it broke"),
	}
	template := &mockTemplate{}
	w := httptest.NewRecorder()
	ctx := api.Context{
		Context: context.Background(),
	}
	r, _ := http.NewRequestWithContext(ctx, http.MethodGet, "", nil)

	data := PageData{
		Data: "test Data",
	}

	sut := route{client: client, tmpl: template, partial: "test"}
	err := sut.execute(w, r, data.Data)
	assert.NotNil(t, err)
	assert.Equal(t, "it broke", err.Error())
}
