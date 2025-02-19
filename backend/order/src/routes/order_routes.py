# src/routes/order_routes.py

"""
Order Routes

Defines API routes for order-related operations, including:
1. Create a new order.
2. Fetch order details by OrderID.
3. Fetch all orders for a specific EventID.
4. Fetch all orders placed by a specific UserID.

Author: Vansh Patel (remotevansh@gmail.com)
Last Updated: February 19, 2025

Related Controllers:
- `create_order`
- `order_by_orderID`
- `orders_by_eventID`
- `orders_by_userID`
"""

from fastapi import APIRouter, Depends
from models.order_model import OrderReq
from middleware.jwt import get_current_user 
from controllers.order_controller import create_order, order_by_orderID_controller, order_by_eventID_controller, order_by_userID_controller
from models.auth_model import UserPayload
from typing import Dict, Any, List

router = APIRouter()

@router.post("/order", response_model=dict)
async def create_new_order(order: OrderReq, user: dict = Depends(get_current_user)):
    """
    API Endpoint: Create a new order.

    Args:
        order (OrderReq): Order request payload.
        user (dict): Authenticated user details (extracted from JWT).

    Returns:
        dict: Confirmation message.
    
    Last Updated: February 19, 2025
    """
    user_payload = UserPayload(**user)  # Convert dictionary to Pydantic model
    return await create_order(order, user_payload)


@router.get("/order/{orderID}", response_model=Dict[str, Any])
async def order_by_orderID(orderID: int, user: dict = Depends(get_current_user)):
    """
    API Endpoint: Retrieve an order by ID.

    Args:
        orderID (int): The unique identifier of the order.
        user (dict): Authenticated user details (extracted from JWT).

    Returns:
        Dict[str, Any]: Order details.
    
    Last Updated: February 19, 2025
    """
    user_payload = UserPayload(**user)  # Convert dictionary to Pydantic model
    return await order_by_orderID_controller(orderID, user_payload)


@router.get("/orders/event/{eventID}", response_model=List[Dict[str, Any]])
async def orders_by_eventID(eventID: int, cursor: int = 0, limit: int = 50, user: dict = Depends(get_current_user)):
    """
    API Endpoint: Admin-Focused Order Retrieval (Paginated by EventID, Ordered by EventID)

    Args:
        eventID (int): The unique identifier of the order.
        user (dict): Authenticated user details (extracted from JWT).
        cursor (int): by default zero
        limit (int): length of data
    """
    user_payload = UserPayload(**user)  # Convert dictionary to Pydantic model
    return await order_by_eventID_controller(eventID, cursor, limit, user_payload)

@router.get("/orders/user/{userID}", response_model=List[Dict[str, Any]])
async def orders_by_userID(userID: int, cursor: int = 0, limit: int = 50, user: dict = Depends(get_current_user)):
    """
    API Endpoint: Admin-Focused Order Retrieval (Paginated by userID, Ordered by userID)

    Args:
        eventID (int): The unique identifier of the order.
        user (dict): Authenticated user details (extracted from JWT).
        cursor (int): by default zero
        limit (int): length of data
    """
    user_payload = UserPayload(**user)  # Convert dictionary to Pydantic model
    return await order_by_userID_controller(userID, cursor, limit, user_payload)

