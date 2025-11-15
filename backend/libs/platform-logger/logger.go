package logger

import (
	"context"

	"go.uber.org/zap"
)

type ctxKey string

const (
	contextKeyRequestID ctxKey = "request_id"
)

// Logger — обёртка вокруг zap с сервисным именем и sugar.
type Logger struct {
	Base    *zap.Logger
	Sugar   *zap.SugaredLogger
	service string
	env     string
}

// New создает логгер для сервиса.
// env: "local", "dev", "prod" и т.п. — влияет на формат логов.
func New(serviceName string, env string) (*Logger, error) {
	var cfg zap.Config

	switch env {
	case "local", "dev":
		cfg = zap.NewDevelopmentConfig()
	default:
		cfg = zap.NewProductionConfig()
		cfg.Encoding = "json"
	}

	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stderr"}

	l, err := cfg.Build(
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.Fields(
			zap.String("service", serviceName),
			zap.String("env", env),
		),
	)
	if err != nil {
		return nil, err
	}

	return &Logger{
		Base:    l,
		Sugar:   l.Sugar(),
		service: serviceName,
		env:     env,
	}, nil
}

// Sync безопасно синкает буферы.
func (l *Logger) Sync() {
	_ = l.Base.Sync()
}

// With добавляет поля + request_id из контекста.
func (l *Logger) With(ctx context.Context, fields ...zap.Field) *zap.Logger {
	log := l.Base
	if ctx != nil {
		if reqID := RequestIDFromContext(ctx); reqID != "" {
			log = log.With(zap.String("request_id", reqID))
		}
	}
	if len(fields) > 0 {
		log = log.With(fields...)
	}
	return log
}

// SugaredWith — sugar-логгер, но всё ещё с request_id.
func (l *Logger) SugaredWith(ctx context.Context, fields ...interface{}) *zap.SugaredLogger {
	log := l.With(ctx)
	if len(fields) > 0 {
		return log.Sugar().With(fields...)
	}
	return log.Sugar()
}
