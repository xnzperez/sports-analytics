package analytics

import (
	"fmt"
	"math/rand"
	"time"
)

type AdvisorResult struct {
	Message string
	Level   string
}

type StatsInput struct {
	WinRate     float64
	TotalBets   int
	TotalProfit float64
	Bankroll    float64
}

func GenerateSmartTip(stats StatsInput) AdvisorResult {
	rand.Seed(time.Now().UnixNano())

	// 1. Fase de Recolección
	if stats.TotalBets < 5 {
		return AdvisorResult{
			Message: "Fase de aprendizaje: Estoy analizando tus primeros movimientos. Necesito 5 registros para activar el motor de rentabilidad.",
			Level:   "info",
		}
	}

	// 2. Gestión de Crisis (Profit Negativo)
	if stats.TotalProfit < 0 {
		if stats.WinRate > 55 {
			return AdvisorResult{
				Message: "⚠️ Paradoja detectada: Ganas muchas apuestas pero pierdes dinero. Estás sobre-apostando a cuotas muy bajas que no compensan el riesgo. ¡Busca más valor!",
				Level:   "warning",
			}
		}
		return AdvisorResult{
			Message: "Alerta de varianza: Tu estrategia actual está drenando el bankroll. Te sugiero bajar el Stake al 1% hasta recuperar el 50% de WinRate.",
			Level:   "warning",
		}
	}

	// 3. Optimización de Ganancias (Profit Positivo)
	if stats.TotalProfit > 0 {
		// Cálculo del Stake Sugerido (Kelly simplificado al 2%)
		suggestedStake := stats.Bankroll * 0.02

		if stats.WinRate < 40 {
			return AdvisorResult{
				Message: fmt.Sprintf("🎯 Estilo Francotirador: Pocos aciertos pero de gran valor. Mantén tu gestión de banca. Tu apuesta ideal hoy es de $%.2f.", suggestedStake),
				Level:   "success",
			}
		}

		// Mensajes aleatorios para éxito para que no sea repetitivo
		successMessages := []string{
			fmt.Sprintf("🚀 Sistema Sólido: Estás batiendo al mercado. Mantén el stake en $%.2f para un crecimiento compuesto.", suggestedStake),
			"🔥 ¡Racha detectada! Tus análisis de E-Sports están siendo precisos. No aumentes el riesgo por euforia.",
			fmt.Sprintf("💰 Gestión eficiente: Tu curva de profit es saludable. Sigue el plan de $%.2f por unidad.", suggestedStake),
		}

		return AdvisorResult{
			Message: successMessages[rand.Intn(len(successMessages))],
			Level:   "success",
		}
	}

	return AdvisorResult{
		Message: "Estás en el punto de equilibrio. Es momento de ser más selectivo con las ligas de e-Sports.",
		Level:   "info",
	}
}
