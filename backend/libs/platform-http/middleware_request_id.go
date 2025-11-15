package http

import (
	"net/http"

	"github.com/google/uuid"
	platformlogger "github.com/lum1ere/todo-saas/backend/libs/platform-logger"
)

// RequestIDMiddleware:
// - читает X-Request-ID из заголовка, если нет — генерирует;
// - записывает его в заголовок ответа;
// - кладёт в контекст (через platform-logger.WithRequestID).
func RequestIDMiddleware(next http.Handler) http.Handler {
	const headerName = "X-Request-ID"

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get(headerName)
		if reqID == "" {
			reqID = uuid.NewString()
		}

		w.Header().Set(headerName, reqID)

		ctx := platformlogger.WithRequestID(r.Context(), reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
