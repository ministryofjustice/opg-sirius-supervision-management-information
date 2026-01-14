package server

import (
	"bytes"
	"github.com/opg-sirius-supervision-management-information/shared"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
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
	form := url.Values{
		"uploadType":   {""},
		"bondProvider": {"1"},
	}

	bondProviders := shared.BondProviders{{Id: 1, Name: "Provider1"}}

	client := mockApiClient{BondProviders: bondProviders}
	ro := &mockRoute{client: client}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/uploads", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	appVars := AppVars{
		Path: "/uploads",
	}

	sut := UploadFileHandler{ro}
	err := sut.render(appVars, w, r)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}
