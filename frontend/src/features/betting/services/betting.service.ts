import { api } from "../../../lib/axios";
import type { CreateBetFormData } from "../betting.schemas";

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
