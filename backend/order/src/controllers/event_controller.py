# src/controllers/event_controller.py

"""
Event Controller

Handles event-related operations, including:
1. Checking if an event is live.

Author: Vansh Patel (remotevansh@gmail.com)
Last Updated: February 19, 2025
"""

from config.cache.cache import Cache

async def get_event_by_id(eventID: int) -> bool:
    """
    Check if an event is live.

    - Retrieves event status from the cache.

    Args:
        eventID (int): The unique event identifier.
    
    Returns:
        bool: True if the event is live, False otherwise.

    Last Updated: February 19, 2025
    """
    key = f"Event:{eventID}"
    status = await Cache.get(key)
    return bool(status)
