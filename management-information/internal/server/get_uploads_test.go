package server

import (
	"github.com/opg-sirius-supervision-management-information/management-information/internal/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUploads(t *testing.T) {
	client := mockApiClient{
		BondProviders: []model.BondProvider{
			{Id: 1, Name: "Provider A"},
			{Id: 2, Name: "Provider B"},
		},
	}
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
		model.UploadTypes,
		[]model.BondProvider{
			{Id: 1, Name: "Provider A"},
			{Id: 2, Name: "Provider B"},
		},
		appVars,
	}
	assert.Equal(t, expected, ro.data)
}
