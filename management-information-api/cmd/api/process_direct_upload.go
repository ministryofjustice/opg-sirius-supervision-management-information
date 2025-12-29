package api

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/opg-sirius-supervision-management-information/shared"
	"io"
	"net/http"
	"time"
)

func (s *Server) ProcessDirectUpload(w http.ResponseWriter, r *http.Request) error {
	var upload shared.Upload

	defer unchecked(r.Body.Close)

	if err := json.NewDecoder(r.Body).Decode(&upload); err != nil {
		return err
	}

	fileBytes, err := base64.StdEncoding.DecodeString(upload.Base64Data)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s_%s.csv", upload.BondProvider.Name, time.Now().Format("02_01_2006"))
	filePath := fmt.Sprintf("%s/%s", upload.UploadType.Directory(), fileName)

	_, err = s.fileStorage.StreamFile(context.Background(), s.asyncBucket, filePath, io.NopCloser(bytes.NewReader(fileBytes)))

	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	return nil
}
