package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"test-messenger/internal/config"
	apphttp "test-messenger/internal/http"
	"test-messenger/internal/http/handler"
	"test-messenger/internal/repository"
	"test-messenger/internal/service"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Database config (no direct DSN usage)
	dbCfg, err := config.LoadDatabaseConfig()
	if err != nil {
		log.Fatal(err)
	}

	connStr := "postgres://" +
		dbCfg.User + ":" + dbCfg.Password + "@" +
		dbCfg.Host + ":" + dbCfg.Port + "/" +
		dbCfg.Name + "?sslmode=" + dbCfg.SSLMode

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Repositories
	chatRepo := repository.NewChatRepository(db)
	messageRepo := repository.NewMessageRepository(db)

	// Services
	chatService := service.NewChatService(*chatRepo, *messageRepo)
	messageService := service.NewMessageService(*chatRepo, *messageRepo)

	// Handlers
	chatHandler := handler.NewChatHandler(chatService)
	messageHandler := handler.NewMessageHandler(messageService)

	// Router
	router := apphttp.NewRouter(chatHandler, messageHandler)

	// HTTP server
	server := &http.Server{
		Addr:         ":8080",
		Handler:      router.Handler(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Graceful shutdown
	go func() {
		slog.Info("http server started", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("http server error", "error", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("server shutdown failed", "error", err)
	}

	slog.Info("server stopped")
}
