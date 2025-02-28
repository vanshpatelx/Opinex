# src/models/auth_model.py
from pydantic import BaseModel
from typing import Literal

class UserPayload(BaseModel):
    id: int
    type: Literal["ADMIN", "USER"]
