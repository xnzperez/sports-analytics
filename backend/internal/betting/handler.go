package betting

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Handler struct {
	service *Service
}

func NewHandler(db *gorm.DB) *Handler {
	repo := NewRepository(db)
	service := NewService(repo)
	return &Handler{service: service}
}

// PlaceBet crea una nueva apuesta
func (h *Handler) PlaceBet(c *fiber.Ctx) error {
	// 1. Obtener ID del usuario
	userIDStr := c.Locals("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)

	// 2. Parsear el Body
	var req PlaceBetRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	// 3. Validaciones simples
	if req.StakeUnits <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "El stake debe ser mayor a 0"})
	}

	// 4. Llamar al servicio
	bet, err := h.service.PlaceBet(userID, req)
	if err != nil {
		if err.Error() == "saldo insuficiente para realizar esta apuesta" {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Error interno al procesar apuesta"})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Apuesta realizada con éxito",
		"bet":     bet,
	})
}

// ResolveBetRequest define qué esperamos recibir en el JSON
type ResolveBetRequest struct {
	Outcome string `json:"outcome"` // "WON" o "LOST"
}

// ResolveBetHandler gestiona la resolución de la apuesta
func (h *Handler) ResolveBetHandler(c *fiber.Ctx) error {
	betID := c.Params("id")

	var req ResolveBetRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "El formato del JSON es incorrecto",
		})
	}

	if req.Outcome != "WON" && req.Outcome != "LOST" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "El resultado (outcome) debe ser 'WON' o 'LOST'",
		})
	}

	// CORRECCIÓN: Llamamos al service, no al repo directamente
	err := h.service.ResolveBet(betID, req.Outcome)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "Apuesta resuelta correctamente",
		"new_status": req.Outcome,
	})
}

func (h *Handler) GetBetsHandler(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)

	// Leemos los Query Params de la URL
	// Ej: /api/bets?page=2&limit=5&status=WON
	page := c.QueryInt("page", 1)    // Default: 1
	limit := c.QueryInt("limit", 10) // Default: 10
	status := c.Query("status")      // Opcional
	sportKey := c.Query("sport_key") // Opcional

	filters := BetFilters{
		UserID:   userID,
		Page:     page,
		Limit:    limit,
		Status:   status,
		SportKey: sportKey,
	}

	response, err := h.service.GetBets(filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al obtener las apuestas",
		})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *Handler) GetStatsHandler(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)

	stats, err := h.service.GetUserStats(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error calculando estadísticas",
		})
	}

	return c.Status(fiber.StatusOK).JSON(stats)
}

// GetTransactionsHandler maneja la petición del historial
func (h *Handler) GetTransactionsHandler(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)

	// Leer params de la URL (ej: ?page=1&limit=10)
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)

	response, err := h.service.GetTransactions(userID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "No se pudo obtener el historial",
		})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
