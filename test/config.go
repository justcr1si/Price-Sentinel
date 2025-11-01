package test

// package config

// import (
// 	"os"
// 	"path/filepath"
// 	"price_sentinel/config"
// 	"testing"
// 	"time"
// )

// func TestLoad_BuildsConnString(t *testing.T) {
// 	tmp := t.TempDir()
// 	yaml := `
// env: "local"
// http_server:
//   address: "localhost:8080"
//   timeout: 4s
//   idle_timeout: 60s
// database:
//   port: 5432
//   host: "localhost"
//   name: "price_sentinel"
//   user: "postgres"
//   password: "secret"
// `
// 	cfgPath := filepath.Join(tmp, "config.yaml")
// 	if err := os.WriteFile(cfgPath, []byte(yaml), 0o600); err != nil {
// 		t.Fatalf("write config: %v", err)
// 	}

// 	oldWD, _ := os.Getwd()
// 	defer func() { _ = os.Chdir(oldWD) }()
// 	if err := os.Chdir(tmp); err != nil {
// 		t.Fatalf("chdir: %v", err)
// 	}

// 	c, err := config.Load()
// 	if err != nil {
// 		t.Fatalf("Config wasn't loaded")
// 	}

// 	if c.Env != "local" {
// 		t.Fatalf("Env = %q, want local", c.Env)
// 	}
// 	if c.HTTPServer.Address != "localhost:8080" {
// 		t.Fatalf("HTTP addr = %q, want localhost:8080", c.HTTPServer.Address)
// 	}
// 	if c.HTTPServer.Timeout != 4*time.Second {
// 		t.Fatalf("Timeout = %v, want 4s", c.HTTPServer.Timeout)
// 	}
// 	if c.HTTPServer.IdleTimeout != 60*time.Second {
// 		t.Fatalf("IdleTimeout = %v, want 60s", c.HTTPServer.IdleTimeout)
// 	}

// 	if c.Database.Port != 5432 ||
// 		c.Database.Host != "localhost" ||
// 		c.Database.Name != "price_sentinel" ||
// 		c.Database.User != "postgres" ||
// 		c.Database.Password != "secret" {
// 		t.Fatalf("database fields parsed incorrectly: %#v", c.Database)
// 	}

// 	wantDSN := "postgresql://postgres:secret@localhost:5432/price_sentinel"
// 	if c.Database.ConnString != wantDSN {
// 		t.Fatalf("ConnString = %q, want %q", c.Database.ConnString, wantDSN)
// 	}
// }
