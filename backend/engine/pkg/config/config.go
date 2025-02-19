package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Load environment variables
func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

// Config struct to store environment variables
var Config = struct {
	Redis struct {
		Host     string
		Port     string
		Password string
	}
	JWTSecret struct {
		Secret string
	}
	DB struct {
		User     string
		Host     string
		Database string
		Password string
		Port     string
	}
	Port     string
	RabbitMQ struct {
		User      string
		Password  string
		Host      string
		Port      string
		URL       string
		Exchanges []string
	}
}{
	Redis: struct {
		Host     string
		Port     string
		Password string
	}{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
	},
	JWTSecret: struct {
		Secret string
	}{
		Secret: os.Getenv("JWT_SECRET"),
	},
	DB: struct {
		User     string
		Host     string
		Database string
		Password string
		Port     string
	}{
		User:     os.Getenv("DB_USER"),
		Host:     os.Getenv("DB_HOST"),
		Database: os.Getenv("DB_NAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Port:     os.Getenv("DB_PORT"),
	},
	Port: os.Getenv("PORT"),
	RabbitMQ: struct {
		User      string
		Password  string
		Host      string
		Port      string
		URL       string
		Exchanges []string
	}{
		User:      os.Getenv("RABBITMQ_USER"),
		Password:  os.Getenv("RABBITMQ_PASSWORD"),
		Host:      os.Getenv("RABBITMQ_HOST"),
		Port:      os.Getenv("RABBITMQ_PORT"),
		URL:       "amqp://" + os.Getenv("RABBITMQ_USER") + ":" + os.Getenv("RABBITMQ_PASSWORD") + "@" + os.Getenv("RABBITMQ_HOST") + ":" + os.Getenv("RABBITMQ_PORT"),
		Exchanges: []string{"auth_exchange", "dlx_exchange"},
	},
}
