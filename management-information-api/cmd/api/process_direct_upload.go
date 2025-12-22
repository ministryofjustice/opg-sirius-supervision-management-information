package api

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/opg-sirius-supervision-management-information/shared"
	"io"
	"net/http"
	"time"
)

func (s *Server) ProcessDirectUpload(w http.ResponseWriter, r *http.Request) {
	var upload shared.Upload
	logger := telemetry.LoggerFromContext(r.Context())
	logger.Info("processing upload...")

	defer unchecked(r.Body.Close)

	if err := json.NewDecoder(r.Body).Decode(&upload); err != nil {
		logger.Error("backend error", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fileBytes, err := base64.StdEncoding.DecodeString(upload.Base64Data)
	if err != nil {
		logger.Error("backend error", "error", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	fileName := fmt.Sprintf("%s_%s.csv", upload.BondProvider.Name, time.Now().Format("02_01_2006"))
	filePath := fmt.Sprintf("%s/%s", upload.UploadType.Directory(), fileName)

	_, err = s.fileStorage.StreamFile(context.Background(), s.asyncBucket, filePath, io.NopCloser(bytes.NewReader(fileBytes)))

	if err != nil {
		logger.Error("backend error", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Info("file uploaded successfully")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
