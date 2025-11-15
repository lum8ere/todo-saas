package http

import (
	"net/http"
	"time"

	"go.uber.org/zap"

	platformlogger "github.com/lum1ere/todo-saas/backend/libs/platform-logger"
)

// обёртка, чтобы поймать статус и кол-во байт
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode  int
	wroteHeader bool
	bytes       int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	if !lrw.wroteHeader {
		lrw.statusCode = code
		lrw.wroteHeader = true
		lrw.ResponseWriter.WriteHeader(code)
	}
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	if !lrw.wroteHeader {
		lrw.WriteHeader(http.StatusOK)
	}
	n, err := lrw.ResponseWriter.Write(b)
	lrw.bytes += n
	return n, err
}

// LoggingMiddleware — логирует каждый HTTP-запрос:
// метод, путь, статус, длительность, размер ответа, user-agent.
func LoggingMiddleware(l *platformlogger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			lrw := newLoggingResponseWriter(w)

			next.ServeHTTP(lrw, r)

			duration := time.Since(start)

			log := l.With(r.Context(),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", lrw.statusCode),
				zap.Int("bytes", lrw.bytes),
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("user_agent", r.UserAgent()),
				zap.Duration("duration", duration),
			)

			// Уровень в зависимости от статуса
			switch {
			case lrw.statusCode >= 500:
				log.Error("http request")
			case lrw.statusCode >= 400:
				log.Warn("http request")
			default:
				log.Info("http request")
			}
		})
	}
}
