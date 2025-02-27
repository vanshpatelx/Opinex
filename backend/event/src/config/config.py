# src/config/config.py

"""
Configuration Loader

Loads environment variables and provides application-wide configurations.

Features:
1. Loads `.env` variables using `dotenv`.
2. Stores Redis, JWT, Database, and RabbitMQ settings.
3. Provides computed properties for service URLs.

Author: Vansh Patel (remotevansh@gmail.com)
Last Updated: February 19, 2025
"""

import os
from dotenv import load_dotenv

load_dotenv()

class Config:
    """
    Application configuration settings.

    Loads and stores environment variables for:
    - Redis
    - JWT Authentication
    - PostgreSQL Database
    - RabbitMQ (PubSub)
    """
    def __init__(self):
        # Redis settings
        self.REDIS_HOST = os.getenv("REDIS_HOST")
        self.REDIS_PORT = int(os.getenv("REDIS_PORT"))
        self.REDIS_PASSWORD = os.getenv("REDIS_PASSWORD")

        # JWT secret key
        self.JWT_SECRET = os.getenv("JWT_SECRET")

        # PostgreSQL database settings
        self.DB_USER = os.getenv("DB_USER")
        self.DB_HOST = os.getenv("DB_HOST")
        self.DB_NAME = os.getenv("DB_NAME")
        self.DB_PASSWORD = os.getenv("DB_PASSWORD")
        self.DB_PORT = int(os.getenv("DB_PORT"))

        # Application port
        self.PORT = int(os.getenv("PORT"))

        # RabbitMQ settings
        self.RABBITMQ_USER = os.getenv("RABBITMQ_USER")
        self.RABBITMQ_PASSWORD = os.getenv("RABBITMQ_PASSWORD")
        self.RABBITMQ_HOST = os.getenv("RABBITMQ_HOST")
        self.RABBITMQ_PORT = int(os.getenv("RABBITMQ_PORT"))
        self.RABBITMQ_EXCHANGES = [
            exchange.strip() for exchange in os.getenv("RABBITMQ_EXCHANGES", "").split(",")
        ]

    @property
    def RABBITMQ_URL(self):
        """Returns RabbitMQ connection URL."""
        return f"amqp://{self.RABBITMQ_USER}:{self.RABBITMQ_PASSWORD}@{self.RABBITMQ_HOST}:{self.RABBITMQ_PORT}"

    @property
    def DATABASE_URL(self):
        """Returns PostgreSQL connection URL."""
        return f"postgresql://{self.DB_USER}:{self.DB_PASSWORD}@{self.DB_HOST}:{self.DB_PORT}/{self.DB_NAME}"

    @property
    def REDIS_URL(self):
        """Returns Redis connection URL."""
        return f"redis://{self.REDIS_HOST}:{self.REDIS_PORT}"

# Load configuration
config = Config()
