import { useEffect, useState } from "react";
import { Plus } from "lucide-react";
import { useAuthStore } from "../features/auth/auth.store";
import { Button } from "../components/ui/Button";
import { PlaceBetModal } from "../features/betting/components/PlaceBetModal";
import { BetHistory } from "../features/betting/components/BetHistory"; // <--- Importar

export const DashboardPage = () => {
  const { logout, user, fetchUser } = useAuthStore();
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [refreshTrigger, setRefreshTrigger] = useState(0); // <--- Estado para recargar lista

  useEffect(() => {
    fetchUser();
  }, [refreshTrigger]); // Recargar saldo también cuando cambie el trigger

  const handleBetSuccess = () => {
    setIsModalOpen(false);
    setRefreshTrigger((prev) => prev + 1); // <--- Incrementamos para avisar que recarguen
  };

  return (
    <div className="min-h-screen bg-slate-900 text-white p-8">
      <div className="max-w-7xl mx-auto">
        {/* Header (Igual que antes) */}
        <header className="flex justify-between items-center mb-10">
          <div>
            <h1 className="text-3xl font-bold text-emerald-400">Dashboard</h1>
            <p className="text-slate-400">
              Bienvenido,{" "}
              <span className="text-white font-medium">
                {user?.username || "Usuario"}
              </span>
            </p>
          </div>

          <div className="flex gap-4">
            <Button onClick={() => setIsModalOpen(true)} className="gap-2">
              <Plus size={18} />
              Nueva Apuesta
            </Button>
            <Button variant="outline" onClick={logout} className="w-32">
              Cerrar Sesión
            </Button>
          </div>
        </header>

        {/* KPIs (Tarjetas de dinero - Igual que antes) */}
        <main className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
          {/* ... Tus tarjetas de Bankroll y ROI van aquí ... */}
          <div className="bg-slate-800 p-6 rounded-xl border border-slate-700 h-40 flex flex-col justify-center relative overflow-hidden">
            <div className="absolute top-4 right-4 text-emerald-500/20 text-6xl font-bold select-none">
              $
            </div>
            <h3 className="text-slate-400 font-medium z-10">Bankroll Actual</h3>
            <p className="text-4xl font-bold mt-2 text-white z-10">
              {user?.bankroll
                ? new Intl.NumberFormat("en-US", {
                    style: "currency",
                    currency: "USD",
                  }).format(user.bankroll)
                : "$0.00"}
            </p>
          </div>
          {/* Si borraste las otras tarjetas por error, avísame para pasártelas de nuevo, 
                pero asumo que siguen ahí */}
        </main>

        {/* --- NUEVO: Sección de Historial --- */}
        <section className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Columna Izquierda: Historial de Apuestas (Ocupa 2/3 en pantallas grandes) */}
          <div className="lg:col-span-2">
            <BetHistory refreshTrigger={refreshTrigger} />
          </div>

          {/* Columna Derecha: Próximamente Analytics o Gráficas */}
          <div className="bg-slate-800/30 border border-slate-700/50 rounded-xl p-6 h-fit">
            <h3 className="text-lg font-bold text-white mb-2">
              Resumen Rápido
            </h3>
            <p className="text-slate-400 text-sm">
              Aquí verás tus gráficas de rendimiento pronto.
            </p>
            {/* Placeholder para gráficas futuras */}
            <div className="h-40 mt-4 bg-slate-800/50 rounded flex items-center justify-center text-xs text-slate-600">
              Coming Soon: Winrate Chart
            </div>
          </div>
        </section>

        <PlaceBetModal
          isOpen={isModalOpen}
          onClose={() => setIsModalOpen(false)}
          onSuccess={handleBetSuccess} // <--- Usamos la nueva función
        />
      </div>
    </div>
  );
};
