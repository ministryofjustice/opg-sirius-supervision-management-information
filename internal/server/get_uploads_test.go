package server

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUploads(t *testing.T) {
	client := mockApiClient{}
	ro := &mockRoute{client: client}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/uploads", nil)

	appVars := AppVars{
		Path: "/downloads",
		Tabs: []Tab{
			{
				Id:    "downloads",
				Title: "Downloads",
			},
			{
				Id:    "uploads",
				Title: "Uploads",
			},
		},
	}

	sut := GetUploadsHandler{ro}
	err := sut.render(appVars, w, r)

	assert.Nil(t, err)
	assert.True(t, ro.executed)

	expected := UploadsVars{
		appVars,
	}
	assert.Equal(t, expected, ro.data)
}
