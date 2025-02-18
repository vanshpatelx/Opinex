# from fastapi import FastAPI

# app = FastAPI()

# @app.get("/")
# def read_root():
#     return {"message": "Hello, FastAPI!"}

from fastapi import FastAPI
from routes.order_routes import router as order_router
from utils.logger import logger
from utils.init import ServiceInitializer
from config.config import config
from config.load_env import load_env
import uvicorn

app = FastAPI(title="Order API")

app.include_router(order_router)

@app.get("/")
async def root():
    return {"message": "🚀 Order Service is Running!"}

@app.get("/order/health")
async def health_check():
    return {"success": True, "message": "Server is running."}

async def start_server():
    try:
        load_env()
        await ServiceInitializer.init()
        logger.info(f"Server is running on port {config.PORT}")
    except Exception as e:
        logger.error(f"❌ Failed to start server: {e}")
        exit(1)

if __name__ == "__main__":
    import asyncio
    asyncio.run(start_server())
    uvicorn.run(app, host="0.0.0.0", port=config.PORT)



# import asyncio
# from config.db.db import Database
# from config.cache.cache import Cache
# from utils.logger import logger  # Central logging

# async def test_db():
#     """Test PostgreSQL connection."""
#     db = Database()
#     conn = await db.get_pool()
#     logger.info("PostgreSQL Connected Successfully!")

#     # Create a dummy table
#     await conn.execute("""
#         CREATE TABLE IF NOT EXISTS test_table (
#             id SERIAL PRIMARY KEY,
#             name TEXT NOT NULL
#         )
#     """)
#     logger.info("Table Created")

#     # Insert Data (Insert 10000 rows in batches of 1000)
#     batch_size = 1000
#     for i in range(0, 10000, batch_size):
#         records = [("Test User " + str(i + j),) for j in range(1, batch_size + 1)]
#         await conn.executemany("INSERT INTO test_table (name) VALUES ($1)", records)
#         logger.info(f"Inserted batch {i // batch_size + 1} of {batch_size} Test Users")

#     # Fetch Data (Fetch 10000 rows in batches of 1000)
#     rows = await conn.fetch("SELECT * FROM test_table LIMIT 10000")
#     logger.info(f"DB Data: {len(rows)} rows fetched")

# async def test_redis():
#     """Test Redis connection."""
#     # Insert 10000 keys into Redis in batches of 1000
#     batch_size = 1000
#     for i in range(1, 10000, batch_size):
#         batch = {f"test_key_{i + j}": {"message": f"Hello from Redis {i + j}!"} for j in range(batch_size)}
#         # Set keys in batch
#         for key, value in batch.items():
#             await Cache.set(key, value)
#         logger.info(f"Inserted batch {i // batch_size + 1} of {batch_size} keys into Redis")

#     # Fetch 10000 keys from Redis in batches of 1000
#     for i in range(1, 10000, batch_size):
#         for j in range(batch_size):
#             key = f"test_key_{i + j}"
#             data = await Cache.get(key)
#             # if data:
#             #     logger.info(f"Redis Data {i + j}: {data}")
#             # else:
#             #     logger.warning(f"Redis Data {i + j} not found")

# async def main():
#     # Ensure database and cache are initialized properly
#     await Cache.init()
#     await test_db()
#     await test_redis()

# if __name__ == "__main__":
#     asyncio.run(main())


# main.py

# import asyncio
# from config.pubsub.pubsub import PubSub
# from utils.logger import logger

# async def test_pubsub():
#     queue = "test_queue"
#     await PubSub.create_queue()  # Ensure the queue exists
#     for i in range(10000):  # Loop to send 1000 messages
#         message = {"test": f"connection {i}"}  # Create a unique message for each iteration
#         await PubSub.send(message, queue)
#         logger.info(f"✅ Sent message {i} to {queue}")

#     logger.info("✅ PubSub connection test completed!")


# if __name__ == "__main__":
#     asyncio.run(test_pubsub())

# from config.config import config
# from config.load_env import load_env


# def main():
#     load_env()
#     print(f"Starting application on port {config.PORT}")
#     print(f"Connecting to database: {config.DATABASE_URL}")
#     print(f"RabbitMQ URL: {config.RABBITMQ_URL}")

# if __name__ == "__main__":
#     main()
