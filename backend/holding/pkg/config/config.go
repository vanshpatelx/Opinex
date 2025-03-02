/***
Package config - Application Configuration Loader

Handles the loading of environment variables and setting up configuration values 
for database connections, Redis, RabbitMQ, and JWT authentication.

Features:
1. Loads configuration from environment variables, optionally using a `.env` file.
2. Stores configuration values in a singleton `Config` struct.
3. Constructs connection URLs for PostgreSQL, Redis, and RabbitMQ.
4. Supports fallback default values where necessary.

Configuration Includes:
- **JWTSecret**: Secret key for JWT authentication.
- **User & Holding Service Databases**: PostgreSQL connection settings.
- **Redis Caches**: User and Holding service Redis connection settings.
- **RabbitMQ**: Messaging queue connection settings.
- **Server Port**: Defines the port for the application.

Last Updated: March 2, 2025
**/

package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	JWTSecret string

	USER_DBHost     string
	USER_DBPort     string
	USER_DBUser     string
	USER_DBPassword string
	USER_DBName     string

	HOLDING_DBHost     string
	HOLDING_DBPort     string
	HOLDING_DBUser     string
	HOLDING_DBPassword string
	HOLDING_DBName     string

	USER_RedisHost     string
	USER_RedisPort     string
	USER_RedisPassword string

	HOLDING_RedisHost     string
	HOLDING_RedisPort     string
	HOLDING_RedisPassword string

	RabbitMQUser     string
	RabbitMQPassword string
	RabbitMQHost     string
	RabbitMQPort     string

	USER_DBURL    string
	HOLDING_DBURL string

	USER_RedisURL    string
	HOLDING_RedisURL string

	RabbitMQURL string
	Port        string
}

var AppConfig Config

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	AppConfig = Config{
		JWTSecret: getEnv("JWT_SECRET", ""),

		USER_DBHost:     getEnv("USER_DB_HOST", ""),
		USER_DBPort:     getEnv("USER_DB_PORT", ""),
		USER_DBUser:     getEnv("USER_DB_USER", ""),
		USER_DBPassword: getEnv("USER_DB_PASSWORD", ""),
		USER_DBName:     getEnv("USER_DB_NAME", ""),

		HOLDING_DBHost:     getEnv("HOLDING_DB_HOST", ""),
		HOLDING_DBPort:     getEnv("HOLDING_DB_PORT", ""),
		HOLDING_DBUser:     getEnv("HOLDING_DB_USER", ""),
		HOLDING_DBPassword: getEnv("HOLDING_DB_PASSWORD", ""),
		HOLDING_DBName:     getEnv("HOLDING_DB_NAME", ""),

		USER_RedisHost:     getEnv("USER_REDIS_HOST", ""),
		USER_RedisPort:     getEnv("USER_REDIS_PORT", ""),
		USER_RedisPassword: getEnv("USER_REDIS_PASSWORD", ""),

		HOLDING_RedisHost:     getEnv("HOLDING_REDIS_HOST", ""),
		HOLDING_RedisPort:     getEnv("HOLDING_REDIS_PORT", ""),
		HOLDING_RedisPassword: getEnv("HOLDING_REDIS_PASSWORD", ""),

		RabbitMQUser:     getEnv("RABBITMQ_USER", ""),
		RabbitMQPassword: getEnv("RABBITMQ_PASSWORD", ""),
		RabbitMQHost:     getEnv("RABBITMQ_HOST", ""),
		RabbitMQPort:     getEnv("RABBITMQ_PORT", ""),

		Port: getEnv("PORT", "5001"),
	}

	AppConfig.JWTSecret = AppConfig.JWTSecret;

	AppConfig.USER_DBURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		AppConfig.USER_DBUser, AppConfig.USER_DBPassword, AppConfig.USER_DBHost, AppConfig.USER_DBPort, AppConfig.USER_DBName)

	AppConfig.HOLDING_DBURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		AppConfig.HOLDING_DBUser, AppConfig.HOLDING_DBPassword, AppConfig.HOLDING_DBHost, AppConfig.HOLDING_DBPort, AppConfig.HOLDING_DBName)

	AppConfig.USER_RedisURL = fmt.Sprintf("%s:%s", AppConfig.USER_RedisHost, AppConfig.USER_RedisPort)
	
	AppConfig.HOLDING_RedisURL = fmt.Sprintf("%s:%s", AppConfig.HOLDING_RedisHost, AppConfig.HOLDING_RedisPort)

	AppConfig.RabbitMQURL = fmt.Sprintf("amqp://%s:%s@%s:%s",
		AppConfig.RabbitMQUser, AppConfig.RabbitMQPassword, AppConfig.RabbitMQHost, AppConfig.RabbitMQPort)
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
