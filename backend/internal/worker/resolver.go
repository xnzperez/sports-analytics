package worker

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/xnzperez/sports-analytics-backend/internal/betting"
)

// Estructura auxiliar para leer el JSON de details
type BetDetails struct {
	MatchID   string `json:"match_id"`
	Selection string `json:"selection"` // "HOME" o "AWAY"
	TeamName  string `json:"team_name"`
}

func StartScheduler(service *betting.Service) {
	// Ejecutar cada 10 segundos
	ticker := time.NewTicker(10 * time.Second)

	go func() {
		fmt.Println("🤖 [WORKER] Auto-Resolver: Iniciado. Buscando apuestas para liquidar...")
		for range ticker.C {
			processPendingBets(service)
		}
	}()
}

func processPendingBets(service *betting.Service) {
	bets, err := service.GetPendingBets()
	if err != nil {
		fmt.Println("❌ [WORKER] Error buscando apuestas:", err)
		return
	}

	if len(bets) == 0 {
		return
	}

	fmt.Printf("🔍 [WORKER] Analizando %d apuestas pendientes...\n", len(bets))

	for _, bet := range bets {
		var details BetDetails
		if err := json.Unmarshal([]byte(bet.Details), &details); err != nil {
			continue
		}

		// 1. Simulamos quién ganó el partido (HOME o AWAY)
		matchWinner := simulateWinner(bet.ID.String())

		fmt.Printf("🎲 [SIMULACIÓN] Partido %s finalizado. Ganador del Match: %s\n", details.MatchID, matchWinner)

		// 2. LÓGICA DE CORRECCIÓN: Comparamos selección vs ganador
		// Aquí traducimos "HOME/AWAY" a "WON/LOST"
		betOutcome := "LOST" // Por defecto perdió

		if details.Selection == matchWinner {
			betOutcome = "WON" // Si coinciden, ganó
		}

		// 3. Enviamos el estado CORRECTO a la base de datos
		err := service.ResolveBet(bet.ID.String(), betOutcome)

		if err != nil {
			fmt.Printf("❌ [WORKER] Error resolviendo apuesta %s: %v\n", bet.ID, err)
		} else {
			fmt.Printf("💰 [WORKER] Apuesta %s liquidada. Usuario apostó %s -> Resultado: %s\n",
				bet.ID, details.Selection, betOutcome)
		}
	}
}

// simulateWinner decide aleatoriamente quién ganó (HOME o AWAY)
func simulateWinner(seed string) string {
	hash := 0
	for _, char := range seed {
		hash += int(char)
	}
	if hash%2 == 0 {
		return "HOME"
	}
	return "AWAY"
}
