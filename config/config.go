package config

import (
	"fmt"
	"log"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
	// "github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string     `mapstructure:"env"`
	StoragePath string     `mapstructure:"storage_path"`
	HTTPServer  HTTPServer `mapstructure:"http_server"`
	Database    Database   `mapstructure:"database"`
}

type HTTPServer struct {
	Address     string        `mapstructure:"address"`
	Timeout     time.Duration `mapstructure:"timeout"`
	IdleTimeout time.Duration `mapstructure:"idle_timeout"`
}

type Database struct {
	Port       int    `mapstructure:"port"`
	Host       string `mapstructure:"host"`
	Name       string `mapstructure:"name"`
	User       string `mapstructure:"postgres"`
	Password   string `mapstructure:"password"`
	ConnString string `mapstructure:"conn_string"`
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
	if err := v.Unmarshal(&c, func(dc *mapstructure.DecoderConfig) {
		dc.DecodeHook = mapstructure.StringToTimeDurationHookFunc()
	}); err != nil {
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
