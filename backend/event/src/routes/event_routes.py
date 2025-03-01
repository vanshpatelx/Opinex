# src/routes/order_routes.py

"""
Event Routes

Defines API routes for event-related operations, including:
1. Create a new event.
2. Retrieve a specific event by EventID.
3. Retrieve all events with pagination.
4. Settle an event by EventID.

Author: Vansh Patel (remotevansh@gmail.com)
Last Updated: February 27, 2025

Related Controllers:
- `create_event`
- `get_all_event`
- `get_specific_event`
- `settlement_event`
"""

from fastapi import APIRouter, Depends
from src.models.order_model import EventReq
from src.middleware.jwt import get_current_user
from src.controllers.event_controller import create_event, get_all_event, get_specific_event, settlement_event
from src.models.auth_model import UserPayload
from typing import Dict, Any, List

router = APIRouter(prefix="/event")

@router.post("", response_model=dict)
async def create_new_event(event: EventReq, user: dict = Depends(get_current_user)):
    """
    API Endpoint: Create a new event.

    Args:
        event (EventReq): Event request payload.
        user (dict): Authenticated user details (extracted from JWT).

    Returns:
        dict: Confirmation message.

    Last Updated: February 27, 2025
    """
    user_payload = UserPayload(**user)  # Convert dictionary to Pydantic model
    return await create_event(event, user_payload)

@router.get("/", response_model=Dict[str, Any])
async def get_event_by_id(eventID: int, user: dict = Depends(get_current_user)):
    """
    API Endpoint: Retrieve an event by EventID.

    Args:
        eventID (int): The unique identifier of the event (query parameter).
        user (dict): Authenticated user details (extracted from JWT).

    Returns:
        Dict[str, Any]: Event details.

    Last Updated: February 27, 2025
    """
    user_payload = UserPayload(**user)  # Convert dictionary to Pydantic model
    return await get_specific_event(eventID, user_payload)


@router.get("/events", response_model=Dict[str, Any])
async def get_all_events(
    cursor: int = 0, limit: int = 50, user: dict = Depends(get_current_user)
):
    """
    API Endpoint: Retrieve all events with pagination.

    Args:
        user (dict): Authenticated user details (extracted from JWT).
        cursor (int): By default zero.
        limit (int): Number of events to retrieve.

    Returns:
        Dict[str, Any]: Response containing events.
    """
    return await get_all_event(cursor, limit, UserPayload(**user))


@router.post("/settlement/{eventID}", response_model=Dict[str, Any])
async def settle_event(eventID: int, user: dict = Depends(get_current_user)):
    """
    API Endpoint: Settle an event by EventID.

    Args:
        eventID (int): The unique identifier of the event.
        user (dict): Authenticated user details (extracted from JWT).

    Returns:
        Dict[str, Any]: Event details.

    Last Updated: February 27, 2025
    """
    user_payload = UserPayload(**user)  # Convert dictionary to Pydantic model
    return await settlement_event(eventID, user_payload)
