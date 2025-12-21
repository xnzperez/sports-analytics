import { useAuthStore } from "../features/auth/auth.store";
import { Button } from "../components/ui/Button";

export const DashboardPage = () => {
  const logout = useAuthStore((state) => state.logout);
  const user = useAuthStore((state) => state.user); // Aún es null, pero pronto lo llenaremos

  return (
    <div className="min-h-screen bg-slate-900 text-white p-8">
      <div className="max-w-7xl mx-auto">
        <header className="flex justify-between items-center mb-10">
          <h1 className="text-3xl font-bold text-emerald-400">Dashboard</h1>
          <div className="w-32">
            <Button variant="outline" onClick={logout}>
              Cerrar Sesión
            </Button>
          </div>
        </header>

        <main className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {/* Tarjetas de ejemplo (Skeletons visuales) */}
          <div className="bg-slate-800 p-6 rounded-xl border border-slate-700 h-40">
            <h3 className="text-slate-400 font-medium">Bankroll Actual</h3>
            <p className="text-3xl font-bold mt-2 text-white">$0.00</p>
          </div>
          <div className="bg-slate-800 p-6 rounded-xl border border-slate-700 h-40">
            <h3 className="text-slate-400 font-medium">ROI Total</h3>
            <p className="text-3xl font-bold mt-2 text-emerald-400">0%</p>
          </div>
        </main>
      </div>
    </div>
  );
};
