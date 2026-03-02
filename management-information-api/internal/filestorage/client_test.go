package filestorage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	got, err := NewClient(context.Background(), "eu-west-1", "role", "some-endpoint", "key")

	assert.Nil(t, err)

	assert.IsType(t, new(Client), got)
	assert.Equal(t, "eu-west-1", got.s3.Options().Region)
	assert.Equal(t, "some-endpoint", *got.s3.Options().BaseEndpoint)
	assert.Equal(t, "key", got.kmsKey)
}

type mockUploader struct {
	output *transfermanager.UploadObjectOutput
	err    error
}

func (m *mockUploader) UploadObject(ctx context.Context, input *transfermanager.UploadObjectInput, opts ...func(*transfermanager.Options)) (*transfermanager.UploadObjectOutput, error) {
	return m.output, m.err
}

func TestStreamFile(t *testing.T) {
	versionId := "test"
	tests := []struct {
		name         string
		mockUploader *mockUploader
		want         *string
		wantErr      error
	}{
		{
			name: "success",
			mockUploader: &mockUploader{
				output: &transfermanager.UploadObjectOutput{VersionID: &versionId},
			},
			want:    &versionId,
			wantErr: nil,
		},
		{
			name: "fail",
			mockUploader: &mockUploader{
				err: errors.New("error"),
			},
			want:    nil,
			wantErr: fmt.Errorf("error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{uploader: tt.mockUploader}
			got, err := client.StreamFile(context.Background(), "bucket", "filename", io.NopCloser(strings.NewReader("test")))
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
