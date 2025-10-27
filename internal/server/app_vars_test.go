package server

import (
	"github.com/opg-sirius-supervision-management-information/internal/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestNewAppVars(t *testing.T) {
	r, _ := http.NewRequest("GET", "/path", nil)
	r.SetPathValue("clientId", "1")
	r.AddCookie(&http.Cookie{Name: "XSRF-TOKEN", Value: "abc123"})

	envVars := EnvironmentVars{}
	client := mockApiClient{
		User: model.User{
			Id:          1,
			Name:        "Reporting User",
			PhoneNumber: "123456",
			Deleted:     false,
			Email:       "reporting.user@email.com",
			Firstname:   "Reporting",
			Surname:     "User",
			Roles:       []string{"Reporting User", "Case Manager"},
			Locked:      false,
			Suspended:   false,
		},
	}
	vars := NewAppVars(client, r, envVars)

	assert.Equal(t, AppVars{
		Path:            "/path",
		XSRFToken:       "abc123",
		EnvironmentVars: envVars,
		Tabs: []Tab{
			{
				Id:    "downloads",
				Title: "Downloads",
			},
			{
				Id:    "uploads",
				Title: "Uploads",
			},
		},
		User: model.User{
			Id:          1,
			Name:        "Reporting User",
			PhoneNumber: "123456",
			Deleted:     false,
			Email:       "reporting.user@email.com",
			Firstname:   "Reporting",
			Surname:     "User",
			Roles:       []string{"Reporting User", "Case Manager"},
			Locked:      false,
			Suspended:   false,
		},
	}, vars)
}
