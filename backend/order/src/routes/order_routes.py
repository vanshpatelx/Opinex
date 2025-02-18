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
from middleware.jwt import get_current_user  # Middleware to verify JWT token
from controllers.order_controller import create_order  # Import order creation controller
from models.auth_model import UserPayload  # User authentication model

# Initialize FastAPI router
router = APIRouter()

@router.post("/order", response_model=dict)
async def create_new_order(order: OrderReq, user: dict = Depends(get_current_user)):
    """
    API Endpoint: Create a new order.

    - Validates JWT token and extracts user details.
    - Converts user dictionary to Pydantic model (`UserPayload`).
    - Calls `create_order()` to process and publish the order.

    Args:
        order (OrderReq): Order request payload.
        user (dict): Authenticated user details (extracted from JWT).

    Returns:
        dict: Confirmation message.
    
    Last Updated: February 19, 2025
    """
    user_payload = UserPayload(**user)  # Convert dictionary to Pydantic model
    return await create_order(order, user_payload)