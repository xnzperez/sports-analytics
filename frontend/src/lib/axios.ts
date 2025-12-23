import axios from "axios";

// Creamos una instancia dedicada para no ensuciar la global
export const api = axios.create({
  baseURL: "http://localhost:3000",
  timeout: 10000,
  headers: {
    "Content-Type": "application/json",
  },
});

// Antes de que salga la petición, le pegamos el token
api.interceptors.request.use(
  (config) => {
    // Leemos el token directo del localStorage (donde Zustand lo guardó)
    const storage = localStorage.getItem("auth-storage");
    if (storage) {
      const { state } = JSON.parse(storage);
      if (state.token) {
        config.headers.Authorization = `Bearer ${state.token}`;
      }
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Interceptor de Respuesta (Manejo global de errores)
// Esto nos permite capturar errores 401 (token vencido) o 500 en un solo lugar.
api.interceptors.response.use(
  (response) => response,
  (error) => {
    // Si el token venció (401), podríamos redirigir al login aquí
    if (error.response?.status === 401) {
      localStorage.removeItem("auth-storage");
      window.location.href = "/login";
    }
    return Promise.reject(error);
  }
);

export default api;
