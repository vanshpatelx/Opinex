import aio_pika
import asyncio
from utils.logger import logger
from config.config import config

class PubSub:
    _connection_task = None  # ✅ Store RabbitMQ connection in a global task
    connection = None
    channel = None
    reconnect_attempts = 5
    reconnect_delay = 5  # in seconds

    @classmethod
    async def init(cls):
        """Initialize RabbitMQ connection and channel (only once)."""
        if cls.connection and cls.channel:
            return True # ✅ Already initialized, no need to reconnect

        if cls._connection_task:  # ✅ Wait for the existing connection task if it's initializing
            await cls._connection_task
            return True

        async def _connect():
            attempt = 0
            while attempt < cls.reconnect_attempts:
                try:
                    cls.connection = await aio_pika.connect_robust(config.RABBITMQ_URL)
                    cls.channel = await cls.connection.channel()  # Create a channel
                    logger.info("✅ Connected to RabbitMQ (PubSub).")
                    return True
                except Exception as e:
                    attempt += 1
                    logger.warning(f"⚠️ Failed to connect to RabbitMQ (attempt {attempt}/{cls.reconnect_attempts}): {e}")
                    await asyncio.sleep(cls.reconnect_delay)
                    return False

            logger.error("❌ Could not connect to RabbitMQ after multiple attempts.")

        cls._connection_task = asyncio.create_task(_connect())  # ✅ Use asyncio.Task for a single connection
        await cls._connection_task

    @classmethod
    async def send(cls, message: dict, exchange: str, routing_key: str):
        """Send a message to RabbitMQ."""
        await cls.init()  # ✅ Ensure connection is initialized only if needed

        try:
            # Ensure exchange exists
            exchange_instance = await cls.channel.declare_exchange(
                exchange, aio_pika.ExchangeType.DIRECT, durable=True
            )

            # Publish message
            await exchange_instance.publish(
                aio_pika.Message(body=str(message).encode()),
                routing_key=routing_key
            )
            logger.info(f"📩 Message sent to exchange '{exchange}' with routing key '{routing_key}': {message}")

        except Exception as e:
            logger.error(f"❌ Failed to send message to exchange '{exchange}' with routing key '{routing_key}': {e}")

    @classmethod
    async def close(cls):
        """Close RabbitMQ connection and channel."""
        if cls.channel:
            await cls.channel.close()
            cls.channel = None
            logger.info("✅ RabbitMQ channel closed.")
        if cls.connection:
            await cls.connection.close()
            cls.connection = None
            logger.info("✅ RabbitMQ connection closed.")
