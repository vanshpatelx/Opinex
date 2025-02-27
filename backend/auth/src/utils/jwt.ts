// src/utils/jwt.ts
/**
    JWT Token Utility

    Provides helper functions for generating JWT tokens.

    Features:
    - `generateToken(payload, expiresIn)`: Signs a JWT token with a given payload and expiration.

    Dependencies:
    - jsonwebtoken (JWT signing and verification)
    - Configuration for secret management

    Author: Vansh Patel (remotevansh@gmail.com)  
    Date: February 2, 2025  
 */


import jwt, { Secret, SignOptions } from "jsonwebtoken";
import { config } from "../config/config";

const JWT_SECRET: Secret = config.JWT_SECRET.secret;

export const generateToken = (payload: object, expiresIn: number = 86400): string => {
    const options: SignOptions = { expiresIn }; // Explicitly define SignOptions
    return jwt.sign(payload, JWT_SECRET, options);
};
