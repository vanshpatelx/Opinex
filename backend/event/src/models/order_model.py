# src/models/order_model.py

from pydantic import BaseModel
from typing import Literal
from datetime import datetime

class Event(BaseModel):
    id: int
    name: str
    details: str
    status: Literal["RUNNING", "CLOSED", "SETTLEMENT", "PROCESSED"]
    settlement_time: str
    created_at: str

class EventReq(BaseModel):
    name: str
    details: str
    settlement_time: str


class EventResponse(BaseModel):
    id: int
    name: str
    details: str
    status: Literal["RUNNING", "CLOSED", "SETTLEMENT", "PROCESSED"]
    settlement_time: str
