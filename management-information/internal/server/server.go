package server

import (
	"context"
	"github.com/ministryofjustice/opg-go-common/securityheaders"
	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/opg-sirius-supervision-management-information/management-information/internal/auth"
	"github.com/opg-sirius-supervision-management-information/shared"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"html/template"
	"io"
	"log/slog"
	"net/http"
)

type ApiClient interface {
	GetCurrentUserDetails(context.Context) (shared.User, error)
	GetBondProviders(context.Context) (shared.BondProviders, error)
	Upload(context.Context, shared.Upload) error
}

type router interface {
	Client() ApiClient
	execute(http.ResponseWriter, *http.Request, any) error
}

type Template interface {
	Execute(wr io.Writer, data any) error
	ExecuteTemplate(wr io.Writer, name string, data any) error
}

func New(logger *slog.Logger, client ApiClient, templates map[string]*template.Template, envVars EnvironmentVars) http.Handler {
	mux := http.NewServeMux()

	authenticator := auth.Auth{
		Client: client,
		EnvVars: auth.EnvVars{
			SiriusPublicURL: envVars.SiriusPublicURL,
			Prefix:          envVars.Prefix,
		},
	}

	handleMux := func(pattern string, h Handler) {
		errors := wrapHandler(templates["error.gotmpl"], "main", envVars)
		mux.Handle(pattern, authenticator.Authenticate(auth.XsrfCheck(errors(h))))
	}

	handleMux("GET /downloads", &GetDownloadsHandler{&route{client: client, tmpl: templates["downloads.gotmpl"], partial: "downloads"}})
	handleMux("GET /uploads", &GetUploadsHandler{&route{client: client, tmpl: templates["uploads.gotmpl"], partial: "uploads"}})

	handleMux("POST /uploads", &UploadFileHandler{&route{client: client, tmpl: templates["uploads.gotmpl"], partial: "error-summary"}})

	mux.Handle("/health-check", healthCheck())
	mux.Handle("/", http.RedirectHandler(envVars.Prefix+"/downloads", http.StatusFound))

	static := http.FileServer(http.Dir(envVars.WebDir + "/static"))
	mux.Handle("/assets/", static)
	mux.Handle("/javascript/", static)
	mux.Handle("/stylesheets/", static)

	return otelhttp.NewHandler(http.StripPrefix(envVars.Prefix, telemetry.Middleware(logger)(securityheaders.Use(mux))), "supervision-management-information")
}

// unchecked allows errors to be unchecked when deferring a function, e.g. closing a reader where a failure would only
// occur when the process is likely to already be unrecoverable
func unchecked(f func() error) {
	_ = f()
}
