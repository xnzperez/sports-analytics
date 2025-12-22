import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

/**
 * Función de utilidad para mezclar clases de Tailwind.
 */
export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

/**
 * Formatea una fecha a un formato legible en español.
 * Ej: "22 dic, 14:30"
 */
export function formatDate(dateString: string | Date | undefined): string {
  if (!dateString) return "";

  const date = new Date(dateString);

  // Si la fecha no es válida, retornamos string vacío
  if (isNaN(date.getTime())) return "";

  return new Intl.DateTimeFormat("es-CO", {
    day: "numeric",
    month: "short",
    hour: "2-digit",
    minute: "2-digit",
    hour12: true,
  }).format(date);
}

/**
 * Formatea dinero a formato moneda (USD/COP).
 * Ej: $1,500.00
 */
export function formatCurrency(amount: number): string {
  return new Intl.NumberFormat("en-US", {
    style: "currency",
    currency: "USD",
    minimumFractionDigits: 2,
  }).format(amount);
}
