import { z } from "zod";

export const createBetSchema = (currentBankroll: number) =>
  z.object({
    title: z.string().min(3, "El título es requerido"),

    // Usamos z.number() directo con un pre-procesador para asegurar conversión
    stake_units: z.preprocess(
      (val) => Number(val),
      z
        .number()
        .min(1, "Apuesta mínima $1")
        .max(currentBankroll, "Fondos insuficientes")
    ),

    odds: z.preprocess(
      (val) => Number(val),
      z.number().min(1.01, "Cuota mínima 1.01")
    ),

    sport_key: z.string().min(1),
    is_parlay: z.boolean().default(false),
    user_notes: z.string().optional(),

    // Details opcional y permisivo para evitar conflicto de tipos
    details: z.record(z.any()).optional(),
  });

// Definimos el tipo manualmente para que coincida EXACTO con lo que espera el formulario
// Esto rompe el ciclo de inferencia errónea
export type CreateBetFormData = {
  title: string;
  stake_units: number;
  odds: number;
  sport_key: string;
  is_parlay: boolean;
  user_notes?: string;
  details?: Record<string, any>;
};

// Interfaz para mostrar apuestas (lectura)
export interface Bet {
  id: string;
  title: string;
  sport_key: string;
  stake_units: number;
  odds: number;
  status: "pending" | "won" | "lost" | "void";
  created_at: string;
  details?: string;
}
