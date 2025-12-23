package auth

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// GetUserIDFromToken extrae el UUID del usuario desde el token JWT en el contexto
func GetUserIDFromToken(c *fiber.Ctx) (uuid.UUID, error) {
	// Fiber guarda el token en c.Locals("user") gracias al middleware jwtware
	userToken, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return uuid.Nil, errors.New("token JWT no encontrado o inválido")
	}

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, errors.New("claims del token inválidos")
	}

	// Asumiendo que guardaste el ID como "user_id" o "id" en el token al hacer Login
	// Revisa tu auth.service.go (Login) para ver qué clave usaste.
	// Usualmente es "user_id" o "sub". Aquí pruebo con "user_id".
	idStr, ok := claims["user_id"].(string)
	if !ok {
		// Si falló, intentamos con "id" por si acaso
		idStr, ok = claims["id"].(string)
		if !ok {
			return uuid.Nil, errors.New("ID de usuario no encontrado en el token")
		}
	}

	uid, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, errors.New("formato de UUID inválido en el token")
	}

	return uid, nil
}
