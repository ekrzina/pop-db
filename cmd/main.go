package main

import (
	"bytes"
	"context"
	"embed"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/haoli/pop-db/api/server"
	"github.com/haoli/pop-db/internal/dbman"
	"github.com/haoli/pop-db/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	v5 "github.com/swaggest/swgui/v5"
)

//go:embed configs/config.yaml
var configFS embed.FS

func initializeLogger() zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}
	logger := zerolog.New(consoleWriter).With().Timestamp().Logger()
	log.Logger = logger

	return logger
}

func loadConfig() (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigType("yaml")
	data, err := configFS.ReadFile("configs/config.yaml")
	if err != nil {
		return nil, err
	}
	if err := v.ReadConfig(bytes.NewBuffer(data)); err != nil {
		return nil, err
	}
	return v, nil
}

func main() {
	logger := initializeLogger()
	logger.Info().Msg("Starting PopDB...")

	// Load config
	v, err := loadConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to load config")
	}

	// Context for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	manager, err := dbman.NewDbManager(v, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize DB")
	}
	defer func() {
		if err := manager.DB.Close(); err != nil {
			logger.Error().Err(err).Msg("Failed to close DB")
		}
	}()

	repo := repository.NewPersonRepository(manager, &logger)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	apiProvider := server.NewApiProvider(&logger, *repo)

	strictHandler := server.NewStrictHandler(
		apiProvider,
		nil,
	)

	server.RegisterHandlers(router, strictHandler)

	router.GET("/openapi.json", func(c *gin.Context) {
		spec, err := server.GetSwagger()
		if err != nil {
			logger.Error().Err(err).Msg("Failed to load swagger spec")
			c.Status(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, spec)
	})

	swaggerHandler := v5.NewHandler(
		"PopDB API",
		"/openapi.json",
		"/swagger",
	)

	router.GET("/swagger/*any", gin.WrapH(swaggerHandler))
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		logger.Info().Msg("Server running at http://localhost:8080")
		logger.Info().Msg("Swagger UI at http://localhost:8080/swagger/")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("Server failed")
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()
	logger.Info().Msg("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("Shutdown failed")
	}

	logger.Info().Msg("Server stopped cleanly")
}
