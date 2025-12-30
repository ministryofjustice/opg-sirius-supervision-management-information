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
	"time"
)

func Test_processUpload(t *testing.T) {
	tests := []struct {
		name               string
		upload             any
		fileStorageError   error
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
		},
	}
	for _, tt := range tests {
		mockS3 := &mockFileStorage{
			err: tt.fileStorageError,
		}

		server := NewServer(mockS3, "async-bucket", nil)

		var body bytes.Buffer

		var expectedFileName string
		if tt.expectedStatusCode == http.StatusOK {
			expectedDate := time.Now().Format("02_01_2006")
			expectedFileName = fmt.Sprintf("bonds-without-orders/Marsh_%s.csv", expectedDate)
		}

		_ = json.NewEncoder(&body).Encode(tt.upload)
		ctx := telemetry.ContextWithLogger(context.Background(), telemetry.NewLogger("opg-sirius-management-information"))
		r, _ := http.NewRequestWithContext(ctx, http.MethodPost, "/upload", &body)
		w := httptest.NewRecorder()

		err := server.ProcessDirectUpload(w, r)

		if tt.expectedStatusCode == http.StatusOK {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}

		assert.Equal(t, tt.expectedStatusCode, w.Result().StatusCode)
		if expectedFileName != "" {
			assert.Equal(t, expectedFileName, mockS3.fileName)
		}
	}
}
