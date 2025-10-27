package server

import (
	"errors"
	"github.com/opg-sirius-supervision-management-information/internal/api"
	"github.com/opg-sirius-supervision-management-information/internal/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type mockRouteData struct {
	stuff string
	AppVars
}

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
					Id:          0,
					Name:        "",
					PhoneNumber: "",
					Deleted:     false,
					Email:       "",
					Firstname:   "",
					Surname:     "",
					Roles:       test.userRoles,
					Locked:      false,
					Suspended:   false,
				},
			}
			template := &mockTemplate{}

			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodGet, "", nil)
			r.Header.Add("HX-Request", "true")

			data := mockRouteData{
				stuff:   "abc",
				AppVars: AppVars{Path: "/path"},
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
		{userRoles: []string{"Reporting User", "Case Manager"}, expectError: false},
		{userRoles: []string{"System Admin", "Reporting User", "Case Manager"}, expectError: false},
		{userRoles: []string{""}, expectError: true},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i), func(t *testing.T) {
			client := mockApiClient{}
			template := &mockTemplate{}

			w := httptest.NewRecorder()
			ctx := api.Context{}
			r, _ := http.NewRequestWithContext(ctx, http.MethodGet, "", nil)

			client.User = model.User{
				Id:    123,
				Roles: test.userRoles,
			}

			data := PageData{
				Data: mockRouteData{
					stuff:   "abc",
					AppVars: AppVars{Path: "/path/"},
				},
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

//func TestRoute_error(t *testing.T) {
//	client := mockApiClient{}
//	client.error = errors.New("it broke")
//	template := &mockTemplate{}
//
//	w := httptest.NewRecorder()
//	ctx := {}
//	r, _ := http.NewRequestWithContext(ctx, http.MethodGet, "", nil)
//	r.SetPathValue("clientId", "abc")
//
//	data := PageData{
//		Data: mockRouteData{
//			stuff:   "abc",
//			AppVars: AppVars{Path: "/path/"},
//		},
//	}
//
//	sut := route{client: client, tmpl: template, partial: "test"}
//
//	err := sut.execute(w, r, data.Data)
//
//	assert.NotNil(t, err)
//	assert.Equal(t, "it broke", err.Error())
//}
