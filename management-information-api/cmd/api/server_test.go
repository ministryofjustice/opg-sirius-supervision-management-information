package api

import (
	"context"
	"io"
)

type mockFileStorage struct {
	versionId  string
	bucketname string
	filename   string
	data       io.Reader
	err        error
}

func (m *mockFileStorage) StreamFile(ctx context.Context, bucketName string, fileName string, stream io.ReadCloser) (*string, error) {
	return &m.versionId, m.err
}
