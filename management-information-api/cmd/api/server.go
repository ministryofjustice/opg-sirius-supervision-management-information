package api

import (
	"context"
	"github.com/ministryofjustice/opg-go-common/securityheaders"
	"github.com/ministryofjustice/opg-go-common/telemetry"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"io"
	"log/slog"
	"net/http"
)

type Service interface {
	ProcessDirectUpload(ctx context.Context, filename string, fileBytes io.Reader) error
}

type FileStorage interface {
	GetFile(ctx context.Context, bucketName string, filename string) (io.ReadCloser, error)
}

type Server struct {
	service     Service
	fileStorage FileStorage
	envs        *Envs
}

type Envs struct {
	Port string
}

func NewServer(envs Envs) *Server {
	return &Server{envs: &envs}
}

func (s *Server) SetupRoutes(logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /uploads", s.ProcessDirectUpload)

	mux.Handle("/health-check", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	return otelhttp.NewHandler(telemetry.Middleware(logger)(securityheaders.Use(mux)), "supervision-finance-api")
}
