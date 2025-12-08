package api

import (
	"bytes"
	"context"
	"github.com/opg-sirius-supervision-management-information/shared"
)

type mockService struct {
	err              error
	lastCalledParams []interface{}
}

func (s *mockService) ProcessDirectUpload(ctx context.Context, uploadType shared.UploadType, fileName string, fileBytes bytes.Reader) error {
	s.lastCalledParams = []interface{}{ctx, uploadType, fileName, fileBytes}
	return s.err
}
