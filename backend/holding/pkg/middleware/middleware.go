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
	"fmt"
	"holding/pkg/config"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Claims struct for JWT token payload
type Claims struct {
	Id    string `json:"id"`
	Email string `json:"email"`
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
	fmt.Println("Auth Header:", authHeader) // Print received header

	if authHeader == "" {
		fmt.Println("🚨 Missing Authorization header")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing Authorization header"})
	}

	// Extract token (format: "Bearer <token>")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		fmt.Println("🚨 Invalid token format")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token format"})
	}

	// Parse and validate token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("🚨 Unexpected signing method:", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		fmt.Println("JWT Secret Key:", config.AppConfig.JWTSecret) // Print secret key
		return []byte(config.AppConfig.JWTSecret), nil
	})

	// Handle invalid token
	if err != nil {
		fmt.Println("🚨 JWT Parse Error:", err)
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
		fmt.Println("🚨 Invalid token claims")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	// Print extracted claims
	fmt.Println("✅ Token Valid. Extracted Claims:")
	fmt.Println("   - User ID:", claims.Id)
	fmt.Println("   - Email:", claims.Email)
	fmt.Println("   - Role:", claims.Role)
	fmt.Println("   - Expires At:", claims.ExpiresAt)

	// Attach user info to request context
	c.Locals("Id", claims.Id)
	c.Locals("Role", claims.Role)

	return c.Next()
}
