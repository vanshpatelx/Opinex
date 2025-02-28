# src/config/cache/cache.py
"""
Redis Cache Connection

Handles connection and operations for Redis caching.

Features:
1. Manages a Redis connection pool.
2. Supports GET and SET operations.
3. Handles auto-reconnect in case of failures.

Author: Vansh Patel (remotevansh@gmail.com)
Last Updated: February 26, 2025
"""

import redis
import json
import asyncio
from src.utils.logger import logger
from src.config.config import config

class Cache:
    redis = None
    pool = None
    reconnect_attempts = 5
    reconnect_delay = 5  # in seconds

    @classmethod
    async def init(cls):
        """Initialize the Redis connection pool."""
        attempt = 0
        while attempt < cls.reconnect_attempts:
            try:
                cls.pool = redis.ConnectionPool.from_url(config.REDIS_URL, password=config.REDIS_PASSWORD)
                cls.redis = redis.Redis(connection_pool=cls.pool)
                if cls.redis.ping():
                    logger.info("✅ Connected to Redis.")
                    return True
                else:
                    logger.warning("⚠️ Redis ping failed.")
                    raise Exception("Redis ping failed.")

            except Exception as e:
                attempt += 1
                logger.warning(f"⚠️ Failed to connect to Redis (attempt {attempt}/{cls.reconnect_attempts}): {e}")
                await asyncio.sleep(cls.reconnect_delay)

        logger.error("❌ Could not connect to Redis after multiple attempts.")

    @classmethod
    async def get_redis(cls):
        """Return Redis connection."""
        if cls.redis is None:
            await cls.init()
        return cls.redis

    @classmethod
    async def get(cls, key: str):
        """Get data from Redis."""
        redis = await cls.get_redis()
        data = redis.get(key)
        return json.loads(data) if data else None

    @classmethod
    async def set(cls, key: str, value: dict, expire=72000):
        """Set data in Redis."""
        redis = await cls.get_redis()
        redis.set(key, json.dumps(value), ex=expire)

    @classmethod
    async def close(cls):
        """Close the Redis connection."""
        if cls.redis:
            await cls.redis.close()
            logger.info("✅ Redis connection closed.")
        if cls.pool:
            await cls.pool.disconnect()
            logger.info("✅ Redis connection pool closed.")
