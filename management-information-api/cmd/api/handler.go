package api

import (
	"github.com/ministryofjustice/opg-go-common/telemetry"
	"log/slog"
	"net/http"
	"runtime/debug"
)

type handlerFunc func(w http.ResponseWriter, r *http.Request) error

func (f handlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rcv := recover(); rcv != nil {
			logger := telemetry.LoggerFromContext(r.Context())
			logger.Error("panic recovered", slog.Any("error", rcv), slog.String("stack", string(debug.Stack())))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}()

	if err := f(w, r); err != nil {
		logger := telemetry.LoggerFromContext(r.Context())
		logger.Error("error", slog.Any("error", err), slog.String("stack", string(debug.Stack())))
	}
}
