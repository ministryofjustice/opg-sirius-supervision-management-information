package server

import (
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/opg-sirius-supervision-management-information/management-information/internal/model"
	"github.com/opg-sirius-supervision-management-information/shared"
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
	ctx := r.Context()
	// Limit request body size to 10MB to prevent memory exhaustion
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

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

		file, _, err := r.FormFile("fileUpload")
		if err != nil {
			data.ValidationErrors = model.ValidationErrors{
				"FileUpload": map[string]string{"required": "No file uploaded"},
			}
			w.WriteHeader(http.StatusUnprocessableEntity)
			return h.execute(w, r, data)
		}

		unchecked(file.Close)

		fileData, err := io.ReadAll(file)
		if err != nil {
			telemetry.LoggerFromContext(ctx).Error("unable to read file", "error", err)
			data.ValidationErrors = model.ValidationErrors{
				"FileUpload": map[string]string{"invalid": "Failed to read file"},
			}
			w.WriteHeader(http.StatusUnprocessableEntity)
			return h.execute(w, r, data)
		}

		csvReader := csv.NewReader(bytes.NewReader(fileData))
		rec, err := csvReader.ReadAll()
		if err != nil || len(rec) == 0 {
			telemetry.LoggerFromContext(ctx).Error("unable to read csv data", "error", err, "records", len(rec))
			data.ValidationErrors = model.ValidationErrors{
				"FileUpload": map[string]string{"invalid": "File does not contain valid CSV data"},
			}
			w.WriteHeader(http.StatusUnprocessableEntity)
			return h.execute(w, r, data)
		}

		upload := shared.Upload{
			UploadType:   shared.UploadTypeBonds,
			Base64Data:   base64.StdEncoding.EncodeToString(fileData),
			BondProvider: bondProvider,
		}

		upload.Filename, err = uploadedFileName(upload)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return err
		}

		err = h.router.Client().Upload(ctx, upload)
		if err != nil {
			return err
		}

		w.Header().Add("HX-Redirect", fmt.Sprintf("%s/uploads?success=upload", v.EnvironmentVars.Prefix))
	case shared.UploadTypeVisits:
		file, _, err := r.FormFile("fileUpload")
		if err != nil {
			data.ValidationErrors = model.ValidationErrors{
				"FileUpload": map[string]string{"required": "No file uploaded"},
			}
			w.WriteHeader(http.StatusUnprocessableEntity)
			return h.execute(w, r, data)
		}

		unchecked(file.Close)

		fileData, err := io.ReadAll(file)
		if err != nil {
			telemetry.LoggerFromContext(ctx).Error("unable to read file", "error", err)
			data.ValidationErrors = model.ValidationErrors{
				"FileUpload": map[string]string{"invalid": "Failed to read file"},
			}
			w.WriteHeader(http.StatusUnprocessableEntity)
			return h.execute(w, r, data)
		}

		csvReader := csv.NewReader(bytes.NewReader(fileData))
		rec, err := csvReader.ReadAll()
		if err != nil || len(rec) == 0 {
			telemetry.LoggerFromContext(ctx).Error("unable to read csv data", "error", err, "records", len(rec))
			data.ValidationErrors = model.ValidationErrors{
				"FileUpload": map[string]string{"invalid": "File does not contain valid CSV data"},
			}
			w.WriteHeader(http.StatusUnprocessableEntity)
			return h.execute(w, r, data)
		}

		upload := shared.Upload{
			UploadType: shared.UploadTypeVisits,
			Base64Data: base64.StdEncoding.EncodeToString(fileData),
		}

		upload.Filename, err = uploadedFileName(upload)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return err
		}

		err = h.router.Client().Upload(ctx, upload)
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

func uploadedFileName(upload shared.Upload) (string, error) {
	switch upload.UploadType {
	case shared.UploadTypeBonds:
		return fmt.Sprintf("%s_%s.csv", upload.BondProvider.Name, time.Now().Format("02_01_2006")), nil
	case shared.UploadTypeVisits:
		return fmt.Sprintf("visits_%s.csv", time.Now().Format("02_01_2006")), nil
	default:
		return "", fmt.Errorf("upload type is required")
	}
}
