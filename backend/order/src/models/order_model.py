# src/models/order_model.py

from pydantic import BaseModel
from typing import Literal

class Order(BaseModel):
    id: int
    EventID: int
    orderType: Literal["BUY", "SELL"]
    opt: Literal["YES", "NO"]
    price: int
    qty: int
    UserID: int
    status: str
    timestamp: int

class OrderReq(BaseModel):
    EventID: int
    orderType: Literal["BUY", "SELL"]
    opt: Literal["YES", "NO"]
    price: int
    qty: int
