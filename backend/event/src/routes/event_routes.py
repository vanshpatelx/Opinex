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
from models.order_model import EventReq
from middleware.jwt import get_current_user
from src.controllers.event_controller import create_event, get_all_event, get_specific_event, settlement_event
from models.auth_model import UserPayload
from typing import Dict, Any

router = APIRouter(prefix="/event")

@router.post("/", response_model=dict)
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

@router.get("/{eventID}", response_model=Dict[str, Any])
async def get_event_by_id(eventID: int, user: dict = Depends(get_current_user)):
    """
    API Endpoint: Retrieve an event by EventID.

    Args:
        eventID (int): The unique identifier of the event.
        user (dict): Authenticated user details (extracted from JWT).

    Returns:
        Dict[str, Any]: Event details.

    Last Updated: February 27, 2025
    """
    user_payload = UserPayload(**user)  # Convert dictionary to Pydantic model
    return await get_specific_event(eventID, user_payload)

@router.get("s", response_model=Dict[str, Any])
async def get_all_events(user: dict = Depends(get_current_user)):
    """
    API Endpoint: Retrieve all events with pagination.

    Args:
        user (dict): Authenticated user details (extracted from JWT).

    Returns:
        Dict[str, Any]: List of events.

    Last Updated: February 27, 2025
    """
    user_payload = UserPayload(**user)  # Convert dictionary to Pydantic model
    return await get_all_event(user_payload)

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


# @router.get("/orders/event/{eventID}", response_model=List[Dict[str, Any]])
# async def orders_by_eventID(eventID: int, cursor: int = 0, limit: int = 50, user: dict = Depends(get_current_user)):
#     """
#     API Endpoint: Admin-Focused Order Retrieval (Paginated by EventID, Ordered by EventID)

#     Args:
#         eventID (int): The unique identifier of the order.
#         user (dict): Authenticated user details (extracted from JWT).
#         cursor (int): by default zero
#         limit (int): length of data
#     """
#     user_payload = UserPayload(**user)  # Convert dictionary to Pydantic model
#     return await order_by_eventID_controller(eventID, cursor, limit, user_payload)

# @router.get("/orders/user/{userID}", response_model=List[Dict[str, Any]])
# async def orders_by_userID(userID: int, cursor: int = 0, limit: int = 50, user: dict = Depends(get_current_user)):
#     """
#     API Endpoint: Admin-Focused Order Retrieval (Paginated by userID, Ordered by userID)

#     Args:
#         eventID (int): The unique identifier of the order.
#         user (dict): Authenticated user details (extracted from JWT).
#         cursor (int): by default zero
#         limit (int): length of data
#     """
#     user_payload = UserPayload(**user)  # Convert dictionary to Pydantic model
#     return await order_by_userID_controller(userID, cursor, limit, user_payload)

