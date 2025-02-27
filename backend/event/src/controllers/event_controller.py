# src/controllers/event_controller.py

"""
Event Controller

Handles event-related operations, including:
1. Create Event: `create_event`
2. Retrieve All Events: `get_all_event`
3. Retrieve Specific Event by EventID: `get_specific_event`
4. Settle Event by EventID: `settlement_event`

Author: Vansh Patel (remotevansh@gmail.com)
Last Updated: February 27, 2025

Related Models:
- `EventReq`
- `Event`
- `UserPayload`
"""

from fastapi import HTTPException
from config.cache.cache import Cache
from config.pubsub.pubsub import PubSub
from models.order_model import EventReq, Event
from models.auth_model import UserPayload
from utils.id import generate_unique_id
from utils.logger import logger
from datetime import datetime
from config.db.db import Database
from config.db.query import Query
import logging

logger = logging.getLogger(__name__)

async def create_event(event: EventReq, user: UserPayload):
    """
    Creates and processes a new event.

    Steps:
        1. Validate user permissions.
        2. Generate Event ID.
        3. Store event in cache and publish for DB.
        4. Publish to different exchanges.

    Parameters:
        event (EventReq): Event details.
        user (UserPayload): User details.

    Returns:
        dict: Event creation confirmation.
    """
    if user.user_type != "ADMIN":
        logger.warning(f"🚫 Unauthorized access: User {user.UserID} attempted to create an event.")
        raise HTTPException(status_code=401, detail="Unauthorized")

    event_obj = Event(
        id=generate_unique_id(),
        name=event.name,
        details=event.details,
        status="RUNNING",
        settlement_time=event.settlement_time,
        updated_at=int(datetime.utcnow().timestamp()),
        created_at=int(datetime.utcnow().timestamp()),
    )

    event_dict = event_obj.model_dump()
    event_dict["id"] = str(event_obj.id)  # Convert ID to string format

    logger.info(f"📦 Final Event: {event_dict}")

    await PubSub.send(event_dict, 'Order_Exchange', 'Order.add')

    # Minimize data before caching
    cache_event = event_dict.copy()
    cache_event.pop("updated_at", None)
    cache_event.pop("created_at", None)

    await Cache.set(f"event:{event_obj.id}", cache_event)
    logger.info(f"📦 Cached Event: {event_obj.id}")

    logger.info(f"✅ Event {event_obj.id} created successfully")
    return {"message": "Event created successfully", "event_id": event_obj.id}

async def get_all_event(user: UserPayload):
    """
    Retrieves all events.

    Steps:
        1. Check if events exist in cache.
        2. If not, fetch from DB and store in cache.
        3. Return the list of events.

    Parameters:
        user (UserPayload): User details.

    Returns:
        dict: A dictionary containing the retrieved events.
    """
    key = "live_events"

    results = await Cache.get(key)
    if results is not None:
        logger.info(f"✅ Cache hit: Live Events retrieved from Redis.")
        return {"message": "Live events fetched from cache", "data": results}

    logger.info(f"📥 Request: Fetching Live Events for User {user.UserID} ({user.user_type})")

    try:
        events = await Database.fetch_all(Query.get_all_event)
        events = events if events else []

        await Cache.set(key, events, ex=60)
        logger.info(f"✅ All Events cached successfully")

    except Exception as e:
        logger.error(f"❌ Error fetching live events: {str(e)}")
        raise HTTPException(status_code=500, detail="Internal Server Error")

    return {"message": "All events fetched successfully", "data": events}

async def get_specific_event(eventID: str, user: UserPayload):
    """
    Retrieves a specific event by EventID.

    Steps:
        1. Check if event exists in cache.
        2. If not, fetch from DB and store in cache.
        3. Return the event.

    Parameters:
        eventID (str): Event ID.
        user (UserPayload): User details.

    Returns:
        dict: Retrieved event.
    """
    try:
        eventID = int(eventID)
    except ValueError:
        raise HTTPException(status_code=400, detail="Invalid event ID")

    key = f"event:{eventID}"
    results = await Cache.get(key)

    if results:
        logger.info(f"✅ Cache hit: Event {eventID} retrieved from Redis.")
        return {"message": "Event fetched from cache", "data": results}

    logger.info(f"📥 Request: Fetching Event {eventID} for User {user.UserID} ({user.user_type})")

    try:
        event = await Database.fetch_one(Query.get_event, eventID)
        if not event:
            logger.warning(f"❌ Event {eventID} not found in database.")
            raise HTTPException(status_code=404, detail="Event not found")

        await Cache.set(key, event)
        logger.info(f"✅ Event {eventID} cached successfully")

    except Exception as e:
        logger.error(f"❌ Error fetching event {eventID}: {str(e)}")
        raise HTTPException(status_code=500, detail="Internal Server Error")

    return {"message": "Event fetched successfully", "data": event}

async def settlement_event(eventID: str, user: UserPayload):
    """
    Settles a specific event by EventID.

    Steps:
        1. Check if event exists in cache.
        2. If not, fetch from DB and store in cache.
        3. Update status to 'SETTLEMENT' in cache and publish for DB and more operations.

    Parameters:
        eventID (str): Event ID.
        user (UserPayload): User details.

    Returns:
        dict: Updated event.
    """
    if user.user_type != "ADMIN":
        logger.warning(f"🚫 Unauthorized access: User {user.UserID} attempted to settle an event.")
        raise HTTPException(status_code=401, detail="Unauthorized")

    try:
        eventID = int(eventID)
    except ValueError:
        raise HTTPException(status_code=400, detail="Invalid event ID")

    key = f"event:{eventID}"
    event = await Cache.get(key)

    if event is None:
        logger.info(f"📥 Fetching Event {eventID} for User {user.UserID} ({user.user_type})")

        try:
            event = await Database.fetch_one(Query.get_event, eventID)
            if not event:
                logger.warning(f"❌ Event {eventID} not found in database.")
                raise HTTPException(status_code=404, detail="Event not found")

            await Cache.set(key, event)
            logger.info(f"✅ Event {eventID} cached successfully")

        except Exception as e:
            logger.error(f"❌ Error fetching event {eventID}: {str(e)}")
            raise HTTPException(status_code=500, detail="Internal Server Error")

    event["status"] = "SETTLEMENT"

    await PubSub.send(event, 'Order_Exchange', 'Order.add')

    # Update cache
    await Cache.set(key, event)
    logger.info(f"✅ Event {eventID} settled successfully")

    return {"message": "Event settlement applied", "data": event}
