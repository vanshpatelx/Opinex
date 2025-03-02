/***
JWT Middleware Package

Handles JWT-based authentication and token validation.

Features:
1. **JWTMiddleware**:
   - Extracts and validates JWT tokens from requests.
   - Verifies token signature and expiration.
   - Attaches user information (userID, role, email) to request context.

Security Considerations:
- The secret key should be stored securely (e.g., environment variables).
- The token should use proper expiration policies.

Author: Vansh Patel (remotevansh@gmail.com)
Last Updated: March 2, 2025
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

/***
JWTMiddleware validates the JWT token from the Authorization header.

Steps:
1. Check for the presence of the Authorization header.
2. Extract the token from the "Bearer <token>" format.
3. Parse and validate the JWT token signature.
4. Check token expiration and validity.
5. Extract claims (id, role, email) and store them in request context.
6. Proceed with the request if validation is successful.

Returns:
- 401 Unauthorized if token is missing, invalid, or expired.
- Calls next handler if the token is valid.

Last Updated: March 2, 2025
***/

// func JWTMiddleware(c *fiber.Ctx) error {
// 	authHeader := c.Get("Authorization")

// 	// Check if Authorization header is present
// 	if authHeader == "" {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing Authorization header"})
// 	}

// 	// Extract token (format: "Bearer <token>")
// 	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
// 	if tokenString == authHeader {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token format"})
// 	}

// 	// Parse and validate token
// 	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 		}
// 		return []byte(config.AppConfig.JWTSecret), nil
// 	})

// 	// Handle JWT validation errors
// 	if err != nil {
// 		var errMsg string
// 		if errors.Is(err, jwt.ErrTokenMalformed) {
// 			errMsg = "Malformed token"
// 		} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
// 			errMsg = "Invalid token signature"
// 		} else if errors.Is(err, jwt.ErrTokenExpired) {
// 			errMsg = "Token expired"
// 		} else if errors.Is(err, jwt.ErrTokenNotValidYet) {
// 			errMsg = "Token not valid yet"
// 		} else {
// 			errMsg = "Invalid token"
// 		}
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": errMsg})
// 	}

// 	// Extract claims
// 	claims, ok := token.Claims.(*Claims)
// 	if !ok || !token.Valid {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
// 	}

// 	// Attach user info to request context
// 	c.Locals("Id", claims.Id)
// 	c.Locals("Role", claims.Role)
// 	c.Locals("Email", claims.Email)

// 	return c.Next()
// }

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