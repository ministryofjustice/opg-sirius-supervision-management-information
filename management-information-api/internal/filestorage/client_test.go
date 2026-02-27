package filestorage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
)

type mockS3Client struct {
	putObjectInput  *s3.PutObjectInput
	putObjectOutput *s3.PutObjectOutput
	putObjectError  error
}

func (m *mockS3Client) PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	m.putObjectInput = params
	return m.putObjectOutput, m.putObjectError
}

func (m *mockS3Client) Options() s3.Options {
	return s3.Options{}
}

func TestNewClient(t *testing.T) {
	got, err := NewClient(context.Background(), "eu-west-1", "role", "some-endpoint", "key")

	assert.Nil(t, err)

	assert.IsType(t, new(Client), got)
	assert.Equal(t, "eu-west-1", got.s3.Options().Region)
	assert.Equal(t, "some-endpoint", *got.s3.Options().BaseEndpoint)
	assert.Equal(t, "key", got.kmsKey)
}

func TestStreamFile(t *testing.T) {
	versionId := "test"
	tests := []struct {
		name    string
		mockS3  *mockS3Client
		want    *string
		wantErr error
	}{
		{
			name: "success",
			mockS3: &mockS3Client{
				putObjectOutput: &s3.PutObjectOutput{VersionId: &versionId},
			},
			want:    &versionId,
			wantErr: nil,
		},
		{
			name: "fail",
			mockS3: &mockS3Client{
				putObjectError: errors.New("error"),
			},
			want:    nil,
			wantErr: fmt.Errorf("error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{s3: tt.mockS3}
			got, err := client.StreamFile(context.Background(), "bucket", "filename", io.NopCloser(strings.NewReader("test")))
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
