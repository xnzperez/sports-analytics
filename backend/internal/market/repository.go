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

func (r *Repository) SaveMatch(match *Match) error {
	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "external_id"}},
		// CORREGIDO: Usamos "starts_at" que es el nombre real en tu DB
		DoUpdates: clause.AssignmentColumns([]string{"home_odds", "away_odds", "starts_at", "status"}),
	}).Create(match).Error
}

func (r *Repository) GetMatches() ([]Match, error) {
	var matches []Match
	// Eliminamos el debug y el .Unscoped() (ya vimos que no era soft-delete)
	// Filtramos por status 'scheduled' para no mostrar partidos terminados en el modal
	result := r.db.Where("status = ?", "scheduled").Order("starts_at asc").Find(&matches)
	return matches, result.Error
}
