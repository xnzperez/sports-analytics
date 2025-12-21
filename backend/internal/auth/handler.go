package auth

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Handler struct {
	service *Service
}

// NewHandler inicializa todo el módulo de Auth (Repo + Service)
func NewHandler(db *gorm.DB) *Handler {
	repo := NewRepository(db)
	service := NewService(repo)
	return &Handler{service: service}
}

func (h *Handler) Register(c *fiber.Ctx) error {
	var req RegisterRequest

	// 1. Parsear el Body (JSON)
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	// --- CORRECCIÓN ---
	// Si el usuario no envía Username, usamos la parte antes del @ del email o el email completo
	if req.Username == "" {
		req.Username = req.Email
	}
	// ------------------

	// 2. Validaciones básicas (QUITAMOS req.Username de la condición)
	if req.Email == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "El email y la contraseña son obligatorios"})
	}

	// 3. Llamar al servicio
	if err := h.service.RegisterUser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// 4. Éxito
	return c.Status(201).JSON(fiber.Map{
		"message": "Usuario registrado exitosamente",
	})
}

// Login maneja la petición de inicio de sesión
func (h *Handler) Login(c *fiber.Ctx) error {
	var req LoginRequest

	// 1. Parsear JSON
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	// 2. Llamar al servicio
	token, err := h.service.LoginUser(req)
	if err != nil {
		// Retornamos 401 Unauthorized si falla
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	// 3. Responder con el token
	return c.JSON(fiber.Map{
		"message": "Login exitoso",
		"token":   token,
	})
}

// En NewHandler, agregaremos la ruta GET protegida más adelante en main.go
// Por ahora, solo añade la función GetMe al struct Handler:

func (h *Handler) GetMe(c *fiber.Ctx) error {
	// 1. Recuperar el ID del usuario desde el contexto (puesto por el middleware)
	userID := c.Locals("user_id").(string)

	// 2. Consultar al servicio
	user, err := h.service.GetUserProfile(userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Usuario no encontrado"})
	}

	// 3. Responder con los datos (incluyendo Bankroll)
	return c.JSON(fiber.Map{
		"user": user,
	})
}
