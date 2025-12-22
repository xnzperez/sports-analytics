import axios from "axios";

// Definimos la interfaz para que TypeScript nos ayude (puedes omitirlo si usas JS)
export interface Match {
  ID: string; // Nuestro UUID interno
  external_id: string;
  league: string;
  home_team: string;
  away_team: string;
  starts_at: string;
  home_odds: number;
  away_odds: number;
  sport_key: string;
}

const API_URL = "http://localhost:3000/api"; // O tu variable de entorno

export const getAvailableMatches = async (): Promise<Match[]> => {
  try {
    // Llamamos a la ruta pública que acabamos de habilitar
    const response = await axios.get(`${API_URL}/markets`);
    return response.data.data; // Accedemos al array dentro de "data"
  } catch (error) {
    console.error("Error fetching matches:", error);
    return [];
  }
};
