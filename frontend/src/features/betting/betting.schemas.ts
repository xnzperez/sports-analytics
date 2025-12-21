import { z } from "zod";

export const createBetSchema = (currentBankroll: number) =>
  z.object({
    title: z
      .string()
      .min(3, "El título debe ser descriptivo (ej: 'Lakers vs Warriors')"),
    stake_units: z.coerce // coerce convierte el string del input a número automáticamente
      .number()
      .min(1, "La apuesta mínima es $1")
      .max(currentBankroll, "No tienes suficientes fondos"),
    odds: z.coerce.number().min(1.01, "La cuota mínima es 1.01"),
    sport_key: z.enum(["football", "basketball", "esports", "tennis", "other"]),
    is_parlay: z.boolean().default(false),
    user_notes: z.string().optional(),
  });

export type CreateBetFormData = z.infer<ReturnType<typeof createBetSchema>>;

// Interfaz de la Apuesta que viene del Backend
export interface Bet {
  id: string;
  title: string;
  sport_key: string;
  stake_units: number;
  odds: number;
  status: "pending" | "WON" | "LOST" | "VOID";
  potential_payout?: number; // Calculado en el front (stake * odds)
  created_at: string;
  match_date: string;
}
