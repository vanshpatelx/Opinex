# src/main.py
from fastapi import FastAPI
from src.routes.event_routes import router as event_router
from src.utils.logger import logger
from src.utils.init import ServiceInitializer
from src.config.config import config
import uvicorn
from src.config.db.db import Database

app = FastAPI(title="Event API")

@app.get("/event/health")
async def health_check():
    return {"success": True, "message": "🚀🚀 Event Server is running."}

app.include_router(event_router)

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
    logger.info(f"Starting server on port {config.PORT}")  # Debugging
    uvicorn.run("src.main:app", host="0.0.0.0", port=config.PORT, reload=False)
