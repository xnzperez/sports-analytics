import axios from "axios";

// Creamos una instancia dedicada para no ensuciar la global
export const api = axios.create({
  baseURL: "http://localhost:3000", // La dirección de tu backend en Go
  timeout: 10000, // 10 segundos de espera máxima
  headers: {
    "Content-Type": "application/json",
  },
});

// Interceptor de Respuesta (Manejo global de errores)
// Esto nos permite capturar errores 401 (token vencido) o 500 en un solo lugar.
api.interceptors.response.use(
  (response) => response,
  (error) => {
    // Aquí podríamos loguear el error a un servicio externo en el futuro
    console.error("API Error:", error.response?.data || error.message);
    return Promise.reject(error);
  }
);
