import { useEffect, useState } from "react";
// Asegúrate de que esta ruta sea correcta según tu estructura
import api from "../../../lib/axios";
import { formatCurrency } from "../../../lib/utils";
import { BetDetailsModal } from "./BetDetailsModal"; // <--- Asegúrate de tener este componente creado
import { Calendar, Trophy, AlertCircle, Loader2 } from "lucide-react";

// Definición local de la interfaz para asegurar tipos (o impórtala de tus schemas)
interface Bet {
  ID: string; // Ojo: Go devuelve 'ID' (mayúscula) o 'id' según tu JSON tag. Verifica esto.
  id?: string; // Fallback por si acaso
  title: string;
  sport_key: string;
  status: "pending" | "WON" | "LOST";
  stake_units: number;
  odds: number;
  details: string;
  created_at: string;
}

interface Props {
  refreshTrigger: number;
}

export const BetHistory = ({ refreshTrigger }: Props) => {
  const [bets, setBets] = useState<Bet[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);
  const [selectedBet, setSelectedBet] = useState<Bet | null>(null);

  useEffect(() => {
    const fetchBets = async () => {
      try {
        setLoading(true);
        setError(false);

        const response = await api.get("/api/bets");
        const rawData = response.data;

        // --- CORRECCIÓN AQUÍ ---
        // Verificamos si la respuesta es un array directo O si viene dentro de .data
        let betsArray: Bet[] = [];

        if (Array.isArray(rawData)) {
          betsArray = rawData;
        } else if (rawData && Array.isArray(rawData.data)) {
          betsArray = rawData.data; // <--- Aquí es donde estaba fallando antes
        } else {
          console.warn("Formato de respuesta desconocido:", rawData);
        }

        setBets(betsArray);
      } catch (err) {
        console.error("Error fetching bets:", err);
        setError(true);
        setBets([]);
      } finally {
        setLoading(false);
      }
    };

    fetchBets();
  }, [refreshTrigger]);

  if (loading) {
    return (
      <div className="flex items-center justify-center p-8 text-slate-400 gap-2">
        <Loader2 className="animate-spin" size={20} />
        Cargando historial...
      </div>
    );
  }

  if (error) {
    return (
      <div className="p-4 border border-red-900/50 bg-red-900/20 rounded-lg text-red-400 text-sm text-center">
        No se pudo conectar con el servidor. Verifica que el backend esté
        corriendo.
      </div>
    );
  }

  if (bets.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center p-8 border border-dashed border-slate-800 rounded-xl bg-slate-900/50">
        <Trophy className="text-slate-600 mb-2" size={32} />
        <p className="text-slate-400 text-sm">
          Aún no tienes apuestas registradas.
        </p>
      </div>
    );
  }

  return (
    <>
      <div className="space-y-4">
        <h3 className="text-xl font-bold text-white mb-4 flex items-center gap-2">
          <Calendar size={20} className="text-emerald-500" />
          Historial de Apuestas
        </h3>

        <div className="grid grid-cols-1 gap-3">
          {bets.map((bet) => {
            // Normalizar ID (por si Go manda ID o id)
            const uniqueKey = bet.ID || bet.id || Math.random().toString();
            const isWin = bet.status === "WON";
            const isLoss = bet.status === "LOST";

            return (
              <div
                key={uniqueKey}
                onClick={() => setSelectedBet(bet)}
                className={`
                  group relative flex items-center justify-between p-4 rounded-xl border cursor-pointer transition-all
                  ${
                    isWin
                      ? "bg-emerald-950/10 border-emerald-500/20 hover:bg-emerald-900/20"
                      : isLoss
                      ? "bg-red-950/10 border-red-500/20 hover:bg-red-900/20"
                      : "bg-slate-800/50 border-slate-700 hover:border-slate-600 hover:bg-slate-800"
                  }
                `}
              >
                {/* Lado Izquierdo: Info Principal */}
                <div className="flex flex-col">
                  <span className="text-[10px] font-bold text-slate-500 uppercase tracking-wider mb-0.5">
                    {bet.sport_key}
                  </span>
                  <h4 className="font-semibold text-slate-200 group-hover:text-white transition-colors text-sm md:text-base">
                    {bet.title}
                  </h4>
                  <div className="flex items-center gap-2 mt-1">
                    <span
                      className={`text-xs px-2 py-0.5 rounded-full font-bold border ${
                        isWin
                          ? "bg-emerald-500/10 border-emerald-500/20 text-emerald-400"
                          : isLoss
                          ? "bg-red-500/10 border-red-500/20 text-red-400"
                          : "bg-yellow-500/10 border-yellow-500/20 text-yellow-400"
                      }`}
                    >
                      {bet.status}
                    </span>
                    <span className="text-xs text-slate-500">
                      {new Date(bet.created_at).toLocaleDateString()}
                    </span>
                  </div>
                </div>

                {/* Lado Derecho: Dinero */}
                <div className="text-right">
                  <div className="text-xs text-slate-500 mb-1">
                    Cuota:{" "}
                    <span className="text-yellow-500 font-mono">
                      {bet.odds.toFixed(2)}
                    </span>
                  </div>
                  <div
                    className={`font-bold text-base md:text-lg ${
                      isWin
                        ? "text-emerald-400"
                        : isLoss
                        ? "text-red-400"
                        : "text-white"
                    }`}
                  >
                    {isWin
                      ? `+${formatCurrency(
                          bet.stake_units * bet.odds - bet.stake_units
                        )}`
                      : isLoss
                      ? `-${formatCurrency(bet.stake_units)}`
                      : formatCurrency(bet.stake_units)}
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      </div>

      {/* MODAL DE DETALLES */}
      <BetDetailsModal
        isOpen={!!selectedBet}
        bet={selectedBet}
        onClose={() => setSelectedBet(null)}
      />
    </>
  );
};
