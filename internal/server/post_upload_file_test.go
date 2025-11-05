package server

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestUploadFileHandlerSuccess(t *testing.T) {
	form := url.Values{
		"uploadType":   {"Bonds"},
		"bondProvider": {"1"},
	}

	client := mockApiClient{}
	ro := &mockRoute{client: client}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/uploads", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	appVars := AppVars{
		Path: "/uploads",
	}

	appVars.EnvironmentVars.Prefix = "prefix"
	sut := UploadFileHandler{ro}
	err := sut.render(appVars, w, r)
	assert.Nil(t, err)
}

func TestUploadFileHandlerValidationErrors(t *testing.T) {
	form := url.Values{
		"uploadType":   {""},
		"bondProvider": {""},
	}

	client := mockApiClient{}
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
