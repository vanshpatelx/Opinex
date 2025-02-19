# src/main.py
from fastapi import FastAPI
from routes.order_routes import router as order_router
from utils.logger import logger
from utils.init import ServiceInitializer
from config.config import config
import uvicorn
from config.db.db import Database


app = FastAPI(title="Order API")

app.include_router(order_router)

@app.get("/")
async def root():
    return {"message": "🚀 Order Service is Running!"}

@app.get("/order/health")
async def health_check():
    return {"success": True, "message": "Server is running."}

@app.on_event("startup")
async def startup():
    """Initialize all required services on FastAPI startup."""
    success = await ServiceInitializer.init()
    if not success:
        logger.error("❌ Service initialization failed. Exiting...")
        exit(1)

@app.on_event("shutdown")
async def shutdown():
    """Cleanup resources on shutdown."""
    await Database.close()
    logger.info("✅ Services cleaned up successfully.")

if __name__ == "__main__":
    uvicorn.run("main:app", host="0.0.0.0", port=config.PORT, reload=True)

