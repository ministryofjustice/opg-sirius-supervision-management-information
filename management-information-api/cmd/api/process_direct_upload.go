package api

import (
	"fmt"
	"net/http"
)

func (s *Server) ProcessDirectUpload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("We got to the ProcessDirectUpload function!")
}
