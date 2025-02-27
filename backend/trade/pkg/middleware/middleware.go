/***
JWT Middleware Package

Handles JWT-based authentication and token generation.

Features:
1. **JWTMiddleware**:
   - Extracts and validates JWT tokens from requests.
   - Verifies token signature and expiration.
   - Attaches user information (userID, role) to request context.

2. **GenerateJWT**:
   - Generates JWT tokens for authentication.
   - Includes userID, role, and expiration claims.

Security Considerations:
- The secret key should be stored securely (e.g., environment variables).
- The token should use proper expiration policies.

Author: Vansh Patel (remotevansh@gmail.com)
Last Updated: February 26, 2025
***/

package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// JWT Secret Key (Store in env or config)
var jwtSecret = []byte("your-secret-key")

// Claims struct
type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

/***
 JWTMiddleware validates the JWT token from the Authorization header.

Steps:
1. Check for the presence of the Authorization header.
2. Extract the token from the "Bearer <token>" format.
3. Parse and validate the JWT token signature.
4. Check token expiration and validity.
5. Extract claims (userID, role) and store them in request context.
6. Proceed with the request if validation is successful.

Returns:
- 401 Unauthorized if token is missing, invalid, or expired.
- Calls next handler if the token is valid.

Last Updated: February 26, 2025
***/

func JWTMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	// Check if Authorization header is present
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing Authorization header"})
	}

	// Extract token (format: "Bearer <token>")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token format"})
	}

	// Parse and validate token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	// Handle invalid or expired token
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	// Extract claims
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	// Attach user info to request context
	c.Locals("userID", claims.UserID)
	c.Locals("role", claims.Role)

	// Continue request
	return c.Next()
}

// GenerateJWT creates a JWT token for testing
func GenerateJWT(userID, role string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Token valid for 24 hours
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
