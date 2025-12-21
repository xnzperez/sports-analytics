import { api } from "../../../lib/axios";
import type { LoginFormData, RegisterFormData } from "../auth.schemas";

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
