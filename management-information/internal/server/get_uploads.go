package server

import (
	"github.com/opg-sirius-supervision-management-information/shared"
	"net/http"
)

type UploadsVars struct {
	UploadTypes   []shared.UploadType
	BondProviders []shared.BondProvider
	AppVars
}

type GetUploadsHandler struct {
	router
}

func (h *GetUploadsHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	ctx := getContext(r)
	bondProviders, err := h.router.Client().GetBondProviders(ctx)
	if err != nil {
		return err
	}
	data := UploadsVars{
		shared.UploadTypes,
		bondProviders,
		v}
	data.selectTab("uploads")
	return h.execute(w, r, data)
}
