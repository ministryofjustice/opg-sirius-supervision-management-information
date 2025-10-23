package main

import (
	"context"
	"github.com/ministryofjustice/opg-go-common/env"
	"github.com/ministryofjustice/opg-go-common/paginate"
	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/opg-sirius-supervision-management-information/internal/api"
	"github.com/opg-sirius-supervision-management-information/internal/server"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
	"unicode"
)

func main() {
	ctx := context.Background()
	logger := telemetry.NewLogger("opg-sirius-management-information")

	err := run(ctx, logger)
	if err != nil {
		logger.Error("fatal startup error", slog.Any("err", err.Error()))
		os.Exit(1)
	}
}

func run(ctx context.Context, logger *slog.Logger) error {
	exportTraces := env.Get("TRACING_ENABLED", "0") == "1"

	shutdown, err := telemetry.StartTracerProvider(ctx, logger, exportTraces)
	defer shutdown()
	if err != nil {
		return err
	}

	envVars := server.NewEnvironmentVars()
	client := api.NewApiClient(http.DefaultClient)
	templates := createTemplates(envVars)

	s := &http.Server{
		Addr:              ":" + envVars.Port,
		Handler:           server.New(logger, client, templates, envVars),
		ReadHeaderTimeout: 2 * time.Second,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil {
			logger.Error("listen and server error", slog.Any("err", err.Error()))
			os.Exit(1)
		}
	}()

	logger.Info("Running at :" + envVars.Port)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	sig := <-c
	logger.Info("signal received: ", "sig", sig)

	tc, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	return s.Shutdown(tc)
}

func createTemplates(envVars server.EnvironmentVars) map[string]*template.Template {
	templates := map[string]*template.Template{}
	templateFunctions := map[string]interface{}{
		"contains": func(xs []string, needle string) bool {
			for _, x := range xs {
				if x == needle {
					return true
				}
			}

			return false
		},
		"title": func(s string) string {
			r := []rune(s)
			letter := r[0]
			r[0] = unicode.ToUpper(letter)

			return string(r)
		},
		"prefix": func(s string) string {
			return envVars.Prefix + s
		},
		"sirius": func(s string) string {
			return envVars.SiriusPublicURL + s
		},
	}

	templateDirPath := filepath.Clean(envVars.WebDir + "/template")
	templateDir, _ := os.Open(templateDirPath)
	templateDirs, _ := templateDir.Readdir(0)
	_ = templateDir.Close()

	mainTemplates, _ := filepath.Glob(templateDirPath + "/*.gotmpl")

	for _, file := range mainTemplates {
		tmpl := template.New(filepath.Base(file)).Funcs(templateFunctions)
		for _, dir := range templateDirs {
			if dir.IsDir() {
				tmpl, _ = tmpl.ParseGlob(templateDirPath + "/" + dir.Name() + "/*.gotmpl")
			}
		}
		tmpl, _ = tmpl.Parse(paginate.Template)
		templates[tmpl.Name()] = template.Must(tmpl.ParseFiles(file))
	}

	return templates
}
