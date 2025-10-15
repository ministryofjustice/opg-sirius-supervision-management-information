package server

import (
	"net/http"
)

type GetUploadsVars struct {
	AppVars
}

type GetUploadsHandler struct {
	router
}

func (h *GetUploadsHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	data := GetUploadsVars{v}
	data.selectTab("uploads")
	return h.execute(w, r, data)
}
