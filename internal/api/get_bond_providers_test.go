package api

import (
	"bytes"
	"github.com/opg-sirius-supervision-management-information/internal/mocks"
	"github.com/opg-sirius-supervision-management-information/internal/model"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func TestGetBondProviders(t *testing.T) {
	logger, mockClient := SetUpTest()
	client, _ := NewApiClient(&mockClient, "http://localhost:3000", logger)

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

	expectedResponse := []model.BondProvider{
		{Id: 1, Name: "Marsh"},
		{Id: 2, Name: "Howden"},
		{Id: 3, Name: "DBS"},
	}

	bondProviders, err := client.GetBondProviders(getContext(nil))
	assert.NoError(t, err)

	assert.Equal(t, expectedResponse, bondProviders)
}

func TestGetBondProvidersUnauthorised(t *testing.T) {
	logger, mockClient := SetUpTest()
	client, _ := NewApiClient(&mockClient, "http://localhost:3000", logger)

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       io.NopCloser(bytes.NewReader([]byte{})),
		}, nil
	}

	_, err := client.GetBondProviders(getContext(nil))

	assert.Equal(t, ErrUnauthorized.Error(), err.Error())
}
