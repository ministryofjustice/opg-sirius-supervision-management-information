package api

import (
	"bytes"
	"fmt"
	"net/http"
)

func (s *Server) ProcessDirectUpload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In server, calling service.ProcessDirectUpload")
	filename := "test.csv" // Get from r
	fileBytes := []byte{}
	fileReader := bytes.NewReader(fileBytes)
	err := s.service.ProcessDirectUpload(r.Context(), filename, fileReader)
	if err != nil {
		return
	}
}
