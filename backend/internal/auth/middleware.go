package auth

import (
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Protected es el middleware que bloquea accesos sin token válido
func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Obtener el header Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{"error": "No autorizado: Falta token"})
		}

		// 2. El formato debe ser "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(401).JSON(fiber.Map{"error": "Formato de token inválido"})
		}
		tokenString := parts[1]

		// 3. Parsear y Validar el token
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			// Validar el algoritmo de firma (debe ser HMAC)
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("método de firma inesperado: %v", t.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{"error": "Token inválido o expirado"})
		}

		// 4. Extraer Claims (Datos del usuario) e inyectarlos en el Contexto
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			// Guardamos el user_id en c.Locals para usarlo en los controladores
			c.Locals("user_id", claims["user_id"])
		}

		// 5. Dejar pasar a la siguiente función
		return c.Next()
	}
}
