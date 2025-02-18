# src/utils/init.py

"""
Service Initializer

Initializes all required services:
1. Redis cache.
2. PostgreSQL database.
3. RabbitMQ (PubSub) connection.

Ensures all services are ready before starting the application.

Author: Vansh Patel (remotevansh@gmail.com)
Last Updated: February 19, 2025
"""

from utils.logger import logger
from config.config import config
from config.db.db import Database
from config.cache.cache import Cache
import aio_pika

class ServiceInitializer:
    @staticmethod
    async def init():
        """
        Initialize and verify all required services.

        - Attempts to initialize Redis, PostgreSQL, and RabbitMQ.
        - If any service fails, logs an error and exits early.
        
        Returns:
            bool: True if all services are initialized, False otherwise.

        Last Updated: February 19, 2025
        """
        
        # Initialize Redis
        if not await Cache.init():
            logger.error("❌ Redis is not ready after retries.")
            return False
        
        # Initialize Database
        if not await Database.init():
            logger.error("❌ Database is not ready after retries.")
            return False
        
        # Initialize RabbitMQ (PubSub)
        try:
            connection = await aio_pika.connect_robust(config.RABBITMQ_URL)
            async with connection:
                logger.info("✅ PubSub (RabbitMQ) is ready.")
        except Exception as e:
            logger.error(f"❌ PubSub failed to initialize: {e}")
            return False
        
        logger.info("✅ All services (Redis, DB, PubSub) are ready!")
        return True
