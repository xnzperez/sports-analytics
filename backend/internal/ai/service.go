package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Service struct {
	apiKey string
}

func NewService() *Service {
	// Intenta leer la KEY, si no hay, funcionará en modo "Lógica Local"
	return &Service{
		apiKey: os.Getenv("OPENAI_API_KEY"),
	}
}

// Estructuras para hablar con OpenAI
type OpenAIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"` // "system" o "user"
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

// GenerateTip analiza las estadísticas y devuelve un consejo
// Recibe los datos crudos del usuario (WinRate, Deporte más rentable, etc)
func (s *Service) GenerateTip(winRate float64, totalBets int64, topSport string, profit float64) string {

	// 1. Construimos el contexto del usuario
	prompt := fmt.Sprintf(
		"Analiza estos datos de apuestas: WinRate: %.2f%%, Total Apuestas: %d, Deporte Top: %s, Ganancia: $%.2f. Dame un consejo de 1 frase corta y motivadora o de precaución.",
		winRate, totalBets, topSport, profit,
	)

	// 2. Si tenemos API Key, preguntamos a la IA real
	if s.apiKey != "" {
		tip, err := s.callOpenAI(prompt)
		if err == nil {
			return "✨ IA: " + tip
		}
		fmt.Println("Error llamando a OpenAI (usando fallback):", err)
	}

	// 3. FALLBACK (Lógica Local Inteligente)
	// Si no hay API Key o falla, usamos lógica condicional avanzada
	// Esto hace que el sistema parezca inteligente inmediatamente.
	if totalBets == 0 {
		return "🤖 Empieza despacio. Analiza las estadísticas de los equipos antes de tu primera apuesta."
	}
	if winRate == 100 {
		return fmt.Sprintf("🔥 ¡Estás en racha perfecta en %s! Pero cuidado, no te confíes y mantén el stake.", topSport)
	}
	if profit > 0 {
		return fmt.Sprintf("📈 Tu estrategia en %s es rentable. Considera aumentar ligeramente el stake si mantienes el ritmo.", topSport)
	}
	if winRate < 40 {
		return "🛡️ Estás en una mala racha. Tómate un descanso y revisa tus replays."
	}

	return "📊 Diversifica tus apuestas para minimizar el riesgo."
}

func (s *Service) callOpenAI(prompt string) (string, error) {
	reqBody := OpenAIRequest{
		Model: "gpt-3.5-turbo", // O gpt-4
		Messages: []Message{
			{Role: "system", Content: "Eres un experto analista de apuestas deportivas (Esports). Responde conciso."},
			{Role: "user", Content: prompt},
		},
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("status code %d", resp.StatusCode)
	}

	var openAIResp OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openAIResp); err != nil {
		return "", err
	}

	if len(openAIResp.Choices) > 0 {
		return openAIResp.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no response choices")
}
