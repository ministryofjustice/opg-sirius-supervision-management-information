package server

import (
	"net/http"
)

type UploadsVars struct {
	AppVars
}

type GetUploadsHandler struct {
	router
}

func (h *GetUploadsHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	data := UploadsVars{v}
	data.selectTab("uploads")
	return h.execute(w, r, data)
}
