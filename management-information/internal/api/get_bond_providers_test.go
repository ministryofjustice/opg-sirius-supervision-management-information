package api

import (
	"bytes"
	"github.com/opg-sirius-supervision-management-information/management-information/internal/mocks"
	"github.com/opg-sirius-supervision-management-information/shared"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
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

	bondProviders, err := client.GetBondProviders(nil)
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

	_, err := client.GetBondProviders(nil)

	assert.Equal(t, ErrUnauthorized.Error(), err.Error())
}
