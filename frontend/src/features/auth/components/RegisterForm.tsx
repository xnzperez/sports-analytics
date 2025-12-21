import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { toast } from "sonner";
import { registerSchema } from "../auth.schemas";
import type { RegisterFormData } from "../auth.schemas"; // Nota el 'type'
import { registerUser } from "../services/auth.service";
import { Input } from "../../../components/ui/Input";
import { Button } from "../../../components/ui/Button";
import { PasswordRequirements } from "./PasswordRequirements";

export const RegisterForm = () => {
  const {
    register,
    handleSubmit,
    watch, // Importante: para leer el valor en tiempo real
    formState: { errors, isSubmitting },
  } = useForm<RegisterFormData>({
    resolver: zodResolver(registerSchema),
    mode: "onChange", // Valida mientras escribes para mejor feedback
  });

  // Observamos el campo password
  const passwordValue = watch("password");

  const onSubmit = async (data: RegisterFormData) => {
    try {
      await registerUser(data);
      toast.success("¡Cuenta creada exitosamente!");
    } catch (error: any) {
      console.error(error);

      // MEJORA: Obtener el mensaje exacto del backend si existe
      // A veces viene en error.response.data.error o error.response.data.message
      const errorMsg =
        error.response?.data?.error ||
        error.response?.data?.message ||
        "Error al comunicarse con el servidor";

      toast.error("Error al registrarse", {
        description: errorMsg,
      });
    }
  };

  return (
    <div className="w-full max-w-md p-8 space-y-6 bg-slate-900/50 border border-slate-800 rounded-xl shadow-2xl backdrop-blur-sm">
      <div className="text-center space-y-2">
        <h2 className="text-2xl font-bold text-white tracking-tight">
          Crear Cuenta
        </h2>
        <p className="text-slate-400 text-sm">
          Únete a la plataforma profesional de analíticas.
        </p>
      </div>

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
        <Input
          label="Correo Electrónico"
          type="email"
          placeholder="ejemplo@correo.com"
          error={errors.email?.message}
          {...register("email")}
        />

        <div>
          <Input
            label="Contraseña"
            type="password"
            placeholder="••••••••"
            // Ya no mostramos el error de regex aquí para no saturar,
            // pero sí mostramos si está vacío o errores graves
            error={
              errors.password?.message && !passwordValue
                ? errors.password?.message
                : undefined
            }
            {...register("password")}
          />

          {/* Aquí inyectamos el checklist visual */}
          <PasswordRequirements password={passwordValue} />
        </div>

        <Input
          label="Confirmar Contraseña"
          type="password"
          placeholder="••••••••"
          error={errors.confirmPassword?.message}
          {...register("confirmPassword")}
        />

        <div className="pt-2">
          <Button type="submit" isLoading={isSubmitting}>
            Registrarse
          </Button>
        </div>
      </form>

      <p className="text-center text-xs text-slate-500 mt-4">
        Al registrarte aceptas nuestros Términos de Servicio y Política de
        Privacidad.
      </p>
    </div>
  );
};
