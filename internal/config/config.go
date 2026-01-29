package config

import (
	"fmt"
	"os"
	"time"

	"github.com/go-core-fx/config"
)

type http struct {
	Address     string   `koanf:"address"`
	ProxyHeader string   `koanf:"proxy_header"`
	Proxies     []string `koanf:"proxies"`
}

type database struct {
	URL string `koanf:"url"`

	ConnMaxIdleTime time.Duration `koanf:"conn_max_idle_time"`
	ConnMaxLifetime time.Duration `koanf:"conn_max_lifetime"`
	MaxOpenConns    int           `koanf:"max_open_conns"`
	MaxIdleConns    int           `koanf:"max_idle_conns"`
}

type redis struct {
	URL string `koanf:"url"`
}

type Config struct {
	HTTP     http     `koanf:"http"`
	Database database `koanf:"database"`
	Redis    redis    `koanf:"redis"`
}

func Default() Config {
	//nolint:mnd // default values
	return Config{
		HTTP: http{
			Address:     "127.0.0.1:3000",
			ProxyHeader: "X-Forwarded-For",
			Proxies:     []string{},
		},
		Database: database{
			URL: "sqlite://localhost/data/metadata.db?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)",

			ConnMaxIdleTime: 5 * time.Minute,
			ConnMaxLifetime: 30 * time.Minute,
			MaxOpenConns:    1, // SQLite typically benefits from single writer
			MaxIdleConns:    1,
		},
		Redis: redis{
			URL: "redis://localhost:6379/0",
		},
	}
}

func New() (Config, error) {
	cfg := Default()

	options := []config.Option{}
	if yamlPath := os.Getenv("CONFIG_PATH"); yamlPath != "" {
		options = append(options, config.WithLocalYAML(yamlPath))
	}

	if err := config.Load(&cfg, options...); err != nil {
		return Config{}, fmt.Errorf("failed to load config: %w", err)
	}

	return cfg, nil
}
