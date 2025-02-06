import jwt, { Secret, SignOptions } from "jsonwebtoken";

const JWT_SECRET: Secret = process.env.JWT_SECRET || "your_secret_key";

/**
 * Generates a JWT token for a user
 * @param payload - The data to be included in the token (e.g., user ID, email)
 * @param expiresIn - Token expiration time (default: 7 days)
 * @returns A signed JWT token
 */

export const generateToken = (payload: object, expiresIn: number = 86400): string => {
    const options: SignOptions = { expiresIn }; // Explicitly define SignOptions
    return jwt.sign(payload, JWT_SECRET, options);
};
