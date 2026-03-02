package filestorage

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type s3Client interface {
	Options() s3.Options
}

type uploader interface {
	UploadObject(ctx context.Context, input *transfermanager.UploadObjectInput, opts ...func(*transfermanager.Options)) (*transfermanager.UploadObjectOutput, error)
}

type Client struct {
	s3       s3Client
	kmsKey   string
	uploader uploader
}

func NewClient(ctx context.Context, region string, iamRole string, endpoint string, kmsKey string) (*Client, error) {
	awsRegion := region

	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(awsRegion),
	)
	if err != nil {
		return nil, err
	}

	if iamRole != "" {
		client := sts.NewFromConfig(cfg)
		cfg.Credentials = stscreds.NewAssumeRoleProvider(client, iamRole)
	}

	s3Client := s3.NewFromConfig(cfg, func(u *s3.Options) {
		u.UsePathStyle = true
		u.Region = awsRegion

		if endpoint != "" {
			u.BaseEndpoint = &endpoint
		}
	})

	uploader := transfermanager.New(s3Client)

	return &Client{
		s3:       s3Client,
		kmsKey:   kmsKey,
		uploader: uploader,
	}, nil
}

func (c *Client) StreamFile(ctx context.Context, bucketName string, fileName string, stream io.ReadCloser) (*string, error) {
	output, err := c.uploader.UploadObject(ctx, &transfermanager.UploadObjectInput{
		Bucket:               aws.String(bucketName),
		Key:                  aws.String(fileName),
		Body:                 stream,
		ContentType:          aws.String("text/csv"),
		ServerSideEncryption: "aws:kms",
		SSEKMSKeyID:          aws.String(c.kmsKey),
	})

	if err != nil {
		return nil, err
	}

	return output.VersionID, err
}
