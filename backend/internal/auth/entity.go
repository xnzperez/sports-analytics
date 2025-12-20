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
	PasswordHash string    `gorm:"not null" json:"-"` // json:"-" oculta el pass en las respuestas API
	Username     string    `gorm:"unique;not null" json:"username"`

	// Manejo de Dinero (Bankroll)
	// Usamos float64 para mapear DECIMAL en Go, aunque para cálculos
	// precisos se recomienda librerías como 'shopspring/decimal'.
	// Para este MVP, float64 es suficiente si redondeamos bien.
	BankrollUnits    float64 `gorm:"default:0.00" json:"bankroll_units"`
	BankrollCurrency float64 `gorm:"default:0.00" json:"bankroll_currency"`
	CurrencyCode     string  `gorm:"default:'USD'" json:"currency_code"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
