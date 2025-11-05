package server

import (
	"github.com/opg-sirius-supervision-management-information/internal/model"
	"net/http"
)

type UploadFileVars struct {
	UploadTypes      []model.UploadType
	BondProviders    []model.BondProvider
	ValidationErrors model.ValidationErrors
	AppVars
}

type UploadFileHandler struct {
	router
}

func (h *UploadFileHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	ctx := getContext(r)

	var err error

	bondProviders, err := h.router.Client().GetBondProviders(ctx)
	if err != nil {
		return err
	}
	data := UploadFileVars{
		model.UploadTypes,
		bondProviders,
		nil,
		v}
	data.selectTab("uploads")

	uploadType := model.ParseUploadType(r.PostFormValue("uploadType"))

	if !uploadType.Valid() {
		data.ValidationErrors = model.ValidationErrors{
			"UploadType": map[string]string{"required": "Please select a report to upload"},
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		return h.execute(w, r, data)
	}

	switch uploadType {
	case model.UploadTypeBonds:
		bondProvider := r.PostFormValue("bondProvider")
		if bondProvider == "" {
			data.ValidationErrors = model.ValidationErrors{
				"BondProvider": map[string]string{"required": "Please select a bond provider"},
			}
			w.WriteHeader(http.StatusUnprocessableEntity)
			return h.execute(w, r, data)
		}
	}
	return h.execute(w, r, data)
}
