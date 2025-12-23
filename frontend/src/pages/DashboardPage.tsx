import { useEffect, useState } from "react";
import { Plus, TrendingUp, DollarSign, Activity, Target } from "lucide-react";
import { useAuthStore } from "../features/auth/auth.store";
import { Button } from "../components/ui/Button";
import { PlaceBetModal } from "../features/betting/components/PlaceBetModal";
import { BetHistory } from "../features/betting/components/BetHistory";
import { ProfitBySportChart } from "../features/analytics/components/DashboardCharts";
import api from "../lib/axios";
import { formatCurrency } from "../lib/utils";

// Definimos la interfaz de los datos que vienen del Backend
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

  // 1. Definimos la función de carga PRIMERO
  const loadStats = async () => {
    try {
      // Usamos la ruta correcta con /api
      const { data } = await api.get("/api/stats");
      setStats(data);
    } catch (error) {
      console.error("Error cargando estadísticas:", error);
    }
  };

  // 2. Luego usamos el useEffect que llama a esa función
  useEffect(() => {
    fetchUser();
    loadStats();
  }, [refreshTrigger]); // Se ejecutará al inicio y cuando cambie el trigger

  const handleBetSuccess = () => {
    setIsModalOpen(false);
    // Esto forzará que el useEffect se ejecute de nuevo, recargando saldo y gráficas
    setRefreshTrigger((prev) => prev + 1);
  };

  return (
    <div className="min-h-screen bg-slate-950 text-white p-4 md:p-8 font-sans">
      <div className="max-w-7xl mx-auto space-y-8">
        {/* --- HEADER --- */}
        <header className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
          <div>
            <h1 className="text-3xl font-bold bg-gradient-to-r from-emerald-400 to-cyan-400 bg-clip-text text-transparent">
              Panel de Control
            </h1>
            <p className="text-slate-400">
              Estrategia y rendimiento de{" "}
              <span className="text-white font-medium">{user?.username}</span>
            </p>
          </div>

          <div className="flex gap-3 w-full md:w-auto">
            <Button
              onClick={() => setIsModalOpen(true)}
              className="gap-2 flex-1 md:flex-none bg-emerald-600 hover:bg-emerald-700"
            >
              <Plus size={18} /> Nueva Apuesta
            </Button>
            <Button
              variant="outline"
              onClick={logout}
              className="border-slate-700 hover:bg-slate-800"
            >
              Salir
            </Button>
          </div>
        </header>

        {/* --- SECCIÓN KPI (Tarjetas de Métricas) --- */}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          {/* 1. Bankroll Actual */}
          <div className="bg-slate-900/50 p-5 rounded-xl border border-slate-800 backdrop-blur-sm relative overflow-hidden group hover:border-emerald-500/30 transition-all">
            <div className="absolute right-0 top-0 p-4 opacity-10 group-hover:opacity-20 transition-opacity">
              <DollarSign size={48} className="text-emerald-400" />
            </div>
            <p className="text-slate-400 text-sm font-medium flex items-center gap-2">
              <DollarSign size={14} /> Bankroll
            </p>
            <p className="text-3xl font-bold text-white mt-1">
              {formatCurrency(user?.bankroll || 0)}
            </p>
            <div className="text-xs text-emerald-400 mt-2 font-medium">
              Disponible para apostar
            </div>
          </div>

          {/* 2. Profit Total (Pérdidas/Ganancias) */}
          <div className="bg-slate-900/50 p-5 rounded-xl border border-slate-800 backdrop-blur-sm relative overflow-hidden group hover:border-blue-500/30 transition-all">
            <div className="absolute right-0 top-0 p-4 opacity-10 group-hover:opacity-20 transition-opacity">
              <TrendingUp
                size={48}
                className={
                  stats?.total_profit && stats.total_profit >= 0
                    ? "text-blue-400"
                    : "text-red-400"
                }
              />
            </div>
            <p className="text-slate-400 text-sm font-medium flex items-center gap-2">
              <TrendingUp size={14} /> Profit Neto
            </p>
            <p
              className={`text-3xl font-bold mt-1 ${
                stats?.total_profit && stats.total_profit >= 0
                  ? "text-blue-400"
                  : "text-red-400"
              }`}
            >
              {stats ? formatCurrency(stats.total_profit) : "$0.00"}
            </p>
            <div className="text-xs text-slate-500 mt-2">Histórico total</div>
          </div>

          {/* 3. Win Rate */}
          <div className="bg-slate-900/50 p-5 rounded-xl border border-slate-800 backdrop-blur-sm relative overflow-hidden group hover:border-purple-500/30 transition-all">
            <div className="absolute right-0 top-0 p-4 opacity-10 group-hover:opacity-20 transition-opacity">
              <Target size={48} className="text-purple-400" />
            </div>
            <p className="text-slate-400 text-sm font-medium flex items-center gap-2">
              <Target size={14} /> Win Rate
            </p>
            <p className="text-3xl font-bold text-white mt-1">
              {stats?.win_rate ? stats.win_rate.toFixed(1) : "0.0"}%
            </p>
            <div className="text-xs text-purple-400 mt-2">
              {stats?.won_bets || 0} Aciertos
            </div>
          </div>

          {/* 4. Total Apuestas */}
          <div className="bg-slate-900/50 p-5 rounded-xl border border-slate-800 backdrop-blur-sm relative overflow-hidden group hover:border-yellow-500/30 transition-all">
            <div className="absolute right-0 top-0 p-4 opacity-10 group-hover:opacity-20 transition-opacity">
              <Activity size={48} className="text-yellow-400" />
            </div>
            <p className="text-slate-400 text-sm font-medium flex items-center gap-2">
              <Activity size={14} /> Total Jugadas
            </p>
            <p className="text-3xl font-bold text-white mt-1">
              {stats?.total_bets || 0}
            </p>
            <div className="text-xs text-yellow-400 mt-2">
              Registradas en sistema
            </div>
          </div>
        </div>

        {/* --- CONTENIDO PRINCIPAL --- */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Columna Izquierda: Gráficas (Analytics) */}
          <div className="space-y-6">
            <h3 className="text-lg font-bold text-white flex items-center gap-2">
              <Activity className="text-emerald-500" size={20} />
              Analytics
            </h3>

            {/* Gráfica de Barras */}
            <ProfitBySportChart data={stats?.sport_performance || []} />

            {/* Aquí podrías poner otra gráfica pequeña o un tip de IA */}
            <div className="bg-indigo-900/20 border border-indigo-500/30 p-4 rounded-lg">
              <p className="text-indigo-200 text-sm">
                {stats?.ai_tip || "Analizando tus datos para darte consejos..."}
              </p>
            </div>
          </div>

          {/* Columna Derecha: Historial (Ocupa 2 espacios en pantallas grandes) */}
          <div className="lg:col-span-2 space-y-6">
            <h3 className="text-lg font-bold text-white flex items-center gap-2">
              <TrendingUp className="text-blue-500" size={20} />
              Historial Reciente
            </h3>
            <div className="bg-slate-900/50 border border-slate-800 rounded-xl overflow-hidden shadow-sm">
              <BetHistory refreshTrigger={refreshTrigger} />
            </div>
          </div>
        </div>

        {/* Modal */}
        <PlaceBetModal
          isOpen={isModalOpen}
          onClose={() => setIsModalOpen(false)}
          onSuccess={handleBetSuccess}
        />
      </div>
    </div>
  );
};
