import {
  X,
  Trophy,
  AlertCircle,
  Calendar,
  DollarSign,
  Target,
} from "lucide-react";
import { formatCurrency } from "../../../lib/utils"; // Asegúrate de tener esta ruta bien o usa tu util

interface Bet {
  ID: string;
  title: string;
  sport_key: string;
  status: "pending" | "WON" | "LOST";
  stake_units: number;
  odds: number;
  details: string; // Viene como string JSON
  created_at: string;
  ai_tip?: string; // Por si guardamos el tip
}

interface Props {
  bet: Bet | null;
  isOpen: boolean;
  onClose: () => void;
}

export const BetDetailsModal = ({ bet, isOpen, onClose }: Props) => {
  if (!isOpen || !bet) return null;

  // Intentamos parsear los detalles (que vienen como string JSON)
  let detailsObj: any = {};
  try {
    detailsObj = JSON.parse(bet.details);
  } catch (e) {
    console.error("Error parsing details", e);
  }

  // Cálculos
  const potentialWin = bet.stake_units * bet.odds;
  const profit = potentialWin - bet.stake_units;
  const isWinner = bet.status === "WON";
  const isLoser = bet.status === "LOST";

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/80 backdrop-blur-sm p-4 animate-in fade-in duration-200">
      <div className="w-full max-w-md bg-slate-900 border border-slate-800 rounded-xl shadow-2xl overflow-hidden relative">
        {/* Header con Color Dinámico */}
        <div
          className={`h-2 w-full ${
            isWinner
              ? "bg-emerald-500"
              : isLoser
              ? "bg-red-500"
              : "bg-yellow-500"
          }`}
        />

        <button
          onClick={onClose}
          className="absolute top-4 right-4 text-slate-400 hover:text-white transition-colors"
        >
          <X size={20} />
        </button>

        <div className="p-6 space-y-6">
          {/* Título y Estado */}
          <div className="text-center space-y-2">
            <div className="inline-flex items-center justify-center w-12 h-12 rounded-full bg-slate-800 border border-slate-700 mb-2">
              {isWinner ? (
                <Trophy className="text-emerald-500" size={24} />
              ) : isLoser ? (
                <AlertCircle className="text-red-500" size={24} />
              ) : (
                <Calendar className="text-yellow-500" size={24} />
              )}
            </div>
            <h2 className="text-xl font-bold text-white leading-tight">
              {bet.title}
            </h2>
            <p className="text-sm text-slate-400">
              {detailsObj.league || bet.sport_key.toUpperCase()}
            </p>

            <span
              className={`inline-block px-3 py-1 rounded-full text-xs font-bold border ${
                isWinner
                  ? "bg-emerald-500/10 border-emerald-500/20 text-emerald-400"
                  : isLoser
                  ? "bg-red-500/10 border-red-500/20 text-red-400"
                  : "bg-yellow-500/10 border-yellow-500/20 text-yellow-400"
              }`}
            >
              {bet.status === "WON"
                ? "GANADA"
                : bet.status === "LOST"
                ? "PERDIDA"
                : "PENDIENTE"}
            </span>
          </div>

          {/* Detalles de la Selección */}
          <div className="bg-slate-950/50 rounded-lg p-4 border border-slate-800 space-y-3">
            <div className="flex justify-between items-center text-sm">
              <span className="text-slate-400">Tu Selección:</span>
              <span className="font-bold text-white">
                {detailsObj.team_name || "Desconocido"}
              </span>
            </div>
            <div className="flex justify-between items-center text-sm">
              <span className="text-slate-400">Cuota (Odds):</span>
              <span className="font-mono text-yellow-400 font-bold">
                x{bet.odds.toFixed(2)}
              </span>
            </div>
          </div>

          {/* Finanzas */}
          <div className="grid grid-cols-2 gap-4">
            <div className="bg-slate-800/50 p-3 rounded-lg border border-slate-700">
              <div className="flex items-center gap-2 text-slate-400 text-xs mb-1">
                <DollarSign size={12} /> Apostado
              </div>
              <div className="text-lg font-bold text-white">
                {formatCurrency(bet.stake_units)}
              </div>
            </div>
            <div
              className={`p-3 rounded-lg border ${
                isWinner
                  ? "bg-emerald-900/20 border-emerald-500/30"
                  : "bg-slate-800/50 border-slate-700"
              }`}
            >
              <div className="flex items-center gap-2 text-slate-400 text-xs mb-1">
                <Target size={12} />{" "}
                {isWinner ? "Ganancia Neta" : "Retorno Potencial"}
              </div>
              <div
                className={`text-lg font-bold ${
                  isWinner ? "text-emerald-400" : "text-slate-200"
                }`}
              >
                {formatCurrency(isWinner ? profit : potentialWin)}
              </div>
            </div>
          </div>

          {/* Footer / ID */}
          <div className="pt-4 border-t border-slate-800 text-center">
            <p className="text-[10px] text-slate-600 font-mono">ID: {bet.ID}</p>
            <p className="text-[10px] text-slate-600">
              Fecha: {new Date(bet.created_at).toLocaleString()}
            </p>
          </div>
        </div>
      </div>
    </div>
  );
};
