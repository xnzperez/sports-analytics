import { api } from "../../../lib/axios";
import type {
  LoginFormData,
  RegisterFormData,
  UserProfile,
} from "../auth.schemas";

export const registerUser = async (data: RegisterFormData) => {
  // Solo enviamos email y password (confirmPassword se descarta aquí)
  const response = await api.post("/auth/register", {
    email: data.email,
    password: data.password,
  });
  return response.data;
};

export const loginUser = async (data: LoginFormData) => {
  const response = await api.post("/auth/login", data);
  return response.data;
};

// Función para obtener el perfil del usuario logueado
export const getProfile = async () => {
  // Axios enviará automáticamente el token gracias a un interceptor que configuraremos en breve
  const response = await api.get<{ user: UserProfile }>("/api/me");
  return response.data.user;
};
