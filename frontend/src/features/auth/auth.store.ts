import { create } from "zustand";
import { persist } from "zustand/middleware";

interface User {
  id: string;
  email: string;
  // Agregaremos más datos aquí después (ej: bankroll)
}

interface AuthState {
  token: string | null;
  user: User | null;
  isAuthenticated: boolean;
  setToken: (token: string) => void;
  logout: () => void;
}

// Creamos el store con persistencia (se guarda en localStorage solo)
export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      token: null,
      user: null,
      isAuthenticated: false,

      setToken: (token: string) => {
        set({ token, isAuthenticated: true });
        // Aquí podríamos decodificar el token para sacar datos del usuario si fuera necesario
      },

      logout: () => {
        set({ token: null, user: null, isAuthenticated: false });
        localStorage.removeItem("auth-storage"); // Limpieza profunda
      },
    }),
    {
      name: "auth-storage", // Nombre de la llave en localStorage
    }
  )
);
