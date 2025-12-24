package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	RunAdress            string `env:"RUN_ADDRESS"`
	DatabaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	JWTSecret            string `env:"JWT_SECRET"`
	JWTExpireHours       int64  `env:"JWT_EXPIRE_HOURS"`
	Worker               int    `env:"WORKER_COUNT"`
}

// Load loads the configuration from environment variables or command-line flags.
func Load() *Config {
	config := LoadFromEnv()
	if config.RunAdress == "" || config.DatabaseURI == "" || config.AccrualSystemAddress == "" {
		config = LoadFromFlags()
		flag.Parse()
	}
	return config
}

// LoadFromEnv loads configuration from environment variables.
func LoadFromEnv() *Config {
	config := &Config{}
	err := env.Parse(config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

// LoadFromFlags loads configuration from command-line flags.
func LoadFromFlags() *Config {
	config := &Config{}
	flag.StringVar(&config.RunAdress, "a", ":8080", "Run address")
	// flag.StringVar(&config.DatabaseURI, "d", "postgres://admin:admin@192.168.0.20:5432/demo?search_path=gofermart&sslmode=disable", "Data Base DSN")
	flag.StringVar(&config.DatabaseURI, "d", "postgres://gofermart:gofermart@127.0.0.1:5432/gofermart?&sslmode=disable", "Data Base DSN")
	flag.StringVar(&config.AccrualSystemAddress, "r", "http://127.0.0.1:8081", "Accrual System Address")
	flag.StringVar(&config.JWTSecret, "s", "supersecret", "JWT Secret")
	flag.Int64Var(&config.JWTExpireHours, "e", 72, "JWT Expire Hours")
	flag.IntVar(&config.Worker, "w", 5, "Number of Workers")
	return config
}
