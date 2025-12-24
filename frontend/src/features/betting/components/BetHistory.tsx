import { useEffect, useState } from "react";
import api from "../../../lib/axios";
import { formatCurrency } from "../../../lib/utils";
import { BetDetailsModal } from "./BetDetailsModal"; // <--- IMPORTAR

// ... (Interface Bet que ya tenías, asegúrate que coincida)

export const BetHistory = ({ refreshTrigger }: { refreshTrigger: number }) => {
  const [bets, setBets] = useState<any[]>([]); // Usa 'any' o tu interfaz Bet
  const [loading, setLoading] = useState(true);

  // --- NUEVO ESTADO PARA EL MODAL ---
  const [selectedBet, setSelectedBet] = useState<any | null>(null);

  useEffect(() => {
    const fetchBets = async () => {
      try {
        const { data } = await api.get("/api/bets");
        setBets(data);
      } catch (error) {
        console.error("Error fetching bets", error);
      } finally {
        setLoading(false);
      }
    };
    fetchBets();
  }, [refreshTrigger]);

  if (loading)
    return <div className="p-4 text-slate-400">Cargando historial...</div>;

  return (
    <>
      <div className="overflow-x-auto">
        <table className="w-full text-sm text-left">
          <thead className="text-xs text-slate-400 uppercase bg-slate-900/50 border-b border-slate-800">
            <tr>
              <th className="px-6 py-3">Evento</th>
              <th className="px-6 py-3">Selección</th>
              <th className="px-6 py-3">Stake</th>
              <th className="px-6 py-3">Cuota</th>
              <th className="px-6 py-3">Estado</th>
              <th className="px-6 py-3">Fecha</th>
            </tr>
          </thead>
          <tbody>
            {bets.length === 0 ? (
              <tr>
                <td
                  colSpan={6}
                  className="px-6 py-8 text-center text-slate-500"
                >
                  No hay apuestas registradas aún.
                </td>
              </tr>
            ) : (
              bets.map((bet) => (
                <tr
                  key={bet.ID}
                  // --- CLICK PARA ABRIR MODAL ---
                  onClick={() => setSelectedBet(bet)}
                  className="bg-transparent border-b border-slate-800 hover:bg-slate-800/50 transition-colors cursor-pointer group"
                >
                  <td className="px-6 py-4 font-medium text-white">
                    {bet.title}
                    <span className="block text-xs text-slate-500 font-normal mt-0.5">
                      {bet.sport_key.toUpperCase()}
                    </span>
                  </td>
                  <td className="px-6 py-4">
                    {/* Intentamos mostrar el equipo seleccionado de los details */}
                    {(() => {
                      try {
                        const d = JSON.parse(bet.details);
                        return (
                          <span className="text-indigo-300">{d.team_name}</span>
                        );
                      } catch {
                        return "-";
                      }
                    })()}
                  </td>
                  <td className="px-6 py-4">
                    {formatCurrency(bet.stake_units)}
                  </td>
                  <td className="px-6 py-4 text-yellow-500 font-bold">
                    {bet.odds}
                  </td>
                  <td className="px-6 py-4">
                    <span
                      className={`px-2 py-1 rounded text-xs font-bold ${
                        bet.status === "WON"
                          ? "bg-emerald-500/10 text-emerald-400"
                          : bet.status === "LOST"
                          ? "bg-red-500/10 text-red-400"
                          : "bg-yellow-500/10 text-yellow-400"
                      }`}
                    >
                      {bet.status === "WON"
                        ? "GANADA"
                        : bet.status === "LOST"
                        ? "PERDIDA"
                        : "PENDIENTE"}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-slate-400">
                    {new Date(bet.created_at).toLocaleDateString()}
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      {/* --- RENDERIZAR EL MODAL --- */}
      <BetDetailsModal
        isOpen={!!selectedBet}
        bet={selectedBet}
        onClose={() => setSelectedBet(null)}
      />
    </>
  );
};
