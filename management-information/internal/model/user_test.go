package model

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func Test_isReportingUser(t *testing.T) {
	tests := []struct {
		user             User
		expectedResponse bool
	}{
		{
			user: User{
				Id:    1,
				Name:  "Reporting User",
				Roles: []string{"Reporting User", "Case Manager"},
			},
			expectedResponse: true,
		},
		{
			user: User{
				Id:    2,
				Name:  "Non Reporting User",
				Roles: []string{"Case Manager"},
			},
			expectedResponse: false,
		},
		{
			user: User{
				Id:    3,
				Name:  "System Admin",
				Roles: []string{"System Admin", "Case Manager", "Reporting User"},
			},
			expectedResponse: true,
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i), func(t *testing.T) {
			assert.Equal(t, test.expectedResponse, test.user.IsReportingUser())
		})
	}
}
