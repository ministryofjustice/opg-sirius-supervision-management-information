package api

import (
	"bytes"
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ministryofjustice/opg-go-common/securityheaders"
	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/opg-sirius-supervision-management-information/management-information-api/internal/auth"
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

type JWTClient interface {
	Verify(requestToken string) (*jwt.Token, error)
}

type Server struct {
	http        HTTPClient
	fileStorage FileStorage
	asyncBucket string
	JWT         JWTClient
	baseURL     string
}

type Envs struct {
	Port            string
	AwsRegion       string
	IamRole         string
	S3Endpoint      string
	S3EncryptionKey string
	AsyncBucket     string
	JWTSecret       string
	SiriusURL       string
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewServer(httpClient HTTPClient, fileStorage FileStorage, asyncBucket string, jwtClient JWTClient, baseURL string) *Server {
	return &Server{
		http:        httpClient,
		fileStorage: fileStorage,
		asyncBucket: asyncBucket,
		JWT:         jwtClient,
		baseURL:     baseURL,
	}
}

func (s *Server) Logger(ctx context.Context) *slog.Logger {
	return telemetry.LoggerFromContext(ctx)
}

func (s *Server) requestLogger(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context().(auth.Context)
		s.Logger(ctx).Info(
			"API Request",
			"method", r.Method,
			"uri", r.URL.RequestURI(),
			"user-id", ctx.User.ID,
		)
		h.ServeHTTP(w, r)
	}
}

func (s *Server) SetupRoutes(logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	// authFunc is a replacement for mux.HandleFunc
	// which enriches the handler's HTTP instrumentation with the pattern as the http.route.
	authFunc := func(pattern string, role string, h handlerFunc) {
		// Configure the "http.route" for the HTTP instrumentation.
		handler := otelhttp.WithRouteTag(pattern, h)
		mux.Handle(pattern, s.authenticateAPI(s.requestLogger(s.authorise(role)(handler))))
	}

	authFunc("POST /uploads", shared.RoleReportingUser, s.ProcessDirectUpload)

	mux.Handle("/health-check", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	return otelhttp.NewHandler(telemetry.Middleware(logger)(securityheaders.Use(mux)), "supervision-finance-api")
}

func (s *Server) newRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, s.baseURL+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("OPG-Bypass-Membrane", "1")

	return req, err
}

// unchecked allows errors to be unchecked when deferring a function, e.g. closing a reader where a failure would only
// occur when the process is likely to already be unrecoverable
func unchecked(f func() error) {
	_ = f()
}
