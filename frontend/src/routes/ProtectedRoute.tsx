import { Navigate, Outlet } from "react-router-dom";
import { useAuthStore } from "../features/auth/auth.store";

export const ProtectedRoute = () => {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);

  // Si NO está autenticado, redirigir al Login
  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  // Si SÍ está autenticado, renderizar la ruta hija (Outlet)
  return <Outlet />;
};
