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
from src.config.cache.cache import Cache
from src.config.pubsub.pubsub import PubSub
from src.models.order_model import EventReq, Event, EventResponse
from src.models.auth_model import UserPayload
from src.utils.id import generate_unique_id
from src.utils.logger import logger
from datetime import datetime
from src.config.db.db import Database
from src.config.db.query import Query
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
    if user.type != "ADMIN":
        logger.warning(f"{user.type}")
        logger.warning(f"🚫 Unauthorized access: User {user.id} attempted to create an event.")
        raise HTTPException(status_code=401, detail="Unauthorized")

    created_at = datetime.utcnow().strftime("%Y-%m-%d %H:%M:%S")

    logger.warn(str(event.settlement_time))

    event_obj = Event(
        id=generate_unique_id(),
        name=event.name,
        details=event.details,
        status="RUNNING",
        settlement_time=event.settlement_time,
        created_at=created_at,
    )

    event_dict = event_obj.model_dump()
    event_dict["id"] = str(event_obj.id)  # Convert ID to string format

    logger.info(f"📦 Final Event: {event_dict}")

    logger.info(event_dict)
    await PubSub.send(event_dict, 'event_exchange', 'event.registered')

    # Minimize data before caching
    cache_event = event_dict.copy()
    cache_event.pop("updated_at", None)
    cache_event.pop("created_at", None)
    cache_event["settlement_time"] = str(cache_event["settlement_time"])

    await Cache.set(f"event:{event_obj.id}", cache_event, 72000)
    logger.info(f"📦 Cached Event: {event_obj.id}")

    logger.info(f"✅ Event {event_obj.id} created successfully")
    return {"message": "Event created successfully", "event_id": event_obj.id}

async def get_all_event(cursor: int, limit: int, user: UserPayload):
    """
    Retrieves all events, checking the cache first.

    Args:
        cursor (int): Pagination cursor.
        limit (int): Number of events to retrieve.
        user (UserPayload): User details.

    Returns:
        List[Dict[str, Any]]: Retrieved events.
    """
    key = f"live_events:{cursor}:{limit}"
    
    if (cached_events := await Cache.get(key)):
        logger.info("✅ Cache hit: Events retrieved from Redis.")
        return {"message": "Event fetched successfully", "data": cached_events}


    logger.info(f"📥 Fetching Live Events for User {user.id} ({user.type})")

    # Fetch events from the database
    events = await Database.fetch_all(Query.get_all_event, limit, cursor)  # Fixed order
    if not events:
        logger.warning("❌ No events found in the database.")
        raise HTTPException(status_code=404, detail="Events not found")

    # Convert event records into response models
    events_list = [
        EventResponse(
            id=str(event["id"]),
            name=event["name"],
            details=event["details"],
            status=event["status"],
            settlement_time=str(event["settlement_time"]),  # Convert datetime to string
        ).model_dump()
        for event in events
    ]

    # Store in cache for 60 seconds
    await Cache.set(key, events_list, 10)
    logger.info(f"✅ Cached {len(events_list)} events for future requests.")

    return {"message": "Event fetched successfully", "data": events_list}


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
        return {"message": "Event fetched sucessfully", "data": results}

    logger.info(f"📥 Request: Fetching Event {eventID} for User {user.id} ({user.type})")

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
    if user.type != "ADMIN":
        logger.warning(f"🚫 Unauthorized access: User {user.id} attempted to settle an event.")
        raise HTTPException(status_code=401, detail="Unauthorized")

    try:
        eventID = int(eventID)
    except ValueError:
        raise HTTPException(status_code=400, detail="Invalid event ID")

    key = f"event:{eventID}"
    event = await Cache.get(key)

    if event is None:
        logger.info(f"📥 Fetching Event {eventID} for User {user.id} ({user.type})")

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
