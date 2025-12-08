package service

import (
	"bytes"
	"context"
	"fmt"
	"github.com/opg-sirius-supervision-management-information/shared"
	"io"
)

func (s *Service) ProcessDirectUpload(ctx context.Context, uploadType shared.UploadType, fileName string, fileBytes bytes.Reader) error {
	var directory string

	switch uploadType {
	case shared.UploadTypeBonds:
		directory = "bonds-with-orders"
	}

	filePath := fmt.Sprintf("%s/%s", directory, fileName)

	_, err := s.fileStorage.StreamFile(ctx, s.envs.AsyncBucket, filePath, io.NopCloser(&fileBytes))
	if err != nil {
		return err
	}
	return nil
}
