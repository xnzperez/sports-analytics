package betting

import (
	"errors"

	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/xnzperez/sports-analytics-backend/internal/analytics"
	"gorm.io/gorm"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// PlaceBetRequest es el JSON que recibiremos del Frontend
type PlaceBetRequest struct {
	Title      string  `json:"title"`
	SportKey   string  `json:"sport_key"`
	StakeUnits float64 `json:"stake_units"`
	Odds       float64 `json:"odds"`
	IsParlay   bool    `json:"is_parlay"`
	UserNotes  string  `json:"user_notes"`

	// CAMBIO AQUÍ: Usar map[string]interface{} es más seguro para lo que envía Zod
	Details map[string]interface{} `json:"details"`
}

// PlaceBet maneja la creación de la apuesta y el descuento de saldo
func (s *Service) PlaceBet(userID uuid.UUID, req PlaceBetRequest) (*Bet, error) {
	var newBet *Bet

	err := s.repo.RunTransaction(func(tx *gorm.DB) error {
		// 1. Bloqueo y obtención de usuario (Anti-Race Condition)
		user, err := s.repo.GetUserBalanceForUpdate(tx, userID)
		if err != nil {
			return err
		}

		// 2. Verificar Fondos
		if user.Bankroll < req.StakeUnits {
			return errors.New("saldo insuficiente para realizar esta apuesta")
		}

		// 3. Descontar Saldo
		newBalance := user.Bankroll - req.StakeUnits
		if err := s.repo.UpdateUserBalance(tx, user.ID, newBalance); err != nil {
			return err
		}

		// 4. Preparar Datos (JSON y ExternalID)
		detailsJSON := "{}"
		var externalID string

		// Extraemos datos clave del mapa genérico para indexarlos
		if req.Details != nil {
			// Intenta sacar el external_id para guardarlo en la columna indexada
			if val, ok := req.Details["external_id"].(string); ok {
				externalID = val
			}

			bytes, err := json.Marshal(req.Details)
			if err == nil {
				detailsJSON = string(bytes)
			}
		}

		// 5. Crear la Apuesta
		newBet = &Bet{
			UserID:     userID,
			Title:      req.Title,
			SportKey:   req.SportKey,
			StakeUnits: req.StakeUnits,
			Odds:       req.Odds,
			IsParlay:   req.IsParlay,
			Status:     "pending",
			Details:    detailsJSON,
			UserNotes:  req.UserNotes,

			// --- OPTIMIZACIÓN DE ESCALABILIDAD ---
			ExternalID: externalID, // Guardamos el ID real aquí
			Provider:   "pinnacle", // Asumimos pinnacle por defecto
			// -------------------------------------
		}

		// 6. Registrar Transacción (Ledger)
		transaction := &Transaction{
			UserID:      userID,
			Amount:      -req.StakeUnits, // Negativo porque sale dinero
			Type:        "BET_PLACED",
			Description: "Apuesta realizada: " + req.Title,
			ReferenceID: &newBet.ID,
		}

		// Guardamos todo usando la transacción (tx)
		if err := tx.Create(newBet).Error; err != nil {
			return err
		}
		// Nota: Asignamos el ID de la apuesta recién creada a la transacción
		transaction.ReferenceID = &newBet.ID

		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return newBet, nil
}

// ResolveBet conecta el Handler con el Repository para finalizar una apuesta.
func (s *Service) ResolveBet(betID string, outcome string) error {
	if betID == "" {
		return errors.New("el ID de la apuesta es obligatorio")
	}
	return s.repo.ResolveBet(betID, outcome)
}

// BetFilters define los criterios de búsqueda
type BetFilters struct {
	UserID   uuid.UUID
	Status   string // "pending", "WON", "LOST"
	SportKey string // "cs2", "nba"
	Page     int
	Limit    int
}

// GetBetsResponse estructura la respuesta paginada
type GetBetsResponse struct {
	Data  []Bet `json:"data"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}

func (s *Service) GetBets(filters BetFilters) (*GetBetsResponse, error) {
	// Valores por defecto para evitar errores de división por cero o cargas masivas
	if filters.Page <= 0 {
		filters.Page = 1
	}
	if filters.Limit <= 0 || filters.Limit > 100 {
		filters.Limit = 10
	}

	// Llamamos al repositorio
	bets, total, err := s.repo.GetBets(filters)
	if err != nil {
		return nil, err
	}

	return &GetBetsResponse{
		Data:  bets,
		Total: total,
		Page:  filters.Page,
		Limit: filters.Limit,
	}, nil
}

// StatsResponse es el reporte de rendimiento del usuario
type StatsResponse struct {
	TotalBets    int64 `json:"total_bets"`
	TotalWon     int64 `json:"total_won"`
	TotalLost    int64 `json:"total_lost"`
	TotalPending int64 `json:"total_pending"`

	WinRate float64 `json:"win_rate"` // % de aciertos

	TotalWagered  float64 `json:"total_wagered"`  // Total apostado
	TotalReturned float64 `json:"total_returned"` // Total recibido (ganancias + stake devuelto)
	NetProfit     float64 `json:"net_profit"`     // Ganancia/Pérdida neta
	ROI           float64 `json:"roi"`            // Retorno de Inversión (%)
}

// GetUserStats calcula las estadísticas financieras y de rendimiento
func (s *Service) GetUserStats(userID uuid.UUID) (*StatsResponse, error) {
	// 1. Obtener los datos crudos del repositorio
	stats, err := s.repo.GetRawStats(userID)
	if err != nil {
		return nil, err
	}

	// 2. Calcular Métricas Derivadas (Lógica de Negocio)
	response := &StatsResponse{
		TotalBets:    stats.TotalBets,
		TotalWon:     stats.Won,
		TotalLost:    stats.Lost,
		TotalPending: stats.Pending,
		TotalWagered: stats.TotalWagered,
	}

	// A. Calcular Win Rate (Evitar división por cero)
	// Consideramos solo las apuestas resueltas (Ganadas + Perdidas) para el Winrate real
	settledBets := stats.Won + stats.Lost
	if settledBets > 0 {
		response.WinRate = (float64(stats.Won) / float64(settledBets)) * 100
	}

	// B. Calcular Retorno y Profit
	// El repositorio nos dará cuánto dinero "volvió".
	// Profit = Lo que volvió - Lo que aposté (en total)
	response.TotalReturned = stats.TotalReturned
	response.NetProfit = response.TotalReturned - response.TotalWagered

	// C. Calcular ROI (Return On Investment)
	// Fórmula: (Profit / Total Apostado) * 100
	if response.TotalWagered > 0 {
		response.ROI = (response.NetProfit / response.TotalWagered) * 100
	}

	return response, nil
}

// GetTransactionsResponse define cómo entregamos los datos al cliente
type GetTransactionsResponse struct {
	Data  []Transaction `json:"data"`
	Total int64         `json:"total"`
	Page  int           `json:"page"`
	Limit int           `json:"limit"`
}

func (s *Service) GetTransactions(userID uuid.UUID, page, limit int) (*GetTransactionsResponse, error) {
	// Defaults de seguridad
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	txs, total, err := s.repo.GetTransactions(userID, page, limit)
	if err != nil {
		return nil, err
	}

	return &GetTransactionsResponse{
		Data:  txs,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

// GetUserDashboardStats calcula las estadísticas, aplicando filtro opcional de deporte
func (s *Service) GetUserDashboardStats(userID uuid.UUID, sportFilter string) (*DashboardStatsResponse, error) {
	var bets []Bet

	// 1. Construir la Query Base
	query := s.repo.db.Where("user_id = ?", userID)

	// 2. Aplicar el Filtro si existe y no es "all"
	// (El frontend envía "all" cuando quiere ver todo, o vacío)
	if sportFilter != "" && sportFilter != "all" {
		query = query.Where("sport_key = ?", sportFilter)
	}

	// 3. Ejecutar la consulta
	result := query.Order("created_at desc").Find(&bets)
	if result.Error != nil {
		return nil, result.Error
	}

	// 4. Calcular Métricas en Memoria
	var totalBets int64 = int64(len(bets))
	var wonBets int64 = 0
	var totalProfit float64 = 0.0

	// Mapa para agrupar rendimiento por deporte
	sportMap := make(map[string]*SportStat)

	for _, bet := range bets {
		// Inicializar mapa del deporte si no existe
		if _, exists := sportMap[bet.SportKey]; !exists {
			sportMap[bet.SportKey] = &SportStat{SportKey: bet.SportKey}
		}

		currentSportStat := sportMap[bet.SportKey]
		currentSportStat.Bets++

		// Calcular Ganancias/Pérdidas
		// Asumimos que si está "won", ganamos (stake * odds) - stake
		// Si está "lost", perdemos el stake.
		// Si está "pending", no afecta el profit todavía.
		if bet.Status == "WON" {
			wonBets++
			profit := (bet.StakeUnits * bet.Odds) - bet.StakeUnits
			totalProfit += profit
			currentSportStat.Profit += profit
		} else if bet.Status == "LOST" {
			totalProfit -= bet.StakeUnits
			currentSportStat.Profit -= bet.StakeUnits
		}
	}

	// 5. Calcular WinRate
	var winRate float64 = 0
	// Solo contamos apuestas resueltas para el WinRate real (evitamos dividir por pendientes)
	resolvedBets := 0
	for _, b := range bets {
		if b.Status == "WON" || b.Status == "LOST" {
			resolvedBets++
		}
	}

	if resolvedBets > 0 {
		winRate = (float64(wonBets) / float64(resolvedBets)) * 100
	}

	// 6. Obtener Bankroll actual del usuario (siempre el total real)
	user, _ := s.repo.GetUserByID(userID)
	currentBankroll := 0.0
	if user != nil {
		currentBankroll = user.Bankroll
	}

	// Convertir el mapa de deportes a slice para la respuesta
	var sportPerformance []SportStat
	for _, stat := range sportMap {
		sportPerformance = append(sportPerformance, *stat)
	}

	// 7. Generar Tip Inteligente (Lógica Local)
	input := analytics.StatsInput{
		WinRate:     winRate,
		TotalBets:   int(totalBets),
		TotalProfit: totalProfit,
		Bankroll:    currentBankroll,
	}

	advice := analytics.GenerateSmartTip(input)
	aiTip := advice.Message

	return &DashboardStatsResponse{
		TotalBets:        totalBets,
		WonBets:          wonBets,
		WinRate:          winRate,
		TotalProfit:      totalProfit,
		CurrentBankroll:  currentBankroll,
		AiTip:            aiTip,
		SportPerformance: sportPerformance,
	}, nil
}

// Struct auxiliar para leer el JSON que guardamos en 'details'
type BetDetails struct {
	MatchID   string `json:"match_id"`
	Selection string `json:"selection"` // "HOME" o "AWAY"
	TeamName  string `json:"team_name"`
}

// SettleMatch resuelve todas las apuestas de un partido específico
// winner: "HOME" o "AWAY"
// SettleMatch resuelve todas las apuestas de un partido específico
func (s *Service) SettleMatch(matchID uuid.UUID, winner string) error {
	var bets []Bet

	// Traemos solo pendientes para optimizar
	if err := s.repo.db.Where("status = ?", "pending").Find(&bets).Error; err != nil {
		return err
	}

	resolvedCount := 0
	for _, bet := range bets {
		var details BetDetails
		if err := json.Unmarshal([]byte(bet.Details), &details); err != nil {
			continue
		}

		// Match exacto
		if details.MatchID != matchID.String() {
			continue
		}

		// Determinar resultado
		newStatus := "LOST"
		if details.Selection == winner {
			newStatus = "WON"
		}

		// Resolver atómicamente
		if err := s.repo.ResolveBet(bet.ID.String(), newStatus); err == nil {
			resolvedCount++
		}
	}

	if resolvedCount > 0 {
		log.Printf("✅ %d apuestas resueltas para el partido %s", resolvedCount, matchID)
	}

	return nil
}

// GetPendingBets expone las apuestas pendientes para el worker
func (s *Service) GetPendingBets() ([]Bet, error) {
	return s.repo.GetPendingBets()
}
