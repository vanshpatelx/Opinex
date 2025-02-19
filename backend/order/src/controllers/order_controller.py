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

from fastapi import HTTPException
from config.cache.cache import Cache
from config.pubsub.pubsub import PubSub
from models.order_model import OrderReq, Order
from models.auth_model import UserPayload
from utils.id import generate_unique_id
from controllers.event_controller import get_event_by_id
from utils.logger import logger
from datetime import datetime
from config.db.db import Database
from config.db.query import Query

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


async def order_by_orderID_controller(orderID: int, user: UserPayload):
    """
    Retrieves an order by orderID with optimized authorization and caching.

    Optimized Steps:
        1. First, check if `UserID` exists in Redis (to validate access quickly).
        2. If authorized, fetch full order details from Redis.
        3. If order not found in cache, fetch from PostgreSQL.
        4. Store in Redis for future requests.
    
    Returns:
        Dict[str, Any]: Order details.

    Last Updated: February 19, 2025
    """
    logger.info(f"📥 Request: Fetching Order {orderID} for User {user.UserID} ({user.user_type})")
    cache_key = f"order:{orderID}"

    cached_userID = await Cache.get(f"order:{orderID}:UserID")  # Fetch only UserID

    if cached_userID:
        # If user is not admin and doesn't own the order, block access early
        if user.user_type != "ADMIN" and cached_userID != str(user.UserID):
            logger.warning(f"🚫 Unauthorized access: User {user.UserID} attempted to access Order {orderID}.")
            raise HTTPException(status_code=403, detail="Unauthorized access to this order.")

        # Fetch Full Order Data from Redis
        cached_order = await Cache.get(cache_key)
        if cached_order:
            logger.info(f"✅ Cache hit: Order {orderID} retrieved from Redis.")
            return cached_order

    # Fetch from PostgreSQL (if not in Redis)
    order = await Database.fetch_one(Query.get_order_by_orderID, orderID)

    if not order:
        logger.warning(f"❌ Order {orderID} not found in database.")
        raise HTTPException(status_code=404, detail="Order not found")

    order_data = dict(order)

    # Authorization Check After Fetching from DB
    if user.user_type != "ADMIN" and order_data["userid"] != user.UserID:
        logger.warning(f"🚫 Unauthorized access: User {user.UserID} attempted to access Order {orderID}.")
        raise HTTPException(status_code=403, detail="Unauthorized access to this order.")

    # Store UserID Separately for Fast Access Validation
    await Cache.set(f"order:{orderID}:UserID", str(order_data["userid"]))

    order_data["id"] = str(order_data["id"])
    order_data["EventID"] = str(order_data.pop("eventid"))
    order_data["UserID"] = str(order_data.pop("userid"))

    # Store Full Order in Cache
    await Cache.set(cache_key, order_data) 

    logger.info(f"✅ Order {orderID} stored in Redis for future requests.")
    return order_data



async def order_by_eventID_controller(eventID: int, cursor: int, limit: int, user: UserPayload):
    """
    Fetches orders for a given event ID based on user type.

    Steps:
        1. Log request details.
        2. Check user type (ADMIN or USER).
        3. Route request to respective function:
            - Admin: Fetch all orders for the event.
            - User: Fetch only user-specific orders.

    Returns:
        list[dict]: A list of orders for the event.

    Last Updated: February 19, 2025
    """
    logger.info(f"📥 Request: Fetching eventID orders {eventID} for User {user.UserID} ({user.user_type})")

    if user.user_type == "ADMIN":
        return await order_by_eventID_admin(eventID, cursor, limit)
    elif user.user_type == "USER":
        return await order_by_eventID_user(eventID, cursor, limit, user)



async def order_by_eventID_admin(eventID: int, cursor: int, limit: int):
    """
    Fetches all orders for a given event ID (Admin access).

    Steps:
        1. Generate cache key.
        2. Check cache for existing orders.
        3. If found, return cached data.
        4. Query database for orders.
        5. If no orders found, return 404.
        6. Convert order fields to proper format.
        7. Store fetched orders in cache.
        8. Return processed order list.

    Returns:
        list[dict]: Orders for the event.

    Last Updated: February 19, 2025
    """
    cache_key = f"orders:Event{eventID}:{cursor}:{limit}"

    cached_orders = await Cache.get(cache_key)
    if cached_orders:
        logger.info(f"✅ Cache hit: Orders for Event {eventID} retrieved from Redis.")
        return cached_orders
    
    orders = await Database.fetch_all(Query.get_order_by_eventID_pagination_admin, eventID, limit, cursor)

    if not orders:
        logger.warning(f"❌ Orders for Event {eventID} not found in database.")
        raise HTTPException(status_code=404, detail="Orders not found")

    orders_data = [dict(order) for order in orders]

    for order in orders_data:
        order["id"] = str(order["id"])
        order["EventID"] = str(order.pop("eventid"))
        order["UserID"] = str(order.pop("userid"))

    await Cache.set(cache_key, orders_data, expire=300) 

    logger.info(f"✅ Orders of eventID {eventID} cursor:{cursor} limit:{limit} stored in Redis for future requests.")
    return orders_data

