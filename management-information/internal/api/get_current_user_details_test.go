package api

import (
	"bytes"
	"context"
	"github.com/opg-sirius-supervision-management-information/management-information/internal/auth"
	"github.com/opg-sirius-supervision-management-information/management-information/internal/mocks"
	"github.com/opg-sirius-supervision-management-information/shared"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetCurrentUserDetails(t *testing.T) {
	logger, mockClient := SetUpTest()
	mockJwtClient := &mockJWTClient{}
	client, _ := NewApiClient(&mockClient, mockJwtClient, "http://localhost:3000", logger, "")

	json := `{
			   "id":65,
			   "name":"case",
			   "phoneNumber":"12345678",
			   "teams":[{
				  "displayName":"Lay Team 1 - (Supervision)",
				  "id":13
			   }],
			   "displayName":"case manager",
			   "deleted":false,
			   "email":"case.manager@opgtest.com",
			   "firstname":"case",
			   "surname":"manager",
			   "roles":[
				  "Case Manager"
			   ],
			   "locked":false,
			   "suspended":false
			}`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := shared.User{
		ID:          65,
		DisplayName: "case manager",
		Roles:       []string{"Case Manager"},
	}

	ctx := auth.Context{
		User:    &shared.User{ID: 123},
		Context: context.Background(),
	}

	teams, err := client.GetCurrentUserDetails(ctx)
	assert.Equal(t, expectedResponse, teams)
	assert.Equal(t, nil, err)
}

func TestGetCurrentUserDetailsReturnsUnauthorisedClientError(t *testing.T) {
	logger, _ := SetUpTest()
	mockJwtClient := &mockJWTClient{}

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	ctx := auth.Context{
		User:    &shared.User{ID: 123},
		Context: context.Background(),
	}

	client, _ := NewApiClient(http.DefaultClient, mockJwtClient, svr.URL, logger, "")
	_, err := client.GetCurrentUserDetails(ctx)
	assert.Equal(t, ErrUnauthorized, err)
}

func TestMyDetailsReturns500Error(t *testing.T) {
	logger, _ := SetUpTest()
	mockJwtClient := &mockJWTClient{}

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer svr.Close()

	ctx := auth.Context{
		User:    &shared.User{ID: 123},
		Context: context.Background(),
	}

	client, _ := NewApiClient(http.DefaultClient, mockJwtClient, svr.URL, logger, "")

	_, err := client.GetCurrentUserDetails(ctx)
	assert.Equal(t, StatusError{
		Code:   http.StatusInternalServerError,
		URL:    svr.URL + "/v1/users/current",
		Method: http.MethodGet,
	}, err)
}

func TestMyDetailsReturns200(t *testing.T) {
	logger, mockClient := SetUpTest()
	mockJwtClient := &mockJWTClient{}

	client, _ := NewApiClient(&mockClient, mockJwtClient, "http://localhost:3000", logger, "")

	json := `{
		"id": 55,
		"name": "case",
		"phoneNumber": "12345678",
		"teams": [],
		"displayName": "case manager",
		"deleted": false,
		"email": "case.manager@opgtest.com",
		"firstname": "case",
		"surname": "manager",
		"roles": [
			"OPG User",
			"Case Manager"
		],
		"locked": false,
		"suspended": false
    }`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := shared.User{
		ID:          55,
		DisplayName: "case manager",
		Roles:       []string{"OPG User", "Case Manager"},
	}

	ctx := auth.Context{
		User:    &shared.User{ID: 123},
		Context: context.Background(),
	}

	user, err := client.GetCurrentUserDetails(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, user, expectedResponse)
}
