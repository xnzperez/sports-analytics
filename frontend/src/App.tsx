import { Toaster } from "sonner";
import { AppRoutes } from "./routes/AppRoutes";

function App() {
  return (
    <>
      <Toaster position="top-center" richColors theme="dark" />
      <AppRoutes />
    </>
  );
}

export default App;
