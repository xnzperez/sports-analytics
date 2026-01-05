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

// VOLVAMOS A LO SEGURO
const API_URL =
  import.meta.env.VITE_API_URL ||
  "https://env-stakewise.victoriousflower-9df2d478.northcentralus.azurecontainerapps.io/api";

export const getAvailableMatches = async (): Promise<Match[]> => {
  try {
    const response = await axios.get(`${API_URL}/markets`);
    // Si antes te funcionaba con .data.data, déjalo así:
    return response.data.data;
  } catch (error) {
    console.error("Error fetching matches:", error);
    return [];
  }
};
