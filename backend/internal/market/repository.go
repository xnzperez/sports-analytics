package market

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// SaveMatch guarda el partido traído de la API.
// Usa Clauses(clause.OnConflict) para evitar errores si sincronizas dos veces.
func (r *Repository) SaveMatch(match *Match) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "external_id"}},                                    // Si el ID de Pinnacle ya existe...
		DoUpdates: clause.AssignmentColumns([]string{"home_odds", "away_odds", "starts_at"}), // Actualiza cuotas y fecha
	}).Create(match).Error
}

// GetMatches devuelve la lista para que el frontend la vea
func (r *Repository) GetMatches() ([]Match, error) {
	var matches []Match
	// Ordenamos por fecha de inicio
	result := r.db.Order("starts_at asc").Find(&matches)
	return matches, result.Error
}
