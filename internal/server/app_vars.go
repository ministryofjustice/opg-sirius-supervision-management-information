package server

import (
	"github.com/opg-sirius-supervision-management-information/internal/api"
	"github.com/opg-sirius-supervision-management-information/internal/model"
	"net/http"
	"net/url"
)

type AppVars struct {
	Path            string
	XSRFToken       string
	Tabs            []Tab
	EnvironmentVars EnvironmentVars
	Error           string
	User            model.User
}

type Tab struct {
	Title    string
	Id       string
	Selected bool
}

func (t Tab) Path() string {
	return "/" + t.Id
}

func NewAppVars(r *http.Request, envVars EnvironmentVars) AppVars {
	tabs := []Tab{
		{
			Id:    "downloads",
			Title: "Downloads",
		},
		{
			Id:    "uploads",
			Title: "Uploads",
		},
	}

	var token string
	if r.Method == http.MethodGet {
		if cookie, err := r.Cookie("XSRF-TOKEN"); err == nil {
			token, _ = url.QueryUnescape(cookie.Value)
		}
	} else {
		token = r.FormValue("xsrfToken")
	}

	vars := AppVars{
		Path:            r.URL.Path,
		XSRFToken:       token,
		EnvironmentVars: envVars,
		Tabs:            tabs,
	}

	return vars
}

func (a *AppVars) selectTab(s string) {
	for i, tab := range a.Tabs {
		if tab.Id == s {
			a.Tabs[i] = Tab{
				Title:    tab.Title,
				Id:       tab.Id,
				Selected: true,
			}
		}
	}
}

func getContext(r *http.Request) api.Context {
	token := ""

	if r.Method == http.MethodGet {
		if cookie, err := r.Cookie("XSRF-TOKEN"); err == nil {
			token, _ = url.QueryUnescape(cookie.Value)
		}
	} else {
		token = r.FormValue("xsrfToken")
	}

	return api.Context{
		Context:   r.Context(),
		Cookies:   r.Cookies(),
		XSRFToken: token,
	}
}
