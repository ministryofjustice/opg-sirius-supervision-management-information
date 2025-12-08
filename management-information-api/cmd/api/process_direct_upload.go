package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/opg-sirius-supervision-management-information/shared"
	"net/http"
)

func (s *Server) ProcessDirectUpload(w http.ResponseWriter, r *http.Request) {
	var upload shared.Upload
	defer unchecked(r.Body.Close)

	if err := json.NewDecoder(r.Body).Decode(&upload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fileBytes, err := base64.StdEncoding.DecodeString(upload.Base64Data)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	err = s.service.ProcessDirectUpload(r.Context(), upload.UploadType, upload.Filename, *bytes.NewReader(fileBytes))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
