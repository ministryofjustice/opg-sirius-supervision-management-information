package api

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/opg-sirius-supervision-management-information/management-information/internal/auth"
	"github.com/opg-sirius-supervision-management-information/management-information/internal/mocks"
	"github.com/opg-sirius-supervision-management-information/shared"
	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
)

func TestGetBondProviders(t *testing.T) {
	logger, mockClient := SetUpTest()
	mockJwtClient := &mockJWTClient{}
	client, _ := NewApiClient(&mockClient, mockJwtClient, "http://localhost:3000", logger, "")

	json := `[{
			  "id": 1,
			  "name": "Marsh"
			},
			{
			  "id": 2,
			  "name": "Howden"
			},
			{
			  "id": 3,
			  "name": "DBS"
			}]`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       r,
		}, nil
	}

	expectedResponse := shared.BondProviders{
		{Id: 1, Name: "Marsh"},
		{Id: 2, Name: "Howden"},
		{Id: 3, Name: "DBS"},
	}

	ctx := auth.Context{
		User:    &shared.User{ID: 123},
		Context: context.Background(),
	}

	bondProviders, err := client.GetBondProviders(ctx)
	assert.NoError(t, err)

	assert.Equal(t, expectedResponse, bondProviders)
}

func TestGetBondProvidersUnauthorised(t *testing.T) {
	logger, mockClient := SetUpTest()
	mockJwtClient := &mockJWTClient{}

	client, _ := NewApiClient(&mockClient, mockJwtClient, "http://localhost:3000", logger, "")

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       io.NopCloser(bytes.NewReader([]byte{})),
		}, nil
	}

	ctx := auth.Context{
		User:    &shared.User{ID: 123},
		Context: context.Background(),
	}

	_, err := client.GetBondProviders(ctx)

	assert.Equal(t, ErrUnauthorized.Error(), err.Error())
}

func TestGetBondProviders_contract(t *testing.T) {
	pact, err := consumer.NewV2Pact(consumer.MockHTTPProviderConfig{
		Consumer: "sirius-supervision-management-information",
		Provider: "sirius",
		LogDir:   "../../../logs",
		PactDir:  "../../../pacts",
	})
	assert.NoError(t, err)

	err = pact.
		AddInteraction().
		Given("Bond providers exist").
		UponReceiving("A request for the bond providers list").
		WithRequest("GET", "/supervision-api/v1/bond-providers", func(b *consumer.V2RequestBuilder) {
			b.Header("Accept", matchers.S("application/json"))
		}).
		WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
			b.Header("Content-Type", matchers.S("application/json"))
			b.JSONBody(matchers.EachLike(matchers.MapMatcher{
				"id":   matchers.Like(1),
				"name": matchers.Like("Marsh"),
			}, 3))
		}).
		ExecuteTest(t, func(config consumer.MockServerConfig) error {
			logger, mockClient := SetUpTest()
			mockJwtClient := &mockJWTClient{}

			client, _ := NewApiClient(&mockClient, mockJwtClient, fmt.Sprintf("http://%s:%d", config.Host, config.Port), logger, "")

			ctx := auth.Context{
				User:    &shared.User{ID: 123},
				Context: context.Background(),
			}

			bondProviders, _ := client.GetBondProviders(ctx)

			assert.NotEmpty(t, bondProviders)
			assert.EqualValues(t, shared.BondProvider{Id: 1, Name: "Marsh"}, bondProviders[0])

			return nil
		})

	assert.NoError(t, err)
}
