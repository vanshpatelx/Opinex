# src/models/order_model.py

from pydantic import BaseModel
from typing import Literal
from datetime import datetime

class Event(BaseModel):
    id: int
    name: str
    details: str
    status: Literal["RUNNING", "CLOSED", "SETTLEMENT", "PROCESSED"]
    settlement_time: datetime
    updated_at: datetime
    created_at: datetime

class EventReq(BaseModel):
    name: str
    details: str
    settlement_time: datetime
