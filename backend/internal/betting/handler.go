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

// PlaceBet crea una nueva apuesta en el sistema.
// @Summary      Crear una nueva apuesta
// @Description  Permite al usuario registrar una apuesta, descontando saldo automáticamente.
// @Tags         Apuestas
// @Accept       json
// @Produce      json
// @Param        request body PlaceBetRequest true "Datos de la apuesta"
// @Success      201  {object}  map[string]interface{} "Apuesta creada exitosamente"
// @Failure      400  {object}  map[string]interface{} "Error de validación o saldo insuficiente"
// @Failure      500  {object}  map[string]interface{} "Error interno del servidor"
// @Router       /api/bets [post]
// @Security     Bearer
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

// ResolveBetHandler define el resultado de una apuesta.
// @Summary      Resolver una apuesta (Ganada/Perdida)
// @Description  Actualiza el estado de la apuesta. Si es WON, paga al usuario y genera la transacción.
// @Tags         Apuestas
// @Accept       json
// @Produce      json
// @Param        id       path      string             true  "ID de la apuesta (UUID)"
// @Param        request  body      ResolveBetRequest  true  "Resultado (WON/LOST)"
// @Success      200      {object}  map[string]interface{} "Apuesta resuelta correctamente"
// @Failure      400      {object}  map[string]interface{} "Error de validación"
// @Failure      500      {object}  map[string]interface{} "Error interno"
// @Router       /api/bets/{id}/resolve [patch]
// @Security     Bearer
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

// GetBetsHandler obtiene el historial de apuestas con filtros.
// @Summary      Listar apuestas del usuario
// @Description  Obtiene las apuestas paginadas. Permite filtrar por estado o deporte.
// @Tags         Apuestas
// @Accept       json
// @Produce      json
// @Param        page       query     int     false  "Número de página (default 1)"
// @Param        limit      query     int     false  "Items por página (default 10)"
// @Param        status     query     string  false  "Filtrar por estado (pending, WON, LOST)"
// @Param        sport_key  query     string  false  "Filtrar por deporte (cs2, nba)"
// @Success      200        {object}  GetBetsResponse
// @Failure      500        {object}  map[string]interface{} "Error interno"
// @Router       /api/bets [get]
// @Security     Bearer
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

// GetStatsHandler calcula el rendimiento financiero del usuario.
// @Summary      Obtener estadísticas (ROI, WinRate)
// @Description  Calcula métricas clave como ganancia neta, porcentaje de aciertos y total apostado.
// @Tags         Analytics
// @Accept       json
// @Produce      json
// @Success      200      {object}  StatsResponse
// @Failure      500      {object}  map[string]interface{} "Error interno"
// @Router       /api/stats [get]
// @Security     Bearer
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
// GetTransactionsHandler obtiene el extracto bancario.
// @Summary      Historial de Transacciones
// @Description  Muestra los movimientos de dinero (depósitos, apuestas, pagos) paginados.
// @Tags         Analytics
// @Accept       json
// @Produce      json
// @Param        page   query     int  false  "Número de página"
// @Param        limit  query     int  false  "Items por página"
// @Success      200    {object}  GetTransactionsResponse
// @Failure      500    {object}  map[string]interface{} "Error interno"
// @Router       /api/transactions [get]
// @Security     Bearer
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
