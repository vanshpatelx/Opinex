# src/controllers/order_controller.py

"""
Order Controller

Handles order-related operations, including:
1. Create Order: `create_order`
2. Retrieve Order by OrderID: `order_by_orderID`
3. Retrieve Orders by EventID: `orders_by_eventID`
4. Retrieve Orders by UserID: `orders_by_userID`

Author: Vansh Patel (remotevansh@gmail.com)
Last Updated: February 19, 2025

Related Models:
- `OrderReq`
- `Order`
- `UserPayload`
"""

from fastapi import HTTPException, APIRouter, Depends
from config.cache.cache import Cache
from config.pubsub.pubsub import PubSub
from models.order_model import OrderReq, Order
from models.auth_model import UserPayload
from utils.id import generate_unique_id
from controllers.event_controller import get_event_by_id
from utils.logger import logger
from datetime import datetime

router = APIRouter()

@router.post("/order")
async def create_order(order: OrderReq, user: UserPayload):
    """
    Creates and processes a new order.

    Steps:
        1. Validate event status.
        2. Generate Order ID.
        3. Store in cache and publish to DB.
        4. Transform for trading engine.
        5. Publish to trading engine.

    Returns:
        dict: Order creation confirmation.

    Last Updated: February 19, 2025
    """
    # Validate event
    order_data = order.model_dump()
    user_data = user.model_dump()

    event_live = await get_event_by_id(order_data["EventID"])
    if not event_live:
        raise HTTPException(status_code=400, detail="Event is not live")

    id = str(generate_unique_id())

    order_obj = Order(
        id=id,
        EventID=order_data["EventID"],
        UserID=user_data["UserID"],
        orderType=order_data["orderType"],
        opt=order_data["opt"],
        qty=order_data["qty"],
        price=order_data["price"],
        status="PENDING",
        timestamp=int(datetime.now().timestamp())  # Converts to Unix timestamp in seconds
    )

    order_dict = order_obj.model_dump()

    # Convert IDs to string format
    order_dict["id"] = str(order_obj.id)
    order_dict["EventID"] = str(order_obj.EventID)
    order_dict["UserID"] = str(order_obj.UserID)
    logger.info(f"📦 Final Order: {order_dict}")

    # Store order in cache
    await Cache.set(f"order:{order_obj.id}", order_dict)
    logger.info(f"📦 Cache Set Order: {order_dict['id']}")

    # Send to DB server
    await PubSub.send(order_dict, 'Order_Exchange', 'Order.add')

    # Transform for trading engine
    order_type, order_price = transform_order(order_dict)

    trading_order = {
        "OrderID": order_dict["id"],
        "EventID": order_dict["EventID"],
        "orderType": order_type,
        "price": order_price,
        "qty": order_dict["qty"],
    }
    await PubSub.send(trading_order, 'Trade_Exchange', 'Trade.add')

    logger.info(f"✅ Order {order_obj.id} created successfully")
    return {"message": "Order created successfully"}

def transform_order(order_data):
    """
    Transforms order type and price based on order options.

    Last Updated: February 19, 2025
    """
    order_type, order_opt, order_price = order_data["orderType"], order_data["opt"], order_data["price"]

    if order_type == "BUY" and order_opt == "YES":
        return "BUY", order_price
    elif order_type == "SELL" and order_opt == "YES":
        return "BUY", 1000 - order_price
    elif order_type == "BUY" and order_opt == "NO":
        return "SELL", 1000 - order_price
    elif order_type == "SELL" and order_opt == "NO":
        return "SELL", order_price

    return order_type, order_price
