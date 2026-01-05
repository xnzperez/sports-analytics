import axios from "axios";

// Definimos la interfaz para que TypeScript nos ayude (puedes omitirlo si usas JS)
export interface Match {
  ID: string; // GORM lo envía así por defecto
  home_team: string;
  away_team: string;
  home_odds: number;
  away_odds: number;
  sport_key: string;
  league: string;
  external_id: string;
  starts_at: string;
}

const API_URL =
  import.meta.env.VITE_API_URL ||
  "https://tu-app-en-azure.azurecontainerapps.io/api";

export const getAvailableMatches = async (): Promise<Match[]> => {
  try {
    // Ahora apuntará a Azure cuando lo subas a Vercel
    const response = await axios.get(`${API_URL}/markets`);

    // IMPORTANTE: Verifica si tu JSON de Go devuelve { "data": [...] }
    // o si devuelve el array directamente.
    // Si en Postman viste que los partidos están dentro de una llave "data", esto está bien:
    return response.data.data;
  } catch (error) {
    console.error("Error fetching matches:", error);
    return [];
  }
};
