package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddress    string
	PostgresConn     string
	PostgresJDBCURL  string
	PostgresUsername string
	PostgresPassword string
	PostgresHost     string
	PostgresPort     int
	PostgresDatabase string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found. Proceeding with environment variables.")
	} else {
		log.Println(".env file loaded successfully.")
	}

	cfg := &Config{
		ServerAddress:    os.Getenv("SERVER_ADDRESS"),
		PostgresConn:     os.Getenv("POSTGRES_CONN"),
		PostgresJDBCURL:  os.Getenv("POSTGRES_JDBC_URL"),
		PostgresUsername: os.Getenv("POSTGRES_USERNAME"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresHost:     os.Getenv("POSTGRES_HOST"),
		PostgresDatabase: os.Getenv("POSTGRES_DATABASE"),
	}

	portStr := os.Getenv("POSTGRES_PORT")
	if portStr == "" {
		log.Println("POSTGRES_PORT not set, using default port 5432")
		cfg.PostgresPort = 5432
	} else {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("invalid POSTGRES_PORT: %v", err)
		}
		cfg.PostgresPort = port
	}

	if cfg.ServerAddress == "" {
		cfg.ServerAddress = ":8080"
	}

	if cfg.PostgresConn == "" && (cfg.PostgresHost == "" || cfg.PostgresUsername == "" || cfg.PostgresPassword == "" || cfg.PostgresDatabase == "") {
		return nil, fmt.Errorf("either POSTGRES_CONN or all of POSTGRES_HOST, POSTGRES_USERNAME, POSTGRES_PASSWORD, and POSTGRES_DATABASE must be set")
	}

	log.Printf("ServerAddress: %s", cfg.ServerAddress)
	log.Printf("PostgresHost: %s", cfg.PostgresHost)
	log.Printf("PostgresPort: %d", cfg.PostgresPort)
	log.Printf("PostgresDatabase: %s", cfg.PostgresDatabase)

	return cfg, nil
}
