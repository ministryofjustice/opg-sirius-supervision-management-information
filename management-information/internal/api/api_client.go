package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/opg-sirius-supervision-management-information/management-information/internal/auth"
	"io"
	"log/slog"
	"net/http"
	"strconv"
)

type ClientError string

const ErrUnauthorized ClientError = "unauthorized"

type ValidationErrors map[string]map[string]string

type ValidationError struct {
	Message string
	Errors  ValidationErrors
}

type StatusError struct {
	Code   int    `json:"code"`
	URL    string `json:"url"`
	Method string `json:"method"`
}

func (e StatusError) Error() string {
	return fmt.Sprintf("%s %s returned %d", e.Method, e.URL, e.Code)
}

func (e StatusError) Title() string {
	return "unexpected response from Sirius"
}

func (e StatusError) Data() interface{} {
	return e
}

func newStatusError(resp *http.Response) StatusError {
	return StatusError{
		Code:   resp.StatusCode,
		URL:    resp.Request.URL.String(),
		Method: resp.Request.Method,
	}
}

func (ve ValidationError) Error() string {
	return ve.Message
}

func (e ClientError) Error() string {
	return string(e)
}

func NewApiClient(httpClient HTTPClient, jwt JWTClient, baseURL string, logger *slog.Logger, backendURL string) (*ApiClient, error) {
	return &ApiClient{
		http:       httpClient,
		jwt:        jwt,
		baseURL:    baseURL,
		logger:     logger,
		backendURL: backendURL,
	}, nil
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type FileStorageInterface interface {
	StreamFile(ctx context.Context, bucketName string, fileName string, stream io.ReadCloser) (*string, error)
}

type JWTClient interface {
	CreateJWT(ctx context.Context) string
}

type ApiClient struct {
	http       HTTPClient
	baseURL    string
	logger     *slog.Logger
	backendURL string
	jwt        JWTClient
}

func addXsrfFromContext(ctx context.Context, req *http.Request) {
	req.Header.Add("X-XSRF-TOKEN", ctx.(auth.Context).XSRFToken)
}

func addCookiesFromContext(ctx context.Context, req *http.Request) {
	for _, c := range ctx.(auth.Context).Cookies {
		req.AddCookie(c)
	}
}

func (c *ApiClient) newRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return nil, err
	}

	addCookiesFromContext(ctx, req)
	addXsrfFromContext(ctx, req)
	req.Header.Add("OPG-Bypass-Membrane", "1")

	return req, err
}

func (c *ApiClient) newBackendRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.backendURL+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+c.jwt.CreateJWT(ctx))

	return req, err
}

func (c *ApiClient) logErrorRequest(req *http.Request, err error) {
	c.logger.Info("method: " + req.Method + ", url: " + req.URL.Path)
	if err != nil {
		c.logger.Error(err.Error())
	}
}

func (c *ApiClient) logResponse(req *http.Request, resp *http.Response, err error) {
	response := "None"
	if resp != nil {
		response = strconv.Itoa(resp.StatusCode)
	}
	c.logger.Info("method: " + req.Method + ", url: " + req.URL.Path + ", response: " + response)
	if err != nil && !errors.Is(err, context.Canceled) {
		c.logger.Error(err.Error())
	}
}

// unchecked allows errors to be unchecked when deferring a function, e.g. closing a reader where a failure would only
// occur when the process is likely to already be unrecoverable
func unchecked(f func() error) {
	_ = f()
}

type ExpandedError interface {
	Title() string
	Data() interface{}
}
