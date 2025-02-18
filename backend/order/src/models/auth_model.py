# src/models/auth_model.py
from pydantic import BaseModel
from typing import Literal

class UserPayload(BaseModel):
    UserID: int
    user_type: Literal["ADMIN", "USER"]
