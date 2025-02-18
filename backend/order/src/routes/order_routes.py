# src/routes/order_routes.py
from fastapi import APIRouter, Depends
from models.order_model import OrderReq
from middleware.jwt import get_current_user
from controllers.order_controller import create_order
from models.auth_model import UserPayload

router = APIRouter()

@router.post("/order", response_model=dict)
async def create_new_order(order: OrderReq, user: dict = Depends(get_current_user)):
    """Create a new order"""
    user_payload = UserPayload(**user)  # Convert dictionary to Pydantic model
    return await create_order(order, user_payload)

# @router.get("/order/{orderID}", response_model=dict)
# async def fetch_order_by_orderID(orderID: int, user: dict = Depends(get_current_user)):
#     """Fetch order by OrderID."""
#     return await get_order_by_orderID(orderID)

# @router.get("/order/event/{eventID}", response_model=dict)
# async def fetch_order_by_eventID(eventID: int, user: dict = Depends(get_current_user)):
#     """Fetch order by EventID."""
#     return await get_order_by_eventID(eventID)

# @router.get("/order/user/{userID}", response_model=dict)
# async def fetch_order_by_userID(userID: int, user: dict = Depends(get_current_user)):
#     """Fetch orders by UserID."""
#     return await get_order_by_userID(userID)


