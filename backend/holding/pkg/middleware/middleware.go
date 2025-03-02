/***
 * JWT Middleware Package
 *
 * Provides JWT-based authentication and token validation middleware.
 *
 * Features:
 * - Extracts and validates JWT tokens from request headers.
 * - Verifies token signature and expiration.
 * - Attaches user details (ID, role) to the request context.
 *
 * Security Considerations:
 * - Store the secret key securely (e.g., environment variables).
 * - Implement proper token expiration policies.
 *
 * Author: Vansh Patel (remotevansh@gmail.com)
 * Last Updated: March 2, 2025
 ***/

package middleware

import (
	"errors"
	"log"
	"holding/pkg/config"
	"strings"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Claims struct for JWT token payload
type Claims struct {
	Id    string `json:"id"`
	Role  string `json:"type"`
	jwt.RegisteredClaims
}

func JWTMiddleware(c *fiber.Ctx) error {
    authHeader := c.Get("Authorization")

    if authHeader == "" {
        log.Println("🚨 Missing Authorization header")
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing Authorization header"})
    }

    // Extract token (format: "Bearer <token>")
    tokenString := strings.TrimPrefix(authHeader, "Bearer ")
    if tokenString == authHeader {
        log.Println("🚨 Invalid token format")
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token format"})
    }

    // Parse and validate token
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            log.Printf("🚨 Unexpected signing method: %v", token.Header["alg"])
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(config.AppConfig.JWTSecret), nil
    })

    // Handle invalid token
    if err != nil {
        log.Printf("🚨 JWT Parse Error: %v", err)
        if errors.Is(err, jwt.ErrSignatureInvalid) {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token signature"})
        } else if errors.Is(err, jwt.ErrTokenExpired) {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token expired"})
        }
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
    }

    // Extract claims
    claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid {
        log.Println("🚨 Invalid token claims")
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
    }

    // Attach user info to request context
    c.Locals("Id", claims.Id)
    c.Locals("Role", claims.Role)

    return c.Next()
}