# src/controllers/event_controller.py

from config.cache.cache import Cache

async def get_event_by_id(eventID: int) -> bool:
    key = f"Event:{eventID}"
    status = Cache.get(key)
    return bool(status)
