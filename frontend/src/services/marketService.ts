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
  "https://env-stakewise.victoriousflower-9df2d478.northcentralus.azurecontainerapps.io";

export const getAvailableMatches = async () => {
  try {
    // 2. Aquí le pegamos directo a la ruta que ya confirmaste que funciona en el navegador
    const response = await axios.get(`${API_URL}/api/markets`);

    // 3. Retornamos la data (asegurando que sea el array)
    return response.data.data || response.data;
  } catch (error) {
    console.error("Error cargando partidos:", error);
    return [];
  }
};
