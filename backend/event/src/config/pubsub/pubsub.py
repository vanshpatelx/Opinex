# src/config/pubsub/pubsub.py

"""
RabbitMQ Pub/Sub Handler

Manages RabbitMQ connections for publishing messages.

Features:
1. Initializes a persistent connection to RabbitMQ.
2. Publishes messages to specified exchanges with routing keys.
3. Handles reconnection attempts for robustness.
4. Supports graceful shutdown of connections.

Author: Vansh Patel (remotevansh@gmail.com)
Last Updated: February 19, 2025
"""

import aio_pika
import asyncio
from src.utils.logger import logger
from src.config.config import config

class PubSub:
    """
    Manages RabbitMQ connections for sending messages asynchronously.

    Attributes:
        _connection_task (asyncio.Task | None): Stores the connection task to prevent multiple reconnections.
        connection (aio_pika.Connection | None): The active RabbitMQ connection.
        channel (aio_pika.Channel | None): The active RabbitMQ channel.
        reconnect_attempts (int): Maximum retry attempts for reconnection.
        reconnect_delay (int): Delay between retry attempts in seconds.
    """

    _connection_task = None
    connection = None
    channel = None
    reconnect_attempts = 5
    reconnect_delay = 5  # in seconds

    @classmethod
    async def init(cls):
        """
        Initialize RabbitMQ connection and channel (only once).
        
        - Ensures that the connection is established only once.
        - Uses an `asyncio.Task` to prevent multiple reconnection attempts.
        - Retries connection up to `reconnect_attempts` times.

        Last Updated: February 19, 2025
        """
        if cls.connection and cls.channel:
            return  # ✅ Already initialized, no need to reconnect

        if cls._connection_task:  # ✅ Wait for the existing connection task if it's initializing
            await cls._connection_task
            return

        async def _connect():
            attempt = 0
            while attempt < cls.reconnect_attempts:
                try:
                    cls.connection = await aio_pika.connect_robust(config.RABBITMQ_URL)
                    cls.channel = await cls.connection.channel()  # Create a channel
                    logger.info("✅ Connected to RabbitMQ (PubSub).")
                    return
                except Exception as e:
                    attempt += 1
                    logger.warning(f"⚠️ Failed to connect to RabbitMQ (attempt {attempt}/{cls.reconnect_attempts}): {e}")
                    await asyncio.sleep(cls.reconnect_delay)

            logger.error("❌ Could not connect to RabbitMQ after multiple attempts.")

        # ✅ Use asyncio.Task to ensure only one connection attempt is active at a time
        cls._connection_task = asyncio.create_task(_connect())  
        await cls._connection_task

    @classmethod
    async def send(cls, message: dict, exchange: str, routing_key: str):
        """
        Send a message to RabbitMQ.

        - Ensures RabbitMQ connection is initialized before sending.
        - Declares an exchange before publishing messages.
        - Supports direct exchanges with routing keys.

        Args:
            message (dict): The message payload to send.
            exchange (str): The RabbitMQ exchange name.
            routing_key (str): The routing key for the message.

        Returns:
            None
        
        Last Updated: February 19, 2025
        """
        await cls.init()  # ✅ Ensure connection is initialized

        try:
            # Get the existing exchange without declaring it
            exchange_obj = await cls.channel.get_exchange(exchange)

            # Publish message
            await exchange_obj.publish(
                aio_pika.Message(body=str(message).encode()),
                routing_key=routing_key
            )

            logger.info(f"📩 Message sent to exchange '{exchange}' with routing key '{routing_key}': {message}")

        except Exception as e:
            logger.error(f"❌ Failed to send message to exchange '{exchange}' with routing key '{routing_key}': {e}")

    @classmethod
    async def close(cls):
        """
        Close RabbitMQ connection and channel gracefully.

        - Ensures channels and connections are properly closed.
        - Logs confirmation upon successful closure.

        Last Updated: February 19, 2025
        """
        if cls.channel:
            await cls.channel.close()
            cls.channel = None
            logger.info("✅ RabbitMQ channel closed.")
        if cls.connection:
            await cls.connection.close()
            cls.connection = None
            logger.info("✅ RabbitMQ connection closed.")
