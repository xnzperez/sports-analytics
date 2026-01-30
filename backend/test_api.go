package main

import (
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/xnzperez/sports-analytics-backend/internal/integrations/pinnacle"
)

const (
	eventDisplayLimit = 10
	timeLayoutISO     = "2006-01-02T15:04:05"
)

func main() {
	// 1. Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("Nota: No se encontró archivo .env, verificando variables del sistema...")
	}

	// 2. Inicializar cliente
	client := pinnacle.NewClient()

	fmt.Println("--- Consultando Partidos de E-Sports ---")

	// 3. Obtener datos
	resp, err := client.GetEsportsMarkets()
	if err != nil {
		log.Fatalf("Error obteniendo mercados: %v", err)
	}

	// 4. Mostrar resumen
	fmt.Printf("Deporte: %s (ID: %d)\n", resp.SportName, resp.SportID)
	totalEvents := len(resp.Events)
	fmt.Printf("Total de Eventos Encontrados: %d\n", totalEvents)
	fmt.Println("------------------------------------------------")

	if totalEvents == 0 {
		fmt.Println("No hay eventos disponibles en este momento.")
		return
	}

	// 5. Iterar y mostrar detalles (Max 10)
	shownCount := 0
	for _, ev := range resp.Events {
		if shownCount >= eventDisplayLimit {
			break
		}

		// Parsear fecha de forma segura
		t, err := time.Parse(timeLayoutISO, ev.Starts)
		dateStr := t.Format("02 Jan 15:04")
		if err != nil {
			dateStr = ev.Starts // Fallback al string original si falla el parseo
		}

		homeOdd := ev.Periods.Num0.MoneyLine.Home
		awayOdd := ev.Periods.Num0.MoneyLine.Away

		// Mostrar si las cuotas son validas
		if homeOdd == 0 || awayOdd == 0 {
			continue 
		}

		fmt.Printf("[%d] %s\n", ev.EventID, ev.LeagueName)
		fmt.Printf("    🎮 %s vs %s\n", ev.Home, ev.Away)
		fmt.Printf("    💰 Cuotas: %.2f | %.2f\n", homeOdd, awayOdd)
		fmt.Printf("    📅 Fecha: %s\n\n", dateStr)

		shownCount++
	}
}