async def order_by_eventID_user(eventID: int, cursor: int, limit: int, user: UserPayload):
    """
    Fetches user-specific orders for a given event ID.

    Steps:
        1. Generate cache key specific to user.
        2. Check cache for existing orders.
        3. If found, return cached data.
        4. Query database for user orders.
        5. If no orders found, return 404.
        6. Convert order fields to proper format.
        7. Store fetched orders in cache.
        8. Return processed order list.

    Returns:
        list[dict]: User-specific orders for the event.

    Last Updated: February 19, 2025
    """

    cache_key = f"orders:Event{eventID}:{user.UserID}:{cursor}:{limit}"

    cached_orders = await Cache.get(cache_key)
    if cached_orders:
        logger.info(f"✅ Cache hit: Orders for Event {eventID} retrieved from Redis.")
        return cached_orders

    orders = await Database.fetch_all(Query.get_order_by_eventID_pagination_user, user.UserID, eventID, limit, cursor)
    if not orders:
        logger.warning(f"❌ Orders for Event {eventID} not found in database.")
        raise HTTPException(status_code=404, detail="Orders not found")

    orders_data = [dict(order) for order in orders]

    for order in orders_data:
        order["id"] = str(order["id"])
        order["EventID"] = str(order.pop("eventid"))
        order["UserID"] = str(order.pop("userid"))

    await Cache.set(cache_key, orders_data, expire=300) 

    orders_data = [dict(order) for order in orders]

    logger.info(f"✅ Orders of EventID {eventID}, cursor: {cursor}, limit: {limit} stored in Redis for future requests.")
    return orders_data



async def order_by_userID_controller(userID: int, cursor: int, limit: int, user: UserPayload):
    """
    Fetches orders for a given user ID based on user type.

    Steps:
        1. Log request details.
        2. Check user type (ADMIN or USER).
        3. Route request to respective function:
            - Admin: Fetch all orders for the user.
            - User: Fetch only their own orders.

    Returns:
        list[dict]: A list of orders for the user.

    Last Updated: February 19, 2025
    """
    logger.info(f"📥 Request: Fetching users orders {userID} for User {user.UserID} ({user.user_type})")

    if user.user_type == "ADMIN":
        return await order_by_userID_admin(userID, cursor, limit)
    elif user.user_type == "USER":
        return await order_by_userID_user(user.UserID, cursor, limit, user)


async def order_by_userID_admin(userID: int, cursor: int, limit: int):
    """
    Fetches all orders for a given user ID (Admin access).

    Steps:
        1. Generate cache key.
        2. Check cache for existing orders.
        3. If found, return cached data.
        4. Query database for orders.
        5. If no orders found, return 404.
        6. Convert order fields to proper format.
        7. Store fetched orders in cache.
        8. Return processed order list.

    Returns:
        list[dict]: Orders for the user.

    Last Updated: February 19, 2025
    """

    cache_key = f"orders:User{userID}:{cursor}:{limit}"

    cached_orders = await Cache.get(cache_key)
    if cached_orders:
        logger.info(f"✅ Cache hit: Orders for Event {userID} retrieved from Redis.")
        return cached_orders
    
    orders = await Database.fetch_all(Query.get_order_by_userID_pagination, userID, limit, cursor)

    if not orders:
        logger.warning(f"❌ Orders for User {userID} not found in database.")
        raise HTTPException(status_code=404, detail="Orders not found")

    orders_data = [dict(order) for order in orders]

    for order in orders_data:
        order["id"] = str(order["id"])
        order["EventID"] = str(order.pop("eventid"))
        order["UserID"] = str(order.pop("userid"))

    await Cache.set(cache_key, orders_data, expire=300) 

    logger.info(f"✅ Orders of eventID {userID} cursor:{cursor} limit:{limit} stored in Redis for future requests.")
    return orders_data

async def order_by_userID_user(userID: int, cursor: int, limit: int, user: UserPayload):
    """
    Fetches user-specific orders for a given user ID.

    Steps:
        1. Generate cache key specific to user.
        2. Check cache for existing orders.
        3. If found, return cached data.
        4. Query database for user orders.
        5. If no orders found, return 404.
        6. Convert order fields to proper format.
        7. Store fetched orders in cache.
        8. Return processed order list.

    Returns:
        list[dict]: User-specific orders.

    Last Updated: February 19, 2025
    """

    cache_key = f"orders:User{userID}:{cursor}:{limit}"

    cached_orders = await Cache.get(cache_key)
    if cached_orders:
        logger.info(f"✅ Cache hit: Orders for Event {userID} retrieved from Redis.")
        return cached_orders

    orders = await Database.fetch_all(Query.get_order_by_userID_pagination, userID, limit, cursor)
    if not orders:
        logger.warning(f"❌ Orders for Event {userID} not found in database.")
        raise HTTPException(status_code=404, detail="Orders not found")

    orders_data = [dict(order) for order in orders]

    for order in orders_data:
        order["id"] = str(order["id"])
        order["EventID"] = str(order.pop("eventid"))
        order["UserID"] = str(order.pop("userid"))

    await Cache.set(cache_key, orders_data, expire=300) 

    orders_data = [dict(order) for order in orders]

    logger.info(f"✅ Orders of UserID {userID}, cursor: {cursor}, limit: {limit} stored in Redis for future requests.")
    return orders_data