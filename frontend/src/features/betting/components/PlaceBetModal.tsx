import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { toast } from "sonner";
import { X } from "lucide-react"; // Icono de cerrar
import { createBetSchema } from "../betting.schemas"; // Importamos la función (Valor)
import type { CreateBetFormData } from "../betting.schemas";
import { placeBet } from "../services/betting.service";
import { useAuthStore } from "../../auth/auth.store";
import { Button } from "../../../components/ui/Button";
import { Input } from "../../../components/ui/Input";

interface Props {
  isOpen: boolean;
  onClose: () => void;
  onSuccess: () => void;
}

export const PlaceBetModal = ({ isOpen, onClose, onSuccess }: Props) => {
  const { user, fetchUser } = useAuthStore();

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    reset,
  } = useForm<CreateBetFormData>({
    resolver: zodResolver(createBetSchema(user?.bankroll || 0)), // Pasamos el saldo actual para validar
  });

  if (!isOpen) return null;

  const onSubmit = async (data: CreateBetFormData) => {
    try {
      await placeBet(data);
      toast.success("¡Apuesta registrada!", {
        description: `Has apostado $${data.stake_units} a cuota ${data.odds}`,
      });

      // Actualizar el saldo del usuario inmediatamente
      await fetchUser();

      reset();
      onSuccess(); // Cerrar modal
    } catch (error: any) {
      toast.error("Error al crear apuesta", {
        description: error.response?.data?.error || "Inténtalo de nuevo",
      });
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/80 backdrop-blur-sm p-4">
      <div className="w-full max-w-lg bg-slate-900 border border-slate-800 rounded-xl shadow-2xl animate-in fade-in zoom-in duration-200">
        {/* Header */}
        <div className="flex justify-between items-center p-6 border-b border-slate-800">
          <h2 className="text-xl font-bold text-white">Nueva Apuesta</h2>
          <button
            onClick={onClose}
            className="text-slate-400 hover:text-white transition-colors"
          >
            <X size={24} />
          </button>
        </div>

        {/* Body */}
        <form onSubmit={handleSubmit(onSubmit)} className="p-6 space-y-4">
          <Input
            label="Título del Evento"
            placeholder="Ej: Real Madrid vs Barcelona"
            error={errors.title?.message}
            {...register("title")}
          />

          <div className="grid grid-cols-2 gap-4">
            <Input
              label="Importe ($)"
              type="number"
              step="0.01"
              placeholder="0.00"
              error={errors.stake_units?.message}
              {...register("stake_units")}
            />
            <Input
              label="Cuota (Odds)"
              type="number"
              step="0.01"
              placeholder="1.90"
              error={errors.odds?.message}
              {...register("odds")}
            />
          </div>

          <div className="space-y-2">
            <label className="text-sm font-medium text-slate-300 ml-1">
              Deporte
            </label>
            <select
              className="flex h-10 w-full rounded-md border border-slate-700 bg-slate-800 px-3 py-2 text-sm text-slate-100 focus:ring-2 focus:ring-emerald-500 outline-none"
              {...register("sport_key")}
            >
              <option value="esports">Esports (LoL, Valorant)</option>
              <option value="football">Fútbol</option>
              <option value="basketball">Basketball</option>
              <option value="tennis">Tenis</option>
              <option value="other">Otro</option>
            </select>
          </div>

          <div className="pt-4 flex gap-3">
            <Button type="button" variant="outline" onClick={onClose}>
              Cancelar
            </Button>
            <Button type="submit" isLoading={isSubmitting}>
              Confirmar Apuesta
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
};
