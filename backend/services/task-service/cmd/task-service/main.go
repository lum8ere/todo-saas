package main

import (
	"fmt"
	"os"

	platformconfig "github.com/lum1ere/todo-saas/backend/libs/platform-config"
	platformdb "github.com/lum1ere/todo-saas/backend/libs/platform-db"
	platformhttp "github.com/lum1ere/todo-saas/backend/libs/platform-http"
	platformlogger "github.com/lum1ere/todo-saas/backend/libs/platform-logger"

	"go.uber.org/zap"
)

type Config struct {
	AppEnv   string
	AppName  string
	HTTPPort string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

func loadConfig() Config {
	return Config{
		AppEnv:   platformconfig.GetEnv("APP_ENV", "local"),
		AppName:  platformconfig.GetEnv("APP_NAME", "task-service"),
		HTTPPort: platformconfig.GetEnv("HTTP_PORT", "8082"),

		DBHost:     platformconfig.GetEnv("DB_HOST", "localhost"),
		DBPort:     platformconfig.GetEnv("DB_PORT", "5432"),
		DBUser:     platformconfig.GetEnv("DB_USER", "task_user"),
		DBPassword: platformconfig.GetEnv("DB_PASSWORD", "task_password"),
		DBName:     platformconfig.GetEnv("DB_NAME", "task_db"),
		DBSSLMode:  platformconfig.GetEnv("DB_SSLMODE", "disable"),
	}
}

func main() {
	cfg := loadConfig()

	logger, err := platformlogger.New(cfg.AppName, cfg.AppEnv)
	if err != nil {
		fmt.Println("failed to init logger:", err)
		os.Exit(1)
	}
	defer logger.Sync()

	db, err := platformdb.NewPostgres(platformdb.Config{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
		SSLMode:  cfg.DBSSLMode,
	}, logger.Base)
	if err != nil {
		logger.Base.Fatal("db connection failed", zap.Error(err))
	}

	_ = db // пока не используем

	router := platformhttp.NewDefaultRouter(logger)
	// здесь потом добавишь свои руты:
	// router.Route("/api/tasks", func(r chi.Router) { ... })

	server := platformhttp.NewServer(
		platformhttp.ServerConfig{
			Addr: ":" + cfg.HTTPPort,
		},
		router,
		logger,
	)

	server.Run()
}
