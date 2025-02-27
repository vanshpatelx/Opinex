# src/config/db/db.py

"""
PostgreSQL Database Connection

Manages connection pooling and queries for PostgreSQL using `asyncpg`.

Features:
1. Initializes a connection pool to PostgreSQL.
2. Provides methods for executing queries safely.
3. Supports graceful shutdown.

Author: Vansh Patel (remotevansh@gmail.com)
Last Updated: February 26, 2025
"""

import asyncpg
import asyncio
from utils.logger import logger
from config.config import config

class Database:
    """Handles PostgreSQL connection pooling and queries."""
    
    pool: asyncpg.Pool | None = None
    reconnect_attempts = 5
    reconnect_delay = 5  # in seconds

    @classmethod
    async def init(cls):
        """Initialize PostgreSQL database connection pool."""
        attempt = 0
        while attempt < cls.reconnect_attempts:
            try:
                loop = asyncio.get_running_loop()  # Ensure correct event loop
                cls.pool = await asyncpg.create_pool(dsn=config.DATABASE_URL, loop=loop)
                logger.info("✅ Connected to PostgreSQL.")
                return True
            except Exception as e:
                attempt += 1
                logger.warning(f"⚠️ Failed to connect to PostgreSQL (attempt {attempt}/{cls.reconnect_attempts}): {e}")
                await asyncio.sleep(cls.reconnect_delay)

        logger.error("❌ Could not connect to PostgreSQL after multiple attempts.")

    @classmethod
    async def get_pool(cls) -> asyncpg.Pool:
        """Return the database connection pool."""
        if cls.pool is None:
            await cls.init()
        return cls.pool

    @classmethod
    async def fetch_one(cls, query: str, *args):
        """Fetch a single row using the connection pool safely."""
        pool = await cls.get_pool()
        
        if pool is None:
            logger.error("❌ Database pool is not initialized!")
            return None

        try:
            async with pool.acquire() as conn:  # Use 'async with' to manage connection properly
                result = await conn.fetchrow(query, *args)
                logger.info(f"✅ Query Result: {result}")
                return result
        except Exception as e:
            logger.error(f"⚠️ DB fetch_one error: {e}")
            return None

    @classmethod
    async def fetch_all(cls, query: str, *args):
        """Fetch multiple rows from the database."""
        pool = await cls.get_pool()
        
        if pool is None:
            logger.error("❌ Database pool is not initialized!")
            return None

        try:
            async with pool.acquire() as conn:
                result = await conn.fetch(query, *args)
                logger.info(f"✅ Query Result: {result}")
                return result
        except Exception as e:
            logger.error(f"⚠️ DB fetch_all error: {e}")
            return None

    @classmethod
    async def close(cls):
        """Close the PostgreSQL connection pool."""
        if cls.pool:
            await cls.pool.close()
            logger.info("✅ PostgreSQL connection pool closed.")
