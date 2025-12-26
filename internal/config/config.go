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
	config := Default()
	LoadFromEnv(config)
	LoadFromFlags(config)
	flag.Parse()
	return config
}

// Load default settings
func Default() *Config {
	config := &Config{}
	config.RunAdress = ":8080"
	config.DatabaseURI = "postgres://gofermart:gofermart@127.0.0.1:5432/gofermart?&sslmode=disable" //  "postgres://admin:admin@192.168.0.20:5432/demo?search_path=gofermart&sslmode=disable"
	config.AccrualSystemAddress = "http://127.0.0.1:8081"
	config.JWTSecret = "secret"
	config.JWTExpireHours = 74
	config.Worker = 6
	return config
}

// LoadFromEnv loads configuration from environment variables.
func LoadFromEnv(config *Config) {
	err := env.Parse(config)
	if err != nil {
		log.Fatal(err)
	}
}

// LoadFromFlags loads configuration from command-line flags.
func LoadFromFlags(config *Config) {
	flag.StringVar(&config.RunAdress, "a", config.RunAdress, "Run address")
	flag.StringVar(&config.DatabaseURI, "d", config.DatabaseURI, "Data Base DSN")
	flag.StringVar(&config.AccrualSystemAddress, "r", config.AccrualSystemAddress, "Accrual System Address")
	flag.StringVar(&config.JWTSecret, "s", config.JWTSecret, "JWT Secret")
	flag.Int64Var(&config.JWTExpireHours, "e", config.JWTExpireHours, "JWT Expire Hours")
	flag.IntVar(&config.Worker, "w", config.Worker, "Number of Workers")
}
