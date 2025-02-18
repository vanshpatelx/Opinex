# src/models/trade_model.py

from pydantic import BaseModel
from typing import Literal

class Trade(BaseModel):
    OrderID: int
    EventID: int
    orderType: Literal["BUY", "SELL"]
    price: int
    qty: int
