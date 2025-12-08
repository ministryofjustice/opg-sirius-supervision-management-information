package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/opg-sirius-supervision-management-information/shared"
	"net/http"
)

func (s *Server) ProcessDirectUpload(w http.ResponseWriter, r *http.Request) {
	var upload shared.Upload
	defer unchecked(r.Body.Close)

	if err := json.NewDecoder(r.Body).Decode(&upload); err != nil {
		// Throw an error
		return
	}

	fileBytes, err := base64.StdEncoding.DecodeString(upload.Base64Data)
	if err != nil {
		// throw an error
		return
	}

	err = s.service.ProcessDirectUpload(r.Context(), upload.UploadType, upload.Filename, *bytes.NewReader(fileBytes))
	if err != nil {
		fmt.Println(err)
	}
}
