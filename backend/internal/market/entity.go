package market

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Match representa un partido real traído de la API
type Match struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Identificadores Externos
	ExternalID string `gorm:"uniqueIndex;not null" json:"external_id"` // El ID 1621465540 de Pinnacle
	Provider   string `gorm:"default:'pinnacle'" json:"provider"`

	// Detalles del Partido
	SportKey string    `json:"sport_key"` // "lol", "dota2", "cs2" (Lo inferiremos de la liga)
	League   string    `json:"league"`    // "Demacia Cup"
	HomeTeam string    `json:"home_team"` // "Oh My God"
	AwayTeam string    `json:"away_team"` // "JD Gaming"
	StartsAt time.Time `json:"starts_at"` // Cuándo juega

	// Cuotas (Odds) - Solo guardamos Ganador (MoneyLine) por ahora
	HomeOdds float64 `json:"home_odds"`
	AwayOdds float64 `json:"away_odds"`

	// Estado
	Status string `gorm:"default:'scheduled'" json:"status"` // scheduled, live, finished
}

func (Match) TableName() string {
	return "matches" // <-- ASEGÚRATE de que este sea el nombre exacto en tu pgAdmin
}
