import { api } from "../../../lib/axios";
import type { CreateBetFormData, Bet } from "../betting.schemas"; // <--- Agrega Bet
// Mapeo simple para enviar al backend lo que espera
interface PlaceBetRequest {
  title: string;
  stake_units: number;
  odds: number;
  sport_key: string;
  is_parlay: boolean;
  user_notes?: string;
  match_date: string; // Enviaremos la fecha actual por ahora
  details: object; // Detalles vacíos por ahora
}

export const placeBet = async (data: CreateBetFormData) => {
  const payload: PlaceBetRequest = {
    ...data,
    match_date: new Date().toISOString(),
    details: {}, // En el futuro aquí irán equipos, ligas, etc.
  };

  const response = await api.post("/api/bets", payload);
  return response.data;
};

export const getBets = async () => {
  // Pedimos las apuestas ordenadas (el backend ya lo hace)
  const response = await api.get<{ data: Bet[] }>("/api/bets?limit=20");
  return response.data.data;
};
