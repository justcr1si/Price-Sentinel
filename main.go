package main

import (
	"fmt"
	"log"
	"price_sentinel/config"

	"github.com/go-chi/chi/v5"
	// _ "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Config wasn't loaded: %v", err)
	}
	fmt.Println(cfg)

	router := chi.NewRouter()
	fmt.Println(router)
}
