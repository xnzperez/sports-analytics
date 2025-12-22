package betting

import (
	"time"

	"github.com/google/uuid"
)

// Bet representa una apuesta en el sistema.
// Usamos tags de GORM para definir la estructura exacta en la base de datos
// y tags JSON para la respuesta de la API.
type Bet struct {
	ID     uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`

	Title    string `gorm:"not null" json:"title"`
	IsParlay bool   `gorm:"default:false" json:"is_parlay"`
	SportKey string `gorm:"not null" json:"sport_key"` // Ej: "cs2", "nba"
	Status   string `gorm:"default:'pending'" json:"status"`

	StakeUnits float64 `gorm:"not null" json:"stake_units"`
	Odds       float64 `gorm:"not null" json:"odds"`

	// Details se guarda como JSONB en Postgres para poder hacer consultas avanzadas dentro del JSON en el futuro.
	// En Go lo manejamos como string (o []byte) conteniendo el JSON crudo.
	Details string `gorm:"type:jsonb" json:"details"`

	UserNotes string `json:"user_notes"`

	// --- CORRECCIÓN ---
	// Usamos `column:ai_prediction` para evitar que GORM genere "a_iprediction".
	// Usamos *string (puntero) para que si no hay predicción, se guarde como NULL en la BD.
	AIPrediction *string `gorm:"column:ai_prediction" json:"ai_prediction,omitempty"`

	CreatedAt  time.Time  `gorm:"autoCreateTime" json:"created_at"`
	ResultedAt *time.Time `json:"resulted_at,omitempty"` // Puntero porque inicialmente es NULL

	ExternalID string `json:"external_id" gorm:"index"` // Index para búsquedas rápidas
	Provider   string `json:"provider"`                 // 'pinnacle', 'api-sports', etc.
}

// TableName anula la pluralización por defecto de GORM si fuera necesario,
// aunque GORM suele usar "bets" por defecto, es buena práctica ser explícito.
func (Bet) TableName() string {
	return "bets"
}

// Transaction representa cualquier movimiento de dinero en la cuenta del usuario.
// Esto es vital para auditoría y para mostrar el "Extracto Bancario".
type Transaction struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Amount      float64   `gorm:"not null" json:"amount"`
	Type        string    `gorm:"not null" json:"type"`
	Description string    `json:"description"`

	// CORREGIDO: Ahora es *uuid.UUID para coincidir con bet.ID
	ReferenceID *uuid.UUID `gorm:"type:uuid" json:"reference_id"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (Transaction) TableName() string {
	return "transactions"
}
