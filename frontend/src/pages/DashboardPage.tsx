import { useEffect, useState } from "react";
import {
  Plus,
  TrendingUp,
  DollarSign,
  Activity,
  Target,
  Filter,
} from "lucide-react"; // Agregamos Filter
import { useAuthStore } from "../features/auth/auth.store";
import { Button } from "../components/ui/Button";
import { PlaceBetModal } from "../features/betting/components/PlaceBetModal";
import { BetHistory } from "../features/betting/components/BetHistory";
import { ProfitBySportChart } from "../features/analytics/components/DashboardCharts";
import api from "../lib/axios";
import { formatCurrency } from "../lib/utils";

interface DashboardStats {
  total_bets: number;
  won_bets: number;
  win_rate: number;
  total_profit: number;
  current_bankroll: number;
  ai_tip: string;
  sport_performance: any[];
}

export const DashboardPage = () => {
  const { logout, user, fetchUser } = useAuthStore();
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [refreshTrigger, setRefreshTrigger] = useState(0);
  const [stats, setStats] = useState<DashboardStats | null>(null);

  // --- NUEVO ESTADO PARA FILTROS ---
  const [sportFilter, setSportFilter] = useState<string>("");

  const loadStats = async () => {
    try {
      // Enviamos el filtro como Query Parameter
      const query = sportFilter ? `?sport=${sportFilter}` : "";
      const { data } = await api.get(`/api/stats${query}`);
      setStats(data);
    } catch (error) {
      console.error("Error cargando estadísticas:", error);
    }
  };

  useEffect(() => {
    fetchUser();
    loadStats();
    // Se recarga cuando cambia el trigger O el filtro de deporte
  }, [refreshTrigger, sportFilter]);

  const handleBetSuccess = () => {
    setIsModalOpen(false);
    setRefreshTrigger((prev) => prev + 1);
  };

  return (
    <div className="min-h-screen bg-slate-950 text-white p-4 md:p-8 font-sans">
      <div className="max-w-7xl mx-auto space-y-8">
        {/* --- HEADER --- */}
        <header className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
          <div>
            <h1 className="text-3xl font-bold bg-gradient-to-r from-emerald-400 to-cyan-400 bg-clip-text text-transparent">
              Dashboard de Estrategia
            </h1>
            <p className="text-slate-400">Análisis para {user?.username}</p>
          </div>

          <div className="flex gap-3 w-full md:w-auto">
            <Button
              onClick={() => setIsModalOpen(true)}
              className="gap-2 bg-emerald-600"
            >
              <Plus size={18} /> Nueva Apuesta
            </Button>
            <Button
              variant="outline"
              onClick={logout}
              className="border-slate-700"
            >
              Salir
            </Button>
          </div>
        </header>

        {/* --- BARRA DE FILTROS --- */}
        <div className="flex flex-wrap items-center justify-between gap-4 bg-slate-900/40 p-4 rounded-xl border border-slate-800">
          <div className="flex items-center gap-2 text-slate-300">
            <Filter size={18} className="text-emerald-500" />
            <span className="text-sm font-medium">Filtrar rendimiento:</span>
          </div>
          <div className="flex gap-2">
            {["", "lol", "valorant", "cs2", "football", "basketball"].map(
              (s) => (
                <button
                  key={s}
                  onClick={() => setSportFilter(s)}
                  className={`px-4 py-1.5 rounded-full text-xs font-bold transition-all ${
                    sportFilter === s
                      ? "bg-emerald-500 text-white shadow-lg shadow-emerald-500/20"
                      : "bg-slate-800 text-slate-400 hover:bg-slate-700"
                  }`}
                >
                  {s === "" ? "TODOS" : s.toUpperCase()}
                </button>
              )
            )}
          </div>
        </div>

        {/* --- KPI CARDS --- */}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          <div className="bg-slate-900/50 p-5 rounded-xl border border-slate-800 group hover:border-emerald-500/30 transition-all">
            <p className="text-slate-400 text-sm flex items-center gap-2">
              <DollarSign size={14} /> Bankroll
            </p>
            <p className="text-3xl font-bold text-white mt-1">
              {formatCurrency(user?.bankroll || 0)}
            </p>
          </div>

          <div className="bg-slate-900/50 p-5 rounded-xl border border-slate-800 group hover:border-blue-500/30 transition-all">
            <p className="text-slate-400 text-sm flex items-center gap-2">
              <TrendingUp size={14} /> Profit Neto
            </p>
            <p
              className={`text-3xl font-bold mt-1 ${
                stats?.total_profit && stats.total_profit >= 0
                  ? "text-emerald-400"
                  : "text-red-400"
              }`}
            >
              {stats ? formatCurrency(stats.total_profit) : "$0.00"}
            </p>
          </div>

          <div className="bg-slate-900/50 p-5 rounded-xl border border-slate-800 group hover:border-purple-500/30 transition-all">
            <p className="text-slate-400 text-sm flex items-center gap-2">
              <Target size={14} /> Win Rate
            </p>
            <p className="text-3xl font-bold text-white mt-1">
              {stats?.win_rate ? stats.win_rate.toFixed(1) : "0.0"}%
            </p>
          </div>

          <div className="bg-slate-900/50 p-5 rounded-xl border border-slate-800 group hover:border-yellow-500/30 transition-all">
            <p className="text-slate-400 text-sm flex items-center gap-2">
              <Activity size={14} /> Total Jugadas
            </p>
            <p className="text-3xl font-bold text-white mt-1">
              {stats?.total_bets || 0}
            </p>
          </div>
        </div>

        {/* --- ANALYTICS & TIPS --- */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          <div className="space-y-6">
            <div className="bg-indigo-900/20 border border-indigo-500/30 p-5 rounded-xl shadow-inner">
              <h4 className="text-indigo-400 text-xs font-bold uppercase tracking-wider mb-2">
                Tip de Inteligencia Artificial
              </h4>
              <p className="text-indigo-100 text-sm leading-relaxed italic">
                "
                {stats?.ai_tip ||
                  "Recopilando datos suficientes para asesorarte..."}
                "
              </p>
            </div>
            <ProfitBySportChart data={stats?.sport_performance || []} />
          </div>

          <div className="lg:col-span-2 space-y-4">
            <h3 className="text-lg font-bold flex items-center gap-2">
              <TrendingUp className="text-blue-500" size={20} /> Historial
              Reciente
            </h3>
            <div className="bg-slate-900/50 border border-slate-800 rounded-xl overflow-hidden">
              <BetHistory refreshTrigger={refreshTrigger} />
            </div>
          </div>
        </div>

        <PlaceBetModal
          isOpen={isModalOpen}
          onClose={() => setIsModalOpen(false)}
          onSuccess={handleBetSuccess}
        />
      </div>
    </div>
  );
};
