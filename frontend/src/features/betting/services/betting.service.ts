import { api } from "../../../lib/axios";
import type { CreateBetFormData, Bet } from "../betting.schemas";

// Mapeo simple para enviar al backend lo que espera
interface PlaceBetRequest {
  title: string;
  stake_units: number;
  odds: number;
  sport_key: string;
  is_parlay: boolean;
  user_notes?: string;
  match_date: string;
  details?: any; // CAMBIO: Permitimos que venga cualquier cosa o sea opcional
}

export const placeBet = async (data: CreateBetFormData) => {
  const payload: PlaceBetRequest = {
    ...data,
    match_date: new Date().toISOString(),
    // details: {},  <--- ¡BORRA ESTA LÍNEA ASESINA!
  };

  const response = await api.post("/api/bets", payload);
  return response.data;
};

export const getBets = async () => {
  const response = await api.get<{ data: Bet[] }>("/api/bets?limit=20");
  return response.data.data;
};
