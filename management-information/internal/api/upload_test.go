package api

import (
	"bytes"
	"encoding/base64"
	"github.com/opg-sirius-supervision-management-information/management-information/internal/mocks"
	"github.com/opg-sirius-supervision-management-information/shared"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func TestUploadSuccess(t *testing.T) {
	logger, mockClient := SetUpTest()

	client, _ := NewApiClient(&mockClient, "http://localhost:3000", logger, "")

	data := shared.Upload{
		UploadType: shared.ParseUploadType("BONDS"),
		Filename:   "test.csv",
		Base64Data: base64.StdEncoding.EncodeToString([]byte("col1, col2\nabc,1")),
	}

	mocks.GetDoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte{})),
		}, nil
	}

	err := client.Upload(getContext(nil), data)
	assert.NoError(t, err)
}

func TestSubmitUploadReturns500Error(t *testing.T) {
	logger, mockClient := SetUpTest()

	client, _ := NewApiClient(&mockClient, "http://localhost:3000", logger, "")

	mocks.GetDoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(bytes.NewReader([]byte{})),
			Request:    req,
		}, nil
	}

	err := client.Upload(getContext(nil), shared.Upload{})

	assert.Equal(t, StatusError{
		Code:   http.StatusInternalServerError,
		URL:    "/uploads",
		Method: http.MethodPost,
	}, err)
}

func TestSubmitUploadReturnsBadRequestError(t *testing.T) {
	logger, mockClient := SetUpTest()

	client, _ := NewApiClient(&mockClient, "http://localhost:3000", logger, "")

	mocks.GetDoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       io.NopCloser(bytes.NewReader([]byte{})),
			Request:    req,
		}, nil
	}

	err := client.Upload(getContext(nil), shared.Upload{})

	assert.Equal(t, StatusError{
		Code:   http.StatusBadRequest,
		URL:    "/uploads",
		Method: http.MethodPost,
	}, err)
}
