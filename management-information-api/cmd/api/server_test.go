package api

import (
	"context"
	"github.com/ministryofjustice/opg-go-common/telemetry"
	"io"
	"log/slog"
	"net/http"
)

type mockFileStorage struct {
	fileName  string
	versionId string
	err       error
}

func (m *mockFileStorage) StreamFile(ctx context.Context, bucketName string, fileName string, stream io.ReadCloser) (*string, error) {
	m.fileName = fileName
	return &m.versionId, m.err
}

func SetUpTest() (*slog.Logger, MockClient) {
	logger := telemetry.NewLogger("opg-sirius-management-information")
	mockClient := MockClient{}
	return logger, mockClient
}

type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

var (
	// GetDoFunc fetches the mock client's `Do` func. Implement this within a test to modify the client's behaviour.
	GetDoFunc func(req *http.Request) (*http.Response, error)
)

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return GetDoFunc(req)
}
