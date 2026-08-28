package api

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/opg-sirius-supervision-management-information/management-information/internal/auth"
	"github.com/opg-sirius-supervision-management-information/management-information/internal/mocks"
	"github.com/opg-sirius-supervision-management-information/shared"
	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
)

func TestGetCurrentUserDetails(t *testing.T) {
	logger, mockClient := SetUpTest()
	mockJwtClient := &mockJWTClient{}
	client := NewApiClient(&mockClient, mockJwtClient, "http://localhost:3000", logger, "")

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

	client := NewApiClient(http.DefaultClient, mockJwtClient, svr.URL, logger, "")
	_, err := client.GetCurrentUserDetails(ctx)
	assert.Equal(t, ErrUnauthorized, err)
}

func TestGetCurrentUserDetailsReturns500Error(t *testing.T) {
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

	client := NewApiClient(http.DefaultClient, mockJwtClient, svr.URL, logger, "")

	_, err := client.GetCurrentUserDetails(ctx)
	assert.Equal(t, StatusError{
		Code:   http.StatusInternalServerError,
		URL:    svr.URL + "/v1/users/current",
		Method: http.MethodGet,
	}, err)
}

func TestGetCurrentUserDetailsReturns200(t *testing.T) {
	logger, mockClient := SetUpTest()
	mockJwtClient := &mockJWTClient{}

	client := NewApiClient(&mockClient, mockJwtClient, "http://localhost:3000", logger, "")

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

func TestGetCurrentUserDetails_contract(t *testing.T) {
	pact, err := consumer.NewV2Pact(consumer.MockHTTPProviderConfig{
		Consumer: "sirius-supervision-management-information",
		Provider: "sirius",
		LogDir:   "../../../logs",
		PactDir:  "../../../pacts",
	})
	assert.NoError(t, err)

	err = pact.
		AddInteraction().
		Given("User exists").
		UponReceiving("A request for the current user").
		WithRequest("GET", "/supervision-api/v1/users/current", func(b *consumer.V2RequestBuilder) {
			b.Header("Accept", matchers.S("application/json"))
		}).
		WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
			b.Header("Content-Type", matchers.S("application/json"))
			b.JSONBody(matchers.MapMatcher{
				"id":          matchers.Like(1),
				"displayName": matchers.Like("Colin Case"),
				"roles":       matchers.EachLike("Case Manager", 1),
			})
		}).
		ExecuteTest(t, func(config consumer.MockServerConfig) error {
			client := NewApiClient(
				http.DefaultClient,
				nil,
				fmt.Sprintf("http://%s:%d/supervision-api", config.Host, config.Port),
				telemetry.NewLogger("opg-sirius-management-information"),
				"",
			)
			ctx := auth.Context{
				User:    &shared.User{ID: 123},
				Context: context.Background(),
			}

			user, err := client.GetCurrentUserDetails(ctx)
			assert.NoError(t, err)

			assert.EqualValues(t, shared.User{
				ID:          1,
				DisplayName: "Colin Case",
				Roles:       []string{"Case Manager"},
			}, user)
			return nil
		})

	assert.NoError(t, err)
}
