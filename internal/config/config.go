package config

import (
	"flag"
	"log/slog"

	"github.com/caarlos0/env/v11"
)

// Master config struct
type Config struct {
	RunAddress           string `env:"RUN_ADDRESS"`
	DatabaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	JWTSecret            string `env:"JWT_SECRET"`
	JWTExpireHours       int64  `env:"JWT_EXPIRE_HOURS"`
	Worker               int    `env:"WORKER_COUNT"`
}

// loads the configuration from environment variables or command-line flags or use defaults
func Load() *Config {
	config := Default()
	LoadFromEnv(config)
	LoadFromFlags(config)
	flag.Parse()
	return config
}

// Load default settings
func Default() *Config {
	config := &Config{}
	config.RunAddress = ":8080"
	config.DatabaseURI = "postgres://gofermart:gofermart@localhost:5432/gofermart?&sslmode=disable"
	config.AccrualSystemAddress = "http://localhost:8081"
	config.JWTSecret = "superSecret"
	config.JWTExpireHours = 96
	config.Worker = 3
	return config
}

// LoadFromEnv loads configuration from environment variables.
func LoadFromEnv(config *Config) {
	err := env.Parse(config)
	if err != nil {
		slog.Error("loads configuration from environment variables", "err", err)
	}
}

// LoadFromFlags loads configuration from command-line flags.
func LoadFromFlags(config *Config) {
	flag.StringVar(&config.RunAddress, "a", config.RunAddress, "Run address")
	flag.StringVar(&config.DatabaseURI, "d", config.DatabaseURI, "Data Base DSN")
	flag.StringVar(&config.AccrualSystemAddress, "r", config.AccrualSystemAddress, "Accrual System Address")
	flag.StringVar(&config.JWTSecret, "s", config.JWTSecret, "JWT Secret")
	flag.Int64Var(&config.JWTExpireHours, "e", config.JWTExpireHours, "JWT Expire Hours")
	flag.IntVar(&config.Worker, "w", config.Worker, "Number of Workers")
}
