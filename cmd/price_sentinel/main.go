package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"price_sentinel/config"
	"price_sentinel/internal/handler"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Config wasn't loaded: %v", err)
	}
	fmt.Println(cfg)
	fmt.Println(cfg.HTTPServer)
	fmt.Println(cfg.Database)

	logger := setupLogger(cfg.Env)
	log.Logger = *logger

	log.Info().Str("env", cfg.Env).Msg("app starting")

	authRouter := chi.NewRouter()
	authRouter.Route("/auth", func(r chi.Router) {
		r.Post("/signup", handler.SignUp)
		r.Post("/login", handler.Login)
		r.Post("/refresh", handler.Refresh)
	})

	fmt.Println(authRouter)

	fmt.Println(cfg.HTTPServer.Address)
	fmt.Println(cfg.HTTPServer.Address[len(cfg.HTTPServer.Address)-5:])
	http.ListenAndServe(cfg.HTTPServer.Address[len(cfg.HTTPServer.Address)-5:], authRouter)
}

func setupLogger(env string) *zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339

	var writer io.Writer
	var level zerolog.Level

	switch env {
	case envLocal, envDev:
		writer = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05"}
		level = zerolog.DebugLevel
	case envProd:
		writer = os.Stdout
		level = zerolog.InfoLevel
	default:
		writer = os.Stdout
		level = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(level)

	l := zerolog.New(writer).
		Level(level).
		With().
		Timestamp().
		Str("env", env).
		Logger()

	return &l
}
