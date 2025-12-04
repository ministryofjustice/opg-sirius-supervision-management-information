package service

import (
	"context"
	"fmt"
	"io"
)

func (s *Service) ProcessDirectUpload(ctx context.Context, filename string, fileBytes io.Reader) error {
	directory := "bonds-without-orders"

	filePath := fmt.Sprintf("%s/%s", directory, filename)

	_, err := s.fileStorage.StreamFile(ctx, s.envs.AsyncBucket, filePath, io.NopCloser(fileBytes))
	if err != nil {
		return err
	}
	return nil
}
