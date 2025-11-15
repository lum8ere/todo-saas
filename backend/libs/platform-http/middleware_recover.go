package http

import (
	"net/http"
	"runtime/debug"

	platformlogger "github.com/lum1ere/todo-saas/backend/libs/platform-logger"
	"go.uber.org/zap"
)

// RecoverMiddleware — ловит panic, логирует и отдаёт 500.
func RecoverMiddleware(l *platformlogger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					l.With(r.Context(),
						zap.Any("panic", rec),
						zap.ByteString("stack", debug.Stack()),
					).Error("panic recovered")

					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
