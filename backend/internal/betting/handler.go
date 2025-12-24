package betting

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/xnzperez/sports-analytics-backend/internal/ai"

	// "auth" lo quitamos porque ya no lo necesitamos aquí
	"gorm.io/gorm"
)

type Handler struct {
	service   *Service
	aiService *ai.Service
}

func NewHandler(db *gorm.DB) *Handler {
	repo := NewRepository(db)
	service := NewService(repo)
	aiService := ai.NewService()
	return &Handler{
		service:   service,
		aiService: aiService,
	}
}

// NOTA: Borramos 'type PlaceBetRequest struct...' de aquí
// porque ya debe estar en service.go (según el paso anterior).

// PlaceBet crea una nueva apuesta
// @Router /api/bets [post]
func (h *Handler) PlaceBet(c *fiber.Ctx) error {
	// 1. Obtener ID del usuario (Forma correcta: string desde Locals)
	userIDStr := c.Locals("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)

	// 2. Parsear el Body (Usando el struct definido en service.go)
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

// ResolveBetHandler define el resultado de una apuesta.
// @Router /api/bets/{id}/resolve [patch]
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

// GetBetsHandler obtiene el historial de apuestas
// @Router /api/bets [get]
func (h *Handler) GetBetsHandler(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	status := c.Query("status")
	sportKey := c.Query("sport_key")

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

// GetStatsHandler calcula el rendimiento (CORREGIDO)
// @Router /api/stats [get]
func (h *Handler) GetStatsHandler(c *fiber.Ctx) error {
	val := c.Locals("user_id")
	if val == nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	userIDStr := val.(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid User ID"})
	}

	// CAPTURAMOS EL FILTRO: Ejemplo /api/stats?sport=lol
	sportFilter := c.Query("sport")

	// 1. Obtenemos estadísticas filtradas (Asegúrate que tu service reciba este string)
	// Si tu service aún no lo recibe, puedes pasarle solo el userID por ahora
	// pero aquí ya preparamos el Handler para el futuro.
	stats, err := h.service.GetUserDashboardStats(userID, sportFilter)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error calculando estadísticas"})
	}

	// 2. Determinar el "Deporte Top"
	topSport := "General"
	if sportFilter != "" {
		topSport = sportFilter
	} else if len(stats.SportPerformance) > 0 {
		topSport = stats.SportPerformance[0].SportKey
	}

	// 3. Generar el Tip de IA
	tip := h.aiService.GenerateTip(stats.WinRate, stats.TotalBets, topSport, stats.TotalProfit)
	stats.AiTip = tip

	return c.JSON(stats)
}

// GetTransactionsHandler obtiene el extracto bancario.
// @Router /api/transactions [get]
func (h *Handler) GetTransactionsHandler(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)

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

// Estructuras de Respuesta para Docs y JSON

type DashboardStatsResponse struct {
	TotalBets        int64       `json:"total_bets"`
	WonBets          int64       `json:"won_bets"`
	WinRate          float64     `json:"win_rate"`
	TotalProfit      float64     `json:"total_profit"`
	CurrentBankroll  float64     `json:"current_bankroll"`
	AiTip            string      `json:"ai_tip"`
	SportPerformance []SportStat `json:"sport_performance"`
}

type SportStat struct {
	SportKey string  `json:"sport_key"`
	Bets     int     `json:"bets"`
	Profit   float64 `json:"profit"`
}

type ResolveMatchRequest struct {
	MatchID string `json:"match_id"`
	Winner  string `json:"winner"` // "HOME" o "AWAY"
}

// SettleMatchHandler (Endpoint Admin)
func (h *Handler) SettleMatchHandler(c *fiber.Ctx) error {
	var req ResolveMatchRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "JSON inválido"})
	}

	matchUUID, err := uuid.Parse(req.MatchID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID de partido inválido"})
	}

	if req.Winner != "HOME" && req.Winner != "AWAY" {
		return c.Status(400).JSON(fiber.Map{"error": "El ganador debe ser HOME o AWAY"})
	}

	err = h.service.SettleMatch(matchUUID, req.Winner)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message":  "Proceso de liquidación completado",
		"match_id": req.MatchID,
		"winner":   req.Winner,
	})
}

// GetService permite acceder al servicio interno (usado por el worker)
func (h *Handler) GetService() *Service {
	return h.service
}
