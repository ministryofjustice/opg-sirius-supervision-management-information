package api

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/opg-sirius-supervision-management-information/shared"
)

func (s *Server) ProcessDirectUpload(w http.ResponseWriter, r *http.Request) error {
	var upload shared.Upload
	ctx := r.Context()

	defer unchecked(r.Body.Close)

	if err := json.NewDecoder(r.Body).Decode(&upload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	fileBytes, err := base64.StdEncoding.DecodeString(upload.Base64Data)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return err
	}

	if upload.Filename == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return fmt.Errorf("filename is required")
	}

	filePath := fmt.Sprintf("%s/%s", upload.UploadType.Directory(), upload.Filename)

	_, err = s.fileStorage.StreamFile(context.Background(), s.asyncBucket, filePath, io.NopCloser(bytes.NewReader(fileBytes)))

	if err != nil {
		s.Logger(ctx).Error("Unable to upload file: ", "err", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	return nil
}

