package api

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/opg-sirius-supervision-management-information/shared"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_processUpload(t *testing.T) {
	tests := []struct {
		name               string
		upload             any
		serverError        error
		expectedStatusCode int
	}{
		{
			name: "base64 decode error",
			upload: shared.Upload{
				UploadType: shared.UploadTypeUnknown,
				Base64Data: "Hey! This is not base64!",
				Filename:   "oops",
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:               "json decode error",
			upload:             2,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "server error",
			upload: shared.Upload{
				UploadType: shared.UploadTypeUnknown,
				Base64Data: base64.StdEncoding.EncodeToString([]byte("col1, col2\nabc,1")),
				Filename:   "test.csv",
			},
			serverError:        fmt.Errorf("Oops!"),
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "pass",
			upload: shared.Upload{
				UploadType: shared.UploadTypeBonds,
				Filename:   "data.csv",
				Base64Data: base64.StdEncoding.EncodeToString([]byte("col1, col2\nabc,1")),
			},
			expectedStatusCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		service := &mockService{err: tt.serverError}

		server := NewServer(service)

		var body bytes.Buffer

		_ = json.NewEncoder(&body).Encode(tt.upload)
		r := httptest.NewRequest(http.MethodPost, "/upload", &body)
		ctx := context.Background()
		r = r.WithContext(ctx)
		w := httptest.NewRecorder()

		server.ProcessDirectUpload(w, r)
		assert.Equal(t, tt.expectedStatusCode, w.Result().StatusCode)
	}
}
