package market

import (
	"fmt"
	"strings"
	"time"

	"github.com/xnzperez/sports-analytics-backend/internal/integrations/pinnacle"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Service struct {
	db             *gorm.DB
	pinnacleClient *pinnacle.Client
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		db:             db,
		pinnacleClient: pinnacle.NewClient(),
	}
}

// SyncEsports llama a la API y actualiza nuestra base de datos local
func (s *Service) SyncEsports() (int, error) {
	// 1. Llamar a RapidAPI (Gasta 1 crédito)
	resp, err := s.pinnacleClient.GetEsportsMarkets()
	if err != nil {
		return 0, err
	}

	count := 0

	// 2. Procesar cada evento y guardarlo
	for _, event := range resp.Events {

		// Filtrar cuotas vacías (0.00)
		if event.Periods.Num0.MoneyLine.Home == 0 || event.Periods.Num0.MoneyLine.Away == 0 {
			continue
		}

		// Inferir deporte basado en la liga o nombres (Lógica simple por ahora)
		sportKey := "esports"
		leagueLower := strings.ToLower(event.LeagueName)
		if strings.Contains(leagueLower, "lol") || strings.Contains(leagueLower, "league of legends") || strings.Contains(leagueLower, "lck") || strings.Contains(leagueLower, "lpl") {
			sportKey = "lol"
		} else if strings.Contains(leagueLower, "dota") {
			sportKey = "dota2"
		} else if strings.Contains(leagueLower, "cs2") || strings.Contains(leagueLower, "counter-strike") {
			sportKey = "cs2"
		} else if strings.Contains(leagueLower, "valorant") {
			sportKey = "valorant"
		}

		// Parsear fecha
		startTime, _ := time.Parse("2006-01-02T15:04:05", event.Starts)

		match := Match{
			ExternalID: filterNumericID(event.EventID), // Convertir ID a string
			Provider:   "pinnacle",
			SportKey:   sportKey,
			League:     event.LeagueName,
			HomeTeam:   event.Home,
			AwayTeam:   event.Away,
			StartsAt:   startTime,
			HomeOdds:   event.Periods.Num0.MoneyLine.Home,
			AwayOdds:   event.Periods.Num0.MoneyLine.Away,
			Status:     "scheduled",
		}

		// UPSERT: Si ya existe (mismo ExternalID), actualiza las cuotas. Si no, lo crea.
		err := s.db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "external_id"}},                                     // Buscar por ID externo
			DoUpdates: clause.AssignmentColumns([]string{"home_odds", "away_odds", "updated_at"}), // Actualizar solo cuotas
		}).Create(&match).Error

		if err == nil {
			count++
		}
	}

	return count, nil
}

// GetMatches devuelve los partidos guardados en DB para el frontend
func (s *Service) GetMatches(sport string) ([]Match, error) {
	var matches []Match
	query := s.db.Where("starts_at > ?", time.Now()).Order("starts_at asc")

	if sport != "" && sport != "all" {
		query = query.Where("sport_key = ?", sport)
	}

	err := query.Find(&matches).Error
	return matches, err
}

// Helper para convertir int64 a string
func filterNumericID(id int64) string {
	return fmt.Sprintf("%d", id)
}
