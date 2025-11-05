package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"price_sentinel/config"
	"price_sentinel/internal/db"
	"time"

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

	logger := setupLogger(cfg.Env)
	log.Logger = *logger

	// authRouter := chi.NewRouter().
	// 	Route("/auth", func(r chi.Router) {
	// 		r.Post("/signup", auth.SignUp)
	// 		r.Post("/login", auth.Login)
	// 		r.Post("/refresh", auth.Refresh)
	// 	})
	db, err := db.DBClient(cfg.Database.ConnString)

	if err != nil {
		log.Err(err)
	}

	server := config.CreateServer(db)
	server.MountMiddleware()
	server.MountAuthHandlers()

	fmt.Println(server)

	addr := cfg.Server.Address[len(cfg.Server.Address)-5:]

	log.Info().Str("env", cfg.Env).Str("addr", addr).Msg("app starting")
	http.ListenAndServe(addr, server.Router)
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
