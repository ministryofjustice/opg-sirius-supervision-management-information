package api

import (
	"context"
	"io"
)

type mockFileStorage struct {
	versionId string
	err       error
}

func (m *mockFileStorage) StreamFile(ctx context.Context, bucketName string, fileName string, stream io.ReadCloser) (*string, error) {
	return &m.versionId, m.err
}
