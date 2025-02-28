# src/middleware/jwt.py

"""
JWT Middleware

Handles authentication using JWT tokens.

Features:
1. Extracts and validates JWT token from Authorization header.
2. Decodes token payload and verifies signature.
3. Handles token expiration and invalid token cases.

Author: Vansh Patel (remotevansh@gmail.com)
Last Updated: February 19, 2025
"""

import jwt
from jwt import ExpiredSignatureError, InvalidTokenError
from fastapi import HTTPException, Security
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials
from src.config.config import config  # Import JWT secret

# Define security scheme (Bearer token)
security = HTTPBearer()

def get_current_user(credentials: HTTPAuthorizationCredentials = Security(security)):
    """
    Extract and validate the JWT token from the Authorization header.

    - Decodes the token using the secret key.
    - Checks if the token is expired or invalid.
    
    Args:
        credentials (HTTPAuthorizationCredentials): Extracted bearer token.
    
    Returns:
        dict: Decoded token payload (user details).
    
    Raises:
        HTTPException: If the token is expired or invalid.

    Last Updated: February 19, 2025
    """
    token = credentials.credentials
    try:
        payload = jwt.decode(token, config.JWT_SECRET, algorithms=["HS256"])
        return payload  # Return decoded JWT payload
    except ExpiredSignatureError:
        raise HTTPException(status_code=401, detail="Token has expired")
    except InvalidTokenError:
        raise HTTPException(status_code=401, detail="Invalid token")
