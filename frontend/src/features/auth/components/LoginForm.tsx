import { useNavigate } from "react-router-dom";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { toast } from "sonner";
import { loginSchema } from "../auth.schemas";
import type { LoginFormData } from "../auth.schemas";
import { loginUser } from "../services/auth.service";
import { useAuthStore } from "../auth.store"; // Importamos el store
import { Input } from "../../../components/ui/Input";
import { Button } from "../../../components/ui/Button";

export const LoginForm = () => {
  const navigate = useNavigate();
  const setToken = useAuthStore((state) => state.setToken); // Acción de Zustand

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
  });

  const onSubmit = async (data: LoginFormData) => {
    try {
      const response = await loginUser(data);

      // 1. Guardar el token en el estado global
      setToken(response.token);

      // 2. Notificar éxito
      toast.success("¡Bienvenido de vuelta!", {
        description: "Has iniciado sesión correctamente.",
      });

      // 3. Redirigir al Dashboard
      navigate("/dashboard");
    } catch (error: any) {
      console.error(error);
      const errorMsg =
        error.response?.data?.error || "Credenciales incorrectas";
      toast.error("Error de acceso", {
        description: errorMsg,
      });
    }
  };

  return (
    <div className="w-full max-w-md p-8 space-y-6 bg-slate-900/50 border border-slate-800 rounded-xl shadow-2xl backdrop-blur-sm">
      <div className="text-center space-y-2">
        <h2 className="text-2xl font-bold text-white tracking-tight">
          Iniciar Sesión
        </h2>
        <p className="text-slate-400 text-sm">Accede a tu panel de control.</p>
      </div>

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
        <Input
          label="Correo Electrónico"
          type="email"
          placeholder="tu@email.com"
          error={errors.email?.message}
          {...register("email")}
        />

        <Input
          label="Contraseña"
          type="password"
          placeholder="••••••••"
          error={errors.password?.message}
          {...register("password")}
        />

        <div className="pt-2">
          <Button type="submit" isLoading={isSubmitting}>
            Ingresar
          </Button>
        </div>
      </form>

      <div className="text-center mt-4">
        <a
          href="#"
          className="text-xs text-emerald-500 hover:text-emerald-400 transition-colors"
        >
          ¿Olvidaste tu contraseña?
        </a>
      </div>
    </div>
  );
};
