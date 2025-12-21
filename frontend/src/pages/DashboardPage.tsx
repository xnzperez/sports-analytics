import { useEffect, useState } from "react"; // <--- Importa useState
import { Plus } from "lucide-react"; // <--- Importa icono Plus
import { useAuthStore } from "../features/auth/auth.store";
import { Button } from "../components/ui/Button";
import { PlaceBetModal } from "../features/betting/components/PlaceBetModal"; // <--- Importa Modal

export const DashboardPage = () => {
  const { logout, user, fetchUser } = useAuthStore();
  const [isModalOpen, setIsModalOpen] = useState(false); // <--- Estado del modal

  useEffect(() => {
    fetchUser();
  }, []);

  return (
    <div className="min-h-screen bg-slate-900 text-white p-8">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
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
            {/* Botón Nueva Apuesta */}
            <Button onClick={() => setIsModalOpen(true)} className="gap-2">
              <Plus size={18} />
              Nueva Apuesta
            </Button>

            <Button variant="outline" onClick={logout} className="w-32">
              Cerrar Sesión
            </Button>
          </div>
        </header>

        {/* ... (Aquí van tus tarjetas de Bankroll y ROI igual que antes) ... */}
        <main className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {/* Tarjeta Bankroll */}
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
          {/* ... otras tarjetas ... */}
        </main>

        {/* Modal de Apuestas */}
        <PlaceBetModal
          isOpen={isModalOpen}
          onClose={() => setIsModalOpen(false)}
          onSuccess={() => setIsModalOpen(false)}
        />
      </div>
    </div>
  );
};
