# src/config/config.py
import os
from dotenv import load_dotenv

load_dotenv()

class Config:
    def __init__(self):
        self.REDIS_HOST = os.getenv("REDIS_HOST")
        self.REDIS_PORT = int(os.getenv("REDIS_PORT"))
        self.REDIS_PASSWORD = os.getenv("REDIS_PASSWORD")

        self.JWT_SECRET = os.getenv("JWT_SECRET")

        self.DB_USER = os.getenv("DB_USER")
        self.DB_HOST = os.getenv("DB_HOST")
        self.DB_NAME = os.getenv("DB_NAME")
        self.DB_PASSWORD = os.getenv("DB_PASSWORD")
        self.DB_PORT = int(os.getenv("DB_PORT"))

        self.PORT = int(os.getenv("PORT"))

        self.RABBITMQ_USER = os.getenv("RABBITMQ_USER")
        self.RABBITMQ_PASSWORD = os.getenv("RABBITMQ_PASSWORD")
        self.RABBITMQ_HOST = os.getenv("RABBITMQ_HOST")
        self.RABBITMQ_PORT = int(os.getenv("RABBITMQ_PORT"))
        self.RABBITMQ_EXCHANGES = [
            exchange.strip() for exchange in os.getenv("RABBITMQ_EXCHANGES", "").split(",")
        ]

    @property
    def RABBITMQ_URL(self):
        return f"amqp://{self.RABBITMQ_USER}:{self.RABBITMQ_PASSWORD}@{self.RABBITMQ_HOST}:{self.RABBITMQ_PORT}"

    @property
    def DATABASE_URL(self):
        return f"postgresql://{self.DB_USER}:{self.DB_PASSWORD}@{self.DB_HOST}:{self.DB_PORT}/{self.DB_NAME}"

    @property
    def REDIS_URL(self):
        return f"redis://{self.REDIS_HOST}:{self.REDIS_PORT}"

# Load configuration
config = Config()
