package pinnacle

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Client maneja la comunicación con Pinnacle via RapidAPI
type Client struct {
	apiKey     string
	apiHost    string
	httpClient *http.Client
	baseURL    string
}

func NewClient() *Client {
	return &Client{
		apiKey:     os.Getenv("RAPIDAPI_KEY"),
		apiHost:    os.Getenv("RAPIDAPI_HOST_PINNACLE"),
		baseURL:    "https://pinnacle-odds.p.rapidapi.com",
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) makeRequest(method, endpoint string) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, endpoint)

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("x-rapidapi-key", c.apiKey)
	req.Header.Add("x-rapidapi-host", c.apiHost)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	if res.StatusCode != 200 {
		// Imprimimos el error para debuggear
		fmt.Printf("⚠️  RAPIDAPI ERROR: %s\n", string(body))
		return nil, fmt.Errorf("API Error: Status %d", res.StatusCode)
	}

	return body, nil
}

// --- ESTRUCTURAS BASADAS EN LA DOCUMENTACIÓN ---

type MoneyLine struct {
	Home float64 `json:"home"`
	Away float64 `json:"away"`
	Draw float64 `json:"draw,omitempty"` // En Esports a veces no hay empate, por eso omitempty
}

type Period0 struct {
	MoneyLine MoneyLine `json:"money_line"`
	Cutoff    string    `json:"cutoff"` // Fecha límite para apostar
}

type Periods struct {
	Num0 Period0 `json:"num_0"` // num_0 = Match Winner (Partido completo)
}

type Event struct {
	EventID    int64   `json:"event_id"`
	LeagueName string  `json:"league_name"`
	Starts     string  `json:"starts"` // Fecha ISO
	Home       string  `json:"home"`
	Away       string  `json:"away"`
	Periods    Periods `json:"periods"`
}

type MarketsResponse struct {
	SportID   int     `json:"sport_id"`
	SportName string  `json:"sport_name"`
	Last      int64   `json:"last"` // Timestamp para paginación 'since'
	Events    []Event `json:"events"`
}

// GetEsportsMarkets trae los partidos de Esports (ID 12)
// Documentación: /kit/v1/markets
func (c *Client) GetEsportsMarkets() (*MarketsResponse, error) {
	// CAMBIO: sport_id=10 (Esports según tu descubrimiento)
	endpoint := "/kit/v1/markets?sport_id=10&is_have_odds=true"

	fmt.Println("DEBUG: Consultando Esports ->", endpoint)

	body, err := c.makeRequest("GET", endpoint)
	if err != nil {
		return nil, err
	}

	var response MarketsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetSports obtiene la lista de todos los deportes disponibles y sus IDs
func (c *Client) GetSports() ([]byte, error) {
	// Endpoint estándar según tu documentación: @List of sports
	return c.makeRequest("GET", "/kit/v1/sports")
}
