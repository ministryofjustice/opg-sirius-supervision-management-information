package server

import (
	"github.com/opg-sirius-supervision-management-information/management-information/internal/model"
	"net/http"
)

type UploadsVars struct {
	UploadTypes   []model.UploadType
	BondProviders []model.BondProvider
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
		model.UploadTypes,
		bondProviders,
		v}
	data.selectTab("uploads")
	return h.execute(w, r, data)
}
