package market

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Handler struct {
	service *Service
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{service: NewService(db)}
}

// SyncMarketsHandler dispara la actualización con RapidAPI
func (h *Handler) SyncMarketsHandler(c *fiber.Ctx) error {
	count, err := h.service.SyncEsports()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{
		"message":         "Sincronización completada",
		"matches_updated": count,
	})
}

// ListMarketsHandler devuelve los partidos desde TU base de datos
func (h *Handler) ListMarketsHandler(c *fiber.Ctx) error {
	sport := c.Query("sport") // ?sport=lol
	matches, err := h.service.GetMatches(sport)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error leyendo base de datos"})
	}
	return c.JSON(fiber.Map{"data": matches})
}
