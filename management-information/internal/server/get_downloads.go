package server

import (
	"net/http"
)

type DownloadsVars struct {
	AppVars
}

type GetDownloadsHandler struct {
	router
}

func (h *GetDownloadsHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	data := DownloadsVars{v}
	data.selectTab("downloads")
	return h.execute(w, r, data)
}
