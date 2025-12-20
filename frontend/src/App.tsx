import { Toaster } from "sonner";

function App() {
  return (
    <>
      {/* El Toaster invisible que espera instrucciones */}
      <Toaster position="top-center" richColors theme="dark" />

      {/* Contenido temporal */}
      <div className="min-h-screen flex flex-col items-center justify-center gap-4">
        <h1 className="text-3xl font-bold text-emerald-400">
          Sports Analytics Dashboard 🚀
        </h1>
        <p className="text-slate-400">Sistema seguro listo.</p>
      </div>
    </>
  );
}

export default App;
