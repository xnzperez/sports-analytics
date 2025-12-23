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

    // CAMBIO 1: Usamos string para aceptar 'lol', 'dota2', 'cs2' dinámicamente
    sport_key: z.string().min(1, "El deporte es requerido"),

    is_parlay: z.boolean().default(false),
    user_notes: z.string().optional(),

    // CAMBIO 2: Agregamos details para soportar la integración con la API de Partidos
    // Recibirá un JSON string con { match_id, selection, external_id... }
    details: z.any().optional(),
  });

export type CreateBetFormData = z.infer<ReturnType<typeof createBetSchema>>;

// Interfaz de la Apuesta que viene del Backend (Para mostrar en el historial)
export interface Bet {
  id: string;
  title: string;
  sport_key: string;
  stake_units: number;
  odds: number;
  status: "pending" | "won" | "lost" | "void"; // Minúsculas para coincidir con backend
  potential_payout?: number;
  created_at: string;
  // match_date: string; // A veces no viene si es manual, cuidado aquí
  details?: string; // Para poder leer la metadata en el futuro
}
