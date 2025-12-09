package server

import (
	"encoding/base64"
	"fmt"
	"github.com/opg-sirius-supervision-management-information/management-information/internal/model"
	"github.com/opg-sirius-supervision-management-information/shared"
	"io"
	"net/http"
)

type UploadFileVars struct {
	UploadTypes      []shared.UploadType
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
		shared.UploadTypes,
		bondProviders,
		nil,
		v}
	data.selectTab("uploads")

	uploadType := shared.ParseUploadType(r.PostFormValue("uploadType"))

	if !uploadType.Valid() {
		data.ValidationErrors = model.ValidationErrors{
			"UploadType": map[string]string{"required": "Please select a report to upload"},
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		return h.execute(w, r, data)
	}

	switch uploadType {
	case shared.UploadTypeBonds:
		bondProvider := r.PostFormValue("bondProvider")
		if bondProvider == "" {
			data.ValidationErrors = model.ValidationErrors{
				"BondProvider": map[string]string{"required": "Please select a bond provider"},
			}
			w.WriteHeader(http.StatusUnprocessableEntity)
			return h.execute(w, r, data)
		}

		file, handler, err := r.FormFile("fileUpload")

		if err != nil {

			return err
		}

		defer func() {
			if err := file.Close(); err != nil {
				fmt.Println("Error closing file:", err)
			}
		}()

		fileData, err := io.ReadAll(file)

		if err != nil {
			return err
		}

		data := shared.Upload{
			UploadType: shared.UploadTypeBonds,
			Base64Data: base64.StdEncoding.EncodeToString(fileData),
			Filename:   handler.Filename,
		}

		err = h.router.Client().Upload(ctx, data)
		if err != nil {
			return err
		}

		w.Header().Add("HX-Redirect", fmt.Sprintf("%s/uploads?success=upload", v.EnvironmentVars.Prefix))
	}
	return h.execute(w, r, data)
}
