# src/controllers/order_controller.py

from fastapi import HTTPException
from config.cache.cache import Cache
from config.pubsub.pubsub import PubSub
from models.order_model import OrderReq, Order
from models.auth_model import UserPayload
from utils.id import generate_unique_id
from controllers.event_controller import get_event_by_id
from utils.logger import logger
from datetime import datetime

async def create_order(order: OrderReq, user: UserPayload):
    """Create a new order."""

    # Req_data validations
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

    # Converting long to string for string convert issue
    order_dict["id"] = str(order_obj.id)
    order_dict["EventID"] = str(order_obj.EventID)
    order_dict["UserID"] = str(order_obj.UserID)
    logger.info(f"📦 Final Order Dict: {order_dict}")

    await Cache.set(f"order:{order_obj.id}", order_dict)
    
    # send to Dbserver
    await PubSub.send(order_dict, 'Order_Exchange', 'Order.add')

    # tranform for trade server
    order_type, order_price = transform_order(order_dict)

    trading_order = {
        "OrderID": order_dict["id"],
        "EventID": order_dict["EventID"],
        "orderType": order_type,
        "price": order_price,
        "qty": order_dict["qty"],
    }
    await PubSub.send(trading_order, 'Trade_Exchange', 'Trade.add')

    logger.info(f"Order {order_obj.id} created successfully")
    return {"message": "Order created successfully"}

def transform_order(order_data):
    """Transforms order type and price based on order options."""
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



    # check it's live event or not
    # gen ID
    # add to redis
    # send to dbserver
    # Convert to base form trading engine
    # send to trading engine


# Create order
# - jwt validations
# - zod validations
# - Orders send to DbServers by pubsub
# - Convert to base form for trade engine
# - Orders send to trading engine  by pubsub

