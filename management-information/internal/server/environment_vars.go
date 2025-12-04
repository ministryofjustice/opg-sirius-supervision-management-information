package server

import (
	"os"
)

type EnvironmentVars struct {
	Port            string
	WebDir          string
	SiriusURL       string
	SiriusPublicURL string
	Prefix          string
	AwsRegion       string
	IamRole         string
	S3Endpoint      string
	S3EncryptionKey string
	AsyncBucket     string
}

func NewEnvironmentVars() EnvironmentVars {
	return EnvironmentVars{
		Port:            getEnv("PORT", "1234"),
		WebDir:          getEnv("WEB_DIR", "web"),
		SiriusURL:       getEnv("SIRIUS_URL", "http://localhost:8080"),
		SiriusPublicURL: getEnv("SIRIUS_PUBLIC_URL", ""),
		Prefix:          getEnv("PREFIX", ""),
		AwsRegion:       getEnv("AWS_REGION", ""),
		IamRole:         getEnv("AWS_IAM_ROLE", ""),
		S3Endpoint:      getEnv("AWS_S3_ENDPOINT", ""),
		S3EncryptionKey: getEnv("S3_ENCRYPTION_KEY", ""),
		AsyncBucket:     getEnv("ASYNC_S3_BUCKET", ""),
	}
}

func getEnv(key, def string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return def
}
