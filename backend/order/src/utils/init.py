# src/utils/init.py
from utils.logger import logger
from config.config import config
from config.db.db import Database
from config.cache.cache import Cache
from config.pubsub.pubsub import PubSub

import aio_pika

class ServiceInitializer:
    @staticmethod
    async def init():
        """Initialize all services and check readiness."""
        
        # Try to initialize Redis
        if not await Cache.init():
            logger.error("❌ Redis is not ready after retries.")
            return False  # Exit early if Redis failed
        
        # Try to initialize Database
        if not await Database.init():
            logger.error("❌ Database is not ready after retries.")
            return False  # Exit early if Database failed
        
        # Try to initialize PubSub by checking the connection without sending a message
        if not await PubSub.init():
            logger.error("❌ PubSub is not ready after retries.")
            return False  # Exit early if Database failed

        # try:
        #     # Establish connection to RabbitMQ
        #     connection = await aio_pika.connect_robust(config.RABBITMQ_URL)
        #     async with connection:
        #         logger.info("✅ PubSub (RabbitMQ) is ready.")
        # except Exception as e:
        #     logger.error(f"❌ PubSub failed to initialize: {e}")
        #     return False
        
        logger.info("✅ All services (Redis, DB, PubSub) are ready!")
        return True
