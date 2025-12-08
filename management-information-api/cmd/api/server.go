package api

import (
	"bytes"
	"context"
	"github.com/ministryofjustice/opg-go-common/securityheaders"
	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/opg-sirius-supervision-management-information/shared"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"io"
	"log/slog"
	"net/http"
)

type Service interface {
	ProcessDirectUpload(ctx context.Context, uploadType shared.UploadType, fileName string, fileBytes bytes.Reader) error
}

type FileStorage interface {
	GetFile(ctx context.Context, bucketName string, filename string) (io.ReadCloser, error)
}

type Server struct {
	service Service
}

func NewServer(service Service) *Server {
	return &Server{
		service: service,
	}
}

func (s *Server) SetupRoutes(logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /uploads", s.ProcessDirectUpload)

	mux.Handle("/health-check", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	return otelhttp.NewHandler(telemetry.Middleware(logger)(securityheaders.Use(mux)), "supervision-finance-api")
}

// unchecked allows errors to be unchecked when deferring a function, e.g. closing a reader where a failure would only
// occur when the process is likely to already be unrecoverable
func unchecked(f func() error) {
	_ = f()
}
