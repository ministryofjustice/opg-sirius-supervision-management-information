package server

import (
	"encoding/base64"
	"fmt"
	"github.com/opg-sirius-supervision-management-information/management-information/internal/model"
	"github.com/opg-sirius-supervision-management-information/shared"
	"io"
	"net/http"
	"strconv"
)

type UploadFileVars struct {
	UploadTypes      []shared.UploadType
	BondProviders    []shared.BondProvider
	ValidationErrors model.ValidationErrors
	AppVars
}

type UploadFileHandler struct {
	router
}

func (h *UploadFileHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	ctx := getContext(r)

	bondProviders, err := h.router.Client().GetBondProviders(ctx)
	if err != nil {
		return err
	}

	data := UploadFileVars{
		shared.UploadTypes,
		bondProviders,
		nil,
		v}
	data.selectTab("uploads")

	switch shared.ParseUploadType(r.PostFormValue("uploadType")) {
	case shared.UploadTypeBonds:
		bondProviderId, err := strconv.Atoi(r.PostFormValue("bondProvider"))

		if err != nil {
			data.ValidationErrors = model.ValidationErrors{
				"BondProvider": map[string]string{"required": "Please select a bond provider"},
			}
			w.WriteHeader(http.StatusUnprocessableEntity)
			return h.execute(w, r, data)
		}

		bondProvider := bondProviders.GetById(bondProviderId)
		if bondProvider == nil {
			data.ValidationErrors = model.ValidationErrors{
				"BondProvider": map[string]string{"required": "Bond provider not recognised"},
			}
			w.WriteHeader(http.StatusUnprocessableEntity)
			return h.execute(w, r, data)
		}

		file, handler, err := r.FormFile("fileUpload")
		if err != nil {
			return err
		}

		unchecked(file.Close)

		fileData, err := io.ReadAll(file)
		if err != nil {
			return err
		}

		err = h.router.Client().Upload(ctx, shared.Upload{
			UploadType:   shared.UploadTypeBonds,
			Base64Data:   base64.StdEncoding.EncodeToString(fileData),
			Filename:     handler.Filename,
			BondProvider: *bondProvider,
		})
		if err != nil {
			return err
		}

		w.Header().Add("HX-Redirect", fmt.Sprintf("%s/uploads?success=upload", v.EnvironmentVars.Prefix))
	case shared.UploadTypeUnknown:
		data.ValidationErrors = model.ValidationErrors{
			"UploadType": map[string]string{"required": "Please select a report to upload"},
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		return h.execute(w, r, data)
	}
	return h.execute(w, r, data)
}
