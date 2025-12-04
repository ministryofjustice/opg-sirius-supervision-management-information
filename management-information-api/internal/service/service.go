package service

import (
	"context"
	"io"
	"log/slog"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/telemetry"
)

type FileStorage interface {
	StreamFile(ctx context.Context, bucketName string, fileName string, stream io.ReadCloser) (*string, error)
}

type Envs struct {
	Port string
	AwsRegion string
	IamRole string
	S3Endpoint string
	S3EncryptionKey string
	AsyncBucket string
}

type Service struct {
	fileStorage FileStorage
	env         *Envs
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewService(fileStorage FileStorage, env *Envs) *Service {
	return &Service{
		fileStorage: fileStorage,
		env:         env,
	}
}

func (s *Service) Logger(ctx context.Context) *slog.Logger {
	return telemetry.LoggerFromContext(ctx)
}
