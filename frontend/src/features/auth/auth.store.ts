import { create } from "zustand";
import { persist } from "zustand/middleware";
// CORRECCIÓN AQUÍ: Agregamos 'type'
import type { UserProfile } from "./auth.schemas";
import { getProfile } from "./services/auth.service";

interface AuthState {
  token: string | null;
  user: UserProfile | null;
  isAuthenticated: boolean;

  setToken: (token: string) => void;
  fetchUser: () => Promise<void>;
  logout: () => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      token: null,
      user: null,
      isAuthenticated: false,

      setToken: (token: string) => {
        set({ token, isAuthenticated: true });
      },

      fetchUser: async () => {
        try {
          const user = await getProfile();
          set({ user });
        } catch (error) {
          console.error("Error cargando perfil", error);
          // Si falla, asumimos token inválido
          set({ token: null, user: null, isAuthenticated: false });
        }
      },

      logout: () => {
        set({ token: null, user: null, isAuthenticated: false });
        localStorage.removeItem("auth-storage");
      },
    }),
    {
      name: "auth-storage",
    }
  )
);
