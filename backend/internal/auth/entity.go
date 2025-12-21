package auth

import (
	"time"

	"github.com/google/uuid"
)

// User representa al usuario en nuestro sistema.
// Notarás los 'tags' `gorm:"..."` y `json:"..."`.
// GORM usa los primeros para saber a qué columna mapear.
// Fiber usa los segundos para saber cómo mostrarlo en la API (JSON).
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Email        string    `gorm:"unique;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	Username     string    `gorm:"unique;not null" json:"username"`

	// --- CAMBIO: Simplificación para MVP ---
	// Usamos un solo campo 'Bankroll' para que coincida con el Frontend (json:"bankroll")
	// type:decimal(15,2) asegura precisión monetaria en la base de datos
	Bankroll float64 `gorm:"default:0.00;type:decimal(15,2)" json:"bankroll"`
	// ---------------------------------------

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
