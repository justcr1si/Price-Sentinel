package config

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
	// "github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string     `yaml:"env"`
	StoragePath string     `yaml:"storage_path"`
	HTTPServer  HTTPServer `yaml:"http_server"`
	Database    Database   `yaml:"database"`
}

type HTTPServer struct {
	Address     string        `yaml:"address"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

type Database struct {
	Port       int    `yaml:"port"`
	Host       string `yaml:"host"`
	Name       string `yaml:"name"`
	User       string `yaml:"postgres"`
	Password   string `yaml:"password"`
	ConnString string `yaml:"conn_string"`
}

func Load() (*Config, error) {
	return LoadFrom("config.yaml")
}

func LoadFrom(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var c Config
	if err := v.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	c.Database.ConnString = fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
	)

	return &c, nil
}

func MustLoad() *Config {
	c, err := Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	return c
}
