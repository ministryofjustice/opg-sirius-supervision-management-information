package api

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/opg-sirius-supervision-management-information/shared"
	"github.com/stretchr/testify/assert"
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
				UploadType:   shared.UploadTypeBonds,
				Base64Data:   base64.StdEncoding.EncodeToString([]byte("col1, col2\nabc,1")),
				Filename:     "test.csv",
				BondProvider: &shared.BondProvider{Name: "Marsh"},
			},
			fileStorageError:   fmt.Errorf("file storage error"),
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "bonds successful upload",
			upload: shared.Upload{
				UploadType:   shared.UploadTypeBonds,
				Filename:     fmt.Sprintf("Marsh_%s.csv", time.Now().Format("02_01_2006")),
				Base64Data:   base64.StdEncoding.EncodeToString([]byte("col1, col2\nabc,1")),
				BondProvider: &shared.BondProvider{Name: "Marsh"},
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "visits successful upload",
			upload: shared.Upload{
				UploadType: shared.UploadTypeVisits,
				Filename:   fmt.Sprintf("visits_%s.csv", time.Now().Format("02_01_2006")),
				Base64Data: base64.StdEncoding.EncodeToString([]byte("col1, col2\nabc,1")),
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "missing file name error",
			upload: shared.Upload{
				UploadType: shared.UploadTypeVisits,
				Filename:   "",
				Base64Data: base64.StdEncoding.EncodeToString([]byte("col1, col2\nabc,1")),
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
	}
	for _, tt := range tests {
		_, mockClient := SetUpTest()
		mockS3 := &mockFileStorage{
			err: tt.fileStorageError,
		}

		server := NewServer(&mockClient, mockS3, "async-bucket", nil, "")

		var body bytes.Buffer

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
		if tt.expectedFileName != "" {
			assert.Equal(t, tt.expectedFileName, mockS3.fileName)
		}
	}
}
