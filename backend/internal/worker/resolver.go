package worker

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/xnzperez/sports-analytics-backend/internal/betting"
)

// Estructura auxiliar para leer el JSON de details
type BetDetails struct {
	MatchID  string `json:"match_id"`
	TeamName string `json:"team_name"`
}

func StartScheduler(service *betting.Service) {
	// Ejecutar cada 30 segundos para ver resultados rápido en pruebas
	ticker := time.NewTicker(30 * time.Second)

	go func() {
		fmt.Println("🤖 Auto-Resolver: Iniciado y vigilando...")

		for range ticker.C {
			processPendingBets(service)
		}
	}()
}

func processPendingBets(service *betting.Service) {
	// 1. Obtener pendientes
	bets, err := service.GetPendingBets()
	if err != nil {
		fmt.Println("❌ Error buscando apuestas pendientes:", err)
		return
	}

	if len(bets) == 0 {
		return
	}

	fmt.Printf("🔍 Auto-Resolver: Analizando %d apuestas pendientes...\n", len(bets))

	for _, bet := range bets {
		// 2. Extraer ID del partido de los detalles
		var details BetDetails
		if err := json.Unmarshal([]byte(bet.Details), &details); err != nil {
			fmt.Printf("⚠️ Error leyendo detalles de apuesta %s: %v\n", bet.ID, err)
			continue
		}

		if details.MatchID == "" {
			continue
		}

		// 3. CONSULTAR RESULTADO (Simulación)
		winner := checkExternalResult(details.MatchID, details.TeamName)

		// 4. Si hay resultado, resolvemos la apuesta
		if winner != "" {
			fmt.Printf("✅ PARTIDO FINALIZADO DETECTADO (%s). Ganador: %s\n", details.MatchID, winner)

			err := service.ResolveBet(bet.ID.String(), winner)
			if err != nil {
				fmt.Printf("❌ Error resolviendo apuesta %s: %v\n", bet.ID, err)
			} else {
				fmt.Printf("💰 Apuesta %s actualizada a %s automáticamente.\n", bet.ID, winner)
			}
		}
	}
}

// checkExternalResult simula la API de resultados
func checkExternalResult(matchID string, teamName string) string {
	// ESTE ES EL ID QUE SACAMOS DE TU JSON:
	targetMatchID := "474d6868-a238-4f1c-99b3-305748f1d597"

	if matchID == targetMatchID {
		// Simulamos que la API externa dice que tu selección GANÓ
		return "WON"
	}

	return ""
}
