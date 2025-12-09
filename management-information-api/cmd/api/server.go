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
	StreamFile(ctx context.Context, bucketName string, fileName string, stream io.ReadCloser) (*string, error)
}

type Server struct {
	fileStorage FileStorage
	asyncBucket string
}

type Envs struct {
	Port            string
	AwsRegion       string
	IamRole         string
	S3Endpoint      string
	S3EncryptionKey string
	AsyncBucket     string
}

func NewServer(fileStorage FileStorage, asyncBucket string) *Server {
	return &Server{
		fileStorage: fileStorage,
		asyncBucket: asyncBucket,
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
