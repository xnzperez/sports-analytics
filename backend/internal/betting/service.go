package betting

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
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
	Title      string         `json:"title"`
	MatchDate  time.Time      `json:"match_date"`
	SportKey   string         `json:"sport_key"` // 'nba', 'cs2'
	StakeUnits float64        `json:"stake_units"`
	Odds       float64        `json:"odds"`
	IsParlay   bool           `json:"is_parlay"`
	Details    datatypes.JSON `json:"details"`    // Tus campos (teams, market, legs...)
	UserNotes  string         `json:"user_notes"` // Opcional
}

// PlaceBet maneja la creación de la apuesta y el descuento de saldo
func (s *Service) PlaceBet(userID uuid.UUID, req PlaceBetRequest) (*Bet, error) {
	var newBet *Bet

	err := s.repo.RunTransaction(func(tx *gorm.DB) error {
		// 1. Bloqueo y obtención de usuario
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

		// 4. Crear la Apuesta
		newBet = &Bet{
			UserID:     userID, // CORREGIDO: Pasamos el UUID directo, sin .String()
			Title:      req.Title,
			SportKey:   req.SportKey,
			StakeUnits: req.StakeUnits,
			Odds:       req.Odds,
			IsParlay:   req.IsParlay,
			Status:     "pending",
			Details:    string(req.Details),
			UserNotes:  req.UserNotes,
		}
		if err := s.repo.CreateBet(tx, newBet); err != nil {
			return err
		}

		// 5. Registrar Transacción (Ledger)
		transaction := &Transaction{
			UserID:      userID, // CORREGIDO: UUID directo
			Amount:      -req.StakeUnits,
			Type:        "BET_PLACED",
			Description: "Apuesta realizada: " + req.Title,
			ReferenceID: &newBet.ID, // CORREGIDO: Ahora los tipos coinciden (*uuid.UUID)
		}

		// Nota: Asegúrate de usar tx.Create, no s.repo... para mantener la atomicidad
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
