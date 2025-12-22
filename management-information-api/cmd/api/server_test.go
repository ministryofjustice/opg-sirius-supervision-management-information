package api

import (
	"context"
	"io"
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
