import { useEffect, useState } from "react";
import { getBets } from "../services/betting.service";
import type { Bet } from "../betting.schemas";
import { StatusBadge } from "../../../components/ui/StatusBadge";
import { Calendar, DollarSign, Trophy } from "lucide-react";

interface Props {
  refreshTrigger: number; // Un truco para recargar la lista cuando creamos una apuesta nueva
}

export const BetHistory = ({ refreshTrigger }: Props) => {
  const [bets, setBets] = useState<Bet[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const loadData = async () => {
      try {
        setLoading(true);
        const data = await getBets();
        setBets(data);
      } catch (error) {
        console.error("Error cargando historial", error);
      } finally {
        setLoading(false);
      }
    };
    loadData();
  }, [refreshTrigger]); // Cada vez que cambie este número, recargamos

  if (loading)
    return (
      <div className="text-slate-400 text-sm p-4">Cargando historial...</div>
    );

  if (bets.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center p-8 border border-dashed border-slate-800 rounded-xl bg-slate-900/50">
        <p className="text-slate-400 text-sm">
          Aún no tienes apuestas registradas.
        </p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <h3 className="text-xl font-bold text-white mb-4">Actividad Reciente</h3>
      <div className="grid grid-cols-1 gap-4">
        {bets.map((bet) => (
          <div
            key={bet.id}
            className="group relative bg-slate-800/50 border border-slate-700/50 hover:border-emerald-500/30 hover:bg-slate-800 transition-all rounded-xl p-4 overflow-hidden"
          >
            <div className="flex justify-between items-start mb-3">
              <div>
                <span className="text-xs font-mono text-slate-500 uppercase tracking-wider block mb-1">
                  {bet.sport_key}
                </span>
                <h4 className="font-semibold text-slate-200">{bet.title}</h4>
              </div>
              <StatusBadge status={bet.status} />
            </div>

            <div className="grid grid-cols-3 gap-4 text-sm border-t border-slate-700/50 pt-3 mt-1">
              {/* Cuota */}
              <div className="flex flex-col">
                <span className="text-slate-500 text-xs mb-1">Cuota</span>
                <div className="flex items-center gap-1 text-slate-300 font-medium">
                  <span className="bg-slate-700 px-1.5 rounded text-xs">x</span>
                  {bet.odds.toFixed(2)}
                </div>
              </div>

              {/* Apostado */}
              <div className="flex flex-col">
                <span className="text-slate-500 text-xs mb-1">Apostado</span>
                <div className="flex items-center gap-1 text-slate-300 font-medium">
                  <DollarSign size={12} className="text-slate-500" />
                  {bet.stake_units.toFixed(2)}
                </div>
              </div>

              {/* Retorno Potencial */}
              <div className="flex flex-col">
                <span className="text-slate-500 text-xs mb-1">
                  Ganancia Posible
                </span>
                <div className="flex items-center gap-1 text-emerald-400 font-bold">
                  <Trophy size={12} />
                  {(bet.stake_units * bet.odds).toFixed(2)}
                </div>
              </div>
            </div>

            {/* Fecha absoluta pequeña abajo a la derecha */}
            <div className="absolute bottom-2 right-4 flex items-center gap-1 text-[10px] text-slate-600">
              <Calendar size={10} />
              {new Date(bet.created_at).toLocaleDateString()}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};
