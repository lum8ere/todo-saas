package http

import (
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	platformlogger "github.com/lum1ere/todo-saas/backend/libs/platform-logger"
)

// NewDefaultRouter — chi.Router с базовыми middleware:
// - request ID,
// - real IP,
// - логирование запросов,
// - recover.
func NewDefaultRouter(logger *platformlogger.Logger) chi.Router {
	r := chi.NewRouter()

	r.Use(chimw.RealIP)
	r.Use(chimw.RequestID)
	r.Use(RequestIDMiddleware)
	r.Use(LoggingMiddleware(logger))
	r.Use(RecoverMiddleware(logger))

	return r
}
