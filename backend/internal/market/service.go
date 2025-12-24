package market

import (
	"fmt"
	"strings"
	"time"

	"github.com/xnzperez/sports-analytics-backend/internal/integrations/pinnacle"
	"gorm.io/gorm"
)

type Service struct {
	repo           *Repository // CAMBIO: Ahora usamos el Repository en lugar de db directo
	pinnacleClient *pinnacle.Client
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		repo:           NewRepository(db), // Inicializamos el Repo
		pinnacleClient: pinnacle.NewClient(),
	}
}

// SyncEsports llama a la API y actualiza nuestra base de datos local usando el Repo
func (s *Service) SyncEsports() (int, error) {
	// 1. Llamar a RapidAPI
	resp, err := s.pinnacleClient.GetEsportsMarkets()
	if err != nil {
		return 0, err
	}

	count := 0

	// 2. Procesar cada evento
	for _, event := range resp.Events {

		// Filtrar cuotas vacías (0.00)
		if event.Periods.Num0.MoneyLine.Home == 0 || event.Periods.Num0.MoneyLine.Away == 0 {
			continue
		}

		// Inferir deporte
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
			ExternalID: filterNumericID(event.EventID),
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

		// CAMBIO: Usamos el repositorio para guardar
		if err := s.repo.SaveMatch(&match); err == nil {
			count++
		}
	}

	return count, nil
}

// GetMatches devuelve TODOS los partidos (Delegamos al Repo)
func (s *Service) GetMatches(sport string) ([]Match, error) {
	// NOTA: Ignoramos el filtro de sport por ahora para asegurar
	// que veas partidos aunque no coincidan con 'lol'.
	// Y lo más importante: ¡Ya no filtramos por fecha!
	return s.repo.GetMatches()
}

// GetAvailableMatches es un alias por si tu Handler lo llama con este nombre
func (s *Service) GetAvailableMatches() ([]Match, error) {
	return s.repo.GetMatches()
}

// Helper para convertir int64 a string
func filterNumericID(id int64) string {
	return fmt.Sprintf("%d", id)
}
