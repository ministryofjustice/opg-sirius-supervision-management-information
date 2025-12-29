package server

import (
	"errors"
	"github.com/opg-sirius-supervision-management-information/management-information/internal/auth"
	"golang.org/x/sync/errgroup"
	"net/http"
)

type PageData struct {
	Data           any
	SuccessMessage string
}

type route struct {
	client  ApiClient
	tmpl    Template
	partial string
}

func (r route) Client() ApiClient {
	return r.client
}

// execute is an abstraction of the Template execute functions in order to conditionally render either a full template or
// a block, in response to a header added by HTMX. If the header is not present, the function will also fetch all
// additional data needed by the page for a full page load.
func (r route) execute(w http.ResponseWriter, req *http.Request, data any) error {
	if IsHxRequest(req) {
		return r.tmpl.ExecuteTemplate(w, r.partial, data)
	} else {
		ctx, ok := req.Context().(auth.Context)
		if !ok {
			return errors.New("no auth context in request")
		}
		group, groupCtx := errgroup.WithContext(ctx)
		ctx = ctx.WithContext(groupCtx)
		data := PageData{
			Data:           data,
			SuccessMessage: r.getSuccess(req),
		}

		group.Go(func() error {
			user, err := r.client.GetCurrentUserDetails(ctx)
			if err != nil {
				return err
			}
			if !user.IsReportingUser() {
				return errors.New("not reporting user")
			}
			return nil
		})

		if err := group.Wait(); err != nil {
			return err
		}

		return r.tmpl.Execute(w, data)
	}
}

func IsHxRequest(req *http.Request) bool {
	return req.Header.Get("HX-Request") == "true"
}

func (r route) getSuccess(req *http.Request) string {
	switch req.URL.Query().Get("success") {
	case "upload":
		return "File successfully uploaded"
	}
	return ""
}
