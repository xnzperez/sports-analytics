package betting

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/xnzperez/sports-analytics-backend/internal/auth"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// TransactionFunc define una función que se ejecuta dentro de una transacción
type TransactionFunc func(tx *gorm.DB) error

// RunTransaction es un wrapper para ejecutar bloques atómicos
func (r *Repository) RunTransaction(fn TransactionFunc) error {
	return r.db.Transaction(fn)
}

// GetUserBalanceForUpdate bloquea la fila del usuario y devuelve su saldo actual.
func (r *Repository) GetUserBalanceForUpdate(tx *gorm.DB, userID uuid.UUID) (*auth.User, error) {
	var user auth.User
	if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUserBalance actualiza el saldo del usuario dentro de la transacción
func (r *Repository) UpdateUserBalance(tx *gorm.DB, userID uuid.UUID, newBalance float64) error {
	// --- CORREGIDO: "bankroll_units" -> "bankroll" ---
	return tx.Model(&auth.User{}).Where("id = ?", userID).Update("bankroll", newBalance).Error
}

// CreateBet inserta la apuesta
func (r *Repository) CreateBet(tx *gorm.DB, bet *Bet) error {
	return tx.Create(bet).Error
}

// ResolveBet maneja la lógica de ganar/perder y actualiza el saldo atómicamente
func (r *Repository) ResolveBet(betIDStr string, outcome string) error {

	// 0. Convertir string a UUID (Validación inicial)
	betID, err := uuid.Parse(betIDStr)
	if err != nil {
		return errors.New("ID de apuesta inválido")
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Buscar y bloquear apuesta
		var bet Bet
		// GORM maneja la comparación uuid vs uuid automáticamente aquí
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&bet, "id = ?", betID).Error; err != nil {
			return err
		}

		// 2. Validación
		if bet.Status != "pending" {
			return errors.New("esta apuesta ya ha sido resuelta anteriormente")
		}

		// 3. Actualizar apuesta
		now := time.Now()
		bet.Status = outcome
		bet.ResultedAt = &now
		if err := tx.Save(&bet).Error; err != nil {
			return err
		}

		// 4. Si ganó, pagar y registrar transacción
		if outcome == "WON" {
			payout := bet.StakeUnits * bet.Odds

			// A. Actualizar Saldo Usuario
			// bet.UserID ya es UUID, GORM lo maneja bien en el Where
			if err := tx.Model(&auth.User{}).Where("id = ?", bet.UserID).
				Update("bankroll", gorm.Expr("bankroll + ?", payout)).Error; err != nil {
				return err
			}

			// B. Registrar Transacción (Ledger)
			transaction := &Transaction{
				UserID:      bet.UserID, // UUID directo
				Amount:      payout,
				Type:        "BET_PAYOUT",
				Description: "Ganancia apuesta: " + bet.Title,
				ReferenceID: &bet.ID, // Puntero a UUID (*uuid.UUID)
			}
			if err := tx.Create(transaction).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *Repository) GetBets(f BetFilters) ([]Bet, int64, error) {
	var bets []Bet
	var total int64

	// 1. Iniciar la query base
	query := r.db.Model(&Bet{}).Where("user_id = ?", f.UserID)

	// 2. Aplicar filtros dinámicos
	if f.Status != "" {
		query = query.Where("status = ?", f.Status)
	}
	if f.SportKey != "" {
		query = query.Where("sport_key = ?", f.SportKey)
	}

	// 3. Contar el total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 4. Aplicar Paginación y Ordenamiento
	offset := (f.Page - 1) * f.Limit
	err := query.Limit(f.Limit).Offset(offset).Order("created_at desc").Find(&bets).Error

	return bets, total, err
}

// RawStats es una estructura auxiliar para leer el resultado de la query
type RawStats struct {
	TotalBets     int64
	Won           int64
	Lost          int64
	Pending       int64
	TotalWagered  float64
	TotalReturned float64
}

func (r *Repository) GetRawStats(userID uuid.UUID) (*RawStats, error) {
	var stats RawStats

	// Usamos SQL nativo. He mantenido tu lógica de FILTER, pero ajustado para que sea robusto.
	// Asegúrate de que tu versión de Postgres soporte FILTER (Postgres 9.4+). Si no, usa CASE WHEN.
	err := r.db.Model(&Bet{}).
		Select(`
            COUNT(*) as total_bets,
            COUNT(*) FILTER (WHERE status = 'WON') as won,
            COUNT(*) FILTER (WHERE status = 'LOST') as lost,
            COUNT(*) FILTER (WHERE status = 'pending') as pending,
            COALESCE(SUM(stake_units), 0) as total_wagered,
            COALESCE(SUM(CASE WHEN status = 'WON' THEN stake_units * odds ELSE 0 END), 0) as total_returned
        `).
		Where("user_id = ?", userID).
		Scan(&stats).Error

	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// GetTransactions obtiene el historial financiero paginado
func (r *Repository) GetTransactions(userID uuid.UUID, page, limit int) ([]Transaction, int64, error) {
	var transactions []Transaction
	var total int64

	query := r.db.Model(&Transaction{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error

	return transactions, total, err
}
