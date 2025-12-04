package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/opg-sirius-supervision-management-information/management-information-api/cmd/api"
	"github.com/opg-sirius-supervision-management-information/management-information-api/internal/filestorage"
	"github.com/opg-sirius-supervision-management-information/management-information-api/internal/service"
)

func main() {

	fmt.Println("Hey we're in the main function")

	ctx := context.Background()
	logger := telemetry.NewLogger("opg-sirius-supervision-management-information-api")

	err := run(ctx, logger)
	if err != nil {
		logger.Error("fatal startup error", slog.Any("err", err.Error()))
		os.Exit(1)
	}
}

func run(ctx context.Context, logger *slog.Logger) error {
	fmt.Println("Hey we're in the run function")
	exportTraces := os.Getenv("TRACING_ENABLED") == "1"
	shutdown, err := telemetry.StartTracerProvider(ctx, logger, exportTraces)
	defer shutdown()
	if err != nil {
		return err
	}

	envs := &service.Envs{
		Port: os.Getenv("PORT"),
		AwsRegion: os.Getenv("AWS_REGION"),
		IamRole: os.Getenv("IAM_ROLE"),
		S3Endpoint: os.Getenv("S3_ENDPOINT"),
		S3EncryptionKey: os.Getenv("S3_ENCRYPTION_KEY"),
	}

    //creates new instance of filestorage client and uses stuff from env vars
	fileStorageClient, err := filestorage.NewClient(
		ctx,
		envs.AwsRegion,
		envs.IamRole,
		envs.S3Endpoint,
		envs.S3EncryptionKey,
	)

	if err != nil {
		return err
	}

	// creates a new service 
	Service := service.NewService(fileStorageClient, envs)

	//validator, err := validation.New()
	//if err != nil {
	//	return err
	//}

	server := api.NewServer(Service)

	s := &http.Server{
		Addr:              ":" + envs.Port,
		Handler:           server.SetupRoutes(logger),
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil {
			logger.Error("listen and server error", slog.Any("err", err.Error()))
			os.Exit(1)
		}
	}()
	logger.Info("Running at :" + envs.Port)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	sig := <-c
	logger.Info("signal received: ", "sig", sig)

	tc, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	return s.Shutdown(tc)
}
