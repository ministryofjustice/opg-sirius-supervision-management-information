package api

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/ministryofjustice/opg-go-common/telemetry"
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
		fileStorageError   error
		expectedStatusCode int
		expectedFileName   string
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
			name: "file storage error",
			upload: shared.Upload{
				UploadType: shared.UploadTypeUnknown,
				Base64Data: base64.StdEncoding.EncodeToString([]byte("col1, col2\nabc,1")),
				Filename:   "test.csv",
			},
			fileStorageError:   fmt.Errorf("Oops!"),
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "pass",
			upload: shared.Upload{
				UploadType:   shared.UploadTypeBonds,
				Filename:     "data.csv",
				Base64Data:   base64.StdEncoding.EncodeToString([]byte("col1, col2\nabc,1")),
				BondProvider: shared.BondProvider{Name: "Marsh"},
			},
			expectedStatusCode: http.StatusOK,
			//expectedFileName:   "bonds-without-orders/Marsh_15_12_2025.csv", //Todo - fix so it isn't impacted by current date
		},
	}
	for _, tt := range tests {
		mockS3 := &mockFileStorage{
			err: tt.fileStorageError,
		}

		server := NewServer(mockS3, "async-bucket")

		var body bytes.Buffer

		_ = json.NewEncoder(&body).Encode(tt.upload)
		ctx := telemetry.ContextWithLogger(context.Background(), telemetry.NewLogger("opg-sirius-management-information"))
		r, _ := http.NewRequestWithContext(ctx, http.MethodPost, "/upload", &body)
		w := httptest.NewRecorder()

		server.ProcessDirectUpload(w, r)
		assert.Equal(t, tt.expectedStatusCode, w.Result().StatusCode)
		if tt.expectedFileName != "" {
			assert.Equal(t, tt.expectedFileName, mockS3.fileName)
		}
	}
}
