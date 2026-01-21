package server

import (
	"bytes"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/opg-sirius-supervision-management-information/shared"
	"github.com/stretchr/testify/assert"
)

func TestUploadFileHandlerSuccess(t *testing.T) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	_ = writer.WriteField("uploadType", "Bonds")
	_ = writer.WriteField("bondProvider", "1")

	fileWriter, _ := writer.CreateFormFile("fileUpload", "test.csv")
	_, _ = fileWriter.Write([]byte("col1,col2\nval1,val2\n"))
	_ = writer.Close()

	bondProviders := shared.BondProviders{{Id: 1, Name: "Provider1"}}

	client := mockApiClient{BondProviders: bondProviders}
	ro := &mockRoute{client: client}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/uploads", &body)
	r.Header.Add("Content-Type", writer.FormDataContentType())

	appVars := AppVars{
		Path: "/uploads",
	}

	appVars.EnvironmentVars.Prefix = "prefix"
	sut := UploadFileHandler{ro}
	err := sut.render(appVars, w, r)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "prefix/uploads?success=upload", w.Header().Get("HX-Redirect"))
}

func TestUploadFileHandlerValidationErrors(t *testing.T) {
	tests := []struct {
		name            string
		uploadType      string
		bondProvider    string
		fileContent     *string // nil means no file, empty string means empty file
		expectedField   string
		expectedKey     string
		expectedMessage string
	}{
		{
			name:            "MissingUploadType",
			uploadType:      "",
			bondProvider:    "1",
			fileContent:     nil,
			expectedField:   "UploadType",
			expectedKey:     "required",
			expectedMessage: "Please select a report to upload",
		},
		{
			name:            "InvalidBondProvider",
			uploadType:      "Bonds",
			bondProvider:    "invalid",
			fileContent:     nil,
			expectedField:   "BondProvider",
			expectedKey:     "required",
			expectedMessage: "Please select a bond provider",
		},
		{
			name:            "BondProviderNotRecognised",
			uploadType:      "Bonds",
			bondProvider:    "999",
			fileContent:     nil,
			expectedField:   "BondProvider",
			expectedKey:     "required",
			expectedMessage: "Bond provider not recognised",
		},
		{
			name:            "NoFileUploaded",
			uploadType:      "Bonds",
			bondProvider:    "1",
			fileContent:     nil,
			expectedField:   "FileUpload",
			expectedKey:     "required",
			expectedMessage: "No file uploaded",
		},
		{
			name:            "InvalidCSVData",
			uploadType:      "Bonds",
			bondProvider:    "1",
			fileContent:     ptr("invalid\"csv\"data\n\"unclosed"),
			expectedField:   "FileUpload",
			expectedKey:     "invalid",
			expectedMessage: "File does not contain valid CSV data",
		},
		{
			name:            "EmptyCSVFile",
			uploadType:      "Bonds",
			bondProvider:    "1",
			fileContent:     ptr(""),
			expectedField:   "FileUpload",
			expectedKey:     "invalid",
			expectedMessage: "File does not contain valid CSV data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body bytes.Buffer
			writer := multipart.NewWriter(&body)

			_ = writer.WriteField("uploadType", tt.uploadType)
			_ = writer.WriteField("bondProvider", tt.bondProvider)

			if tt.fileContent != nil {
				fileWriter, _ := writer.CreateFormFile("fileUpload", "test.csv")
				_, _ = fileWriter.Write([]byte(*tt.fileContent))
			}

			_ = writer.Close()

			bondProviders := shared.BondProviders{{Id: 1, Name: "Provider1"}}
			client := mockApiClient{BondProviders: bondProviders}
			ro := &mockRoute{client: client}

			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodPost, "/uploads", &body)
			r.Header.Add("Content-Type", writer.FormDataContentType())

			if tt.fileContent != nil {
				ctx := telemetry.ContextWithLogger(r.Context(), slog.New(slog.NewTextHandler(io.Discard, nil)))
				r = r.WithContext(ctx)
			}

			appVars := AppVars{
				Path: "/uploads",
			}

			sut := UploadFileHandler{ro}
			err := sut.render(appVars, w, r)

			assert.Nil(t, err)
			assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

			data := ro.data.(UploadFileVars)
			assert.Equal(t, tt.expectedMessage, data.ValidationErrors[tt.expectedField][tt.expectedKey])
		})
	}
}

func ptr(s string) *string {
	return &s
}
