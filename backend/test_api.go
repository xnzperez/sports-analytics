package main

import (
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/xnzperez/sports-analytics-backend/internal/integrations/pinnacle"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	client := pinnacle.NewClient()

	fmt.Println("--- Consultando Partidos de E-Sports (ID 10) ---")

	resp, err := client.GetEsportsMarkets()
	if err != nil {
		log.Fatalf("Error crítico: %v", err)
	}

	fmt.Printf("Deporte: %s (ID: %d)\n", resp.SportName, resp.SportID)
	fmt.Printf("Total de Eventos: %d\n", len(resp.Events))
	fmt.Println("------------------------------------------------")

	// Mostramos los primeros 10 partidos para ver si sale LoL/Valorant
	limit := 10
	if len(resp.Events) < limit {
		limit = len(resp.Events)
	}

	for i := 0; i < limit; i++ {
		ev := resp.Events[i]
		ml := ev.Periods.Num0.MoneyLine

		// Convertir fecha fea a legible
		t, _ := time.Parse("2006-01-02T15:04:05", ev.Starts)

		fmt.Printf("[%d] %s\n", ev.EventID, ev.LeagueName)
		fmt.Printf("    🎮 %s vs %s\n", ev.Home, ev.Away)
		fmt.Printf("    💰 Cuotas: %.2f | %.2f\n", ml.Home, ml.Away)
		fmt.Printf("    📅 Fecha: %s\n\n", t.Format("02 Jan 15:04"))
	}
}
