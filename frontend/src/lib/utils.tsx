import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

/**
 * Función de utilidad para mezclar clases de Tailwind condicionalmente.
 * Permite hacer cosas como: cn("bg-red-500", condition && "bg-green-500")
 * y resuelve conflictos automáticamente.
 */
export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}
