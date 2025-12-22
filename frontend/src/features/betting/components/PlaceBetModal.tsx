import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { toast } from "sonner";
import { X, Trophy } from "lucide-react";
import { createBetSchema } from "../betting.schemas";
import type { CreateBetFormData } from "../betting.schemas";
import { placeBet } from "../services/betting.service";
import { getAvailableMatches } from "../../../services/marketService";
// Importamos Match como "type" para evitar el error de Vite
import type { Match } from "../../../services/marketService";
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

  // Estado para los partidos traídos de la API
  const [matches, setMatches] = useState<Match[]>([]);
  const [isLoadingMatches, setIsLoadingMatches] = useState(false);

  // Estado para controlar qué partido y qué equipo seleccionó visualmente
  const [selectedMatchId, setSelectedMatchId] = useState<string>("");
  const [selectedTeam, setSelectedTeam] = useState<"HOME" | "AWAY" | null>(
    null
  );

  const {
    register,
    handleSubmit,
    setValue,
    watch,
    formState: { errors, isSubmitting },
    reset,
  } = useForm<CreateBetFormData>({
    resolver: zodResolver(createBetSchema(user?.bankroll || 0)),
    defaultValues: {
      stake_units: 0,
      odds: 0,
      sport_key: "esports",
    },
  });

  const stakeValue = watch("stake_units");
  const oddsValue = watch("odds");

  // 1. Cargar partidos al abrir el modal
  useEffect(() => {
    if (isOpen) {
      setIsLoadingMatches(true);
      getAvailableMatches()
        .then((data) => setMatches(data))
        .catch(() => toast.error("Error cargando partidos en vivo"))
        .finally(() => setIsLoadingMatches(false));
    }
  }, [isOpen]);

  // 2. Manejar cambio de partido en el Select
  const handleMatchChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const matchId = e.target.value;
    setSelectedMatchId(matchId);
    setSelectedTeam(null);
    setValue("odds", 0);

    // CORRECCIÓN 1: Usamos .ID (Mayúscula)
    const match = matches.find((m) => m.ID === matchId);
    if (match) {
      setValue("title", `${match.home_team} vs ${match.away_team}`);
      setValue("sport_key", match.sport_key);
    }
  };

  // 3. Manejar selección de Ganador (Botones)
  const handleTeamSelect = (team: "HOME" | "AWAY") => {
    // CORRECCIÓN 2: Usamos .ID (Mayúscula)
    const match = matches.find((m) => m.ID === selectedMatchId);
    if (!match) return;

    setSelectedTeam(team);

    const selectedOdds = team === "HOME" ? match.home_odds : match.away_odds;
    const teamName = team === "HOME" ? match.home_team : match.away_team;

    setValue("odds", selectedOdds);

    // Guardamos metadata importante
    const detailsJson = JSON.stringify({
      match_id: match.ID, // CORRECCIÓN 3: Guardamos el ID correcto
      external_id: match.external_id,
      selection: team,
      team_name: teamName,
      league: match.league,
    });

    setValue("details", detailsJson as any);
  };

  const onSubmit = async (data: CreateBetFormData) => {
    try {
      await placeBet(data);
      toast.success("¡Apuesta registrada!", {
        description: `Has apostado $${data.stake_units} a cuota ${data.odds}`,
      });
      await fetchUser();
      reset();
      setSelectedMatchId("");
      setSelectedTeam(null);
      onSuccess();
    } catch (error: any) {
      toast.error("Error al crear apuesta", {
        description: error.response?.data?.error || "Inténtalo de nuevo",
      });
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/80 backdrop-blur-sm p-4 overflow-y-auto">
      <div className="w-full max-w-lg bg-slate-900 border border-slate-800 rounded-xl shadow-2xl animate-in fade-in zoom-in duration-200">
        {/* Header */}
        <div className="flex justify-between items-center p-6 border-b border-slate-800">
          <div className="flex items-center gap-2">
            <Trophy className="text-yellow-500" size={20} />
            <h2 className="text-xl font-bold text-white">Nueva Apuesta</h2>
          </div>
          <button
            onClick={onClose}
            className="text-slate-400 hover:text-white transition-colors"
          >
            <X size={24} />
          </button>
        </div>

        {/* Body */}
        <form onSubmit={handleSubmit(onSubmit)} className="p-6 space-y-6">
          {/* SELECCIONADOR DE PARTIDOS */}
          <div className="space-y-2">
            <label className="text-sm font-medium text-slate-300 ml-1">
              Seleccionar Evento Disponible
            </label>
            <select
              className="flex h-11 w-full rounded-md border border-slate-700 bg-slate-800 px-3 py-2 text-sm text-slate-100 focus:ring-2 focus:ring-emerald-500 outline-none"
              onChange={handleMatchChange}
              value={selectedMatchId}
              disabled={isLoadingMatches}
            >
              <option value="">
                {isLoadingMatches
                  ? "Cargando eventos..."
                  : "-- Selecciona un partido --"}
              </option>
              {matches.map((match) => (
                // CORRECCIÓN 4: Usamos .ID en key y value
                <option key={match.ID} value={match.ID}>
                  [{match.sport_key.toUpperCase()}] {match.home_team} vs{" "}
                  {match.away_team}
                </option>
              ))}
            </select>
            <input type="hidden" {...register("title")} />
          </div>

          {/* SELECCIÓN DE GANADOR (SOLO SI HAY PARTIDO) */}
          {selectedMatchId && (
            <div className="space-y-3 animate-in slide-in-from-top-2 fade-in">
              <p className="text-sm text-slate-400 text-center">
                ¿Quién ganará el encuentro?
              </p>
              <div className="grid grid-cols-2 gap-4">
                {/* Botón Local */}
                <button
                  type="button"
                  onClick={() => handleTeamSelect("HOME")}
                  className={`p-4 rounded-lg border flex flex-col items-center gap-2 transition-all ${
                    selectedTeam === "HOME"
                      ? "bg-emerald-500/10 border-emerald-500 text-emerald-400"
                      : "bg-slate-800 border-slate-700 hover:border-slate-600 text-slate-300"
                  }`}
                >
                  {/* CORRECCIÓN VISUAL: Usamos .ID en el find */}
                  <span className="font-bold text-lg">
                    {matches.find((m) => m.ID === selectedMatchId)?.home_team}
                  </span>
                  <span className="text-xs bg-slate-950 px-2 py-1 rounded text-slate-400">
                    x{matches.find((m) => m.ID === selectedMatchId)?.home_odds}
                  </span>
                </button>

                {/* Botón Visitante */}
                <button
                  type="button"
                  onClick={() => handleTeamSelect("AWAY")}
                  className={`p-4 rounded-lg border flex flex-col items-center gap-2 transition-all ${
                    selectedTeam === "AWAY"
                      ? "bg-emerald-500/10 border-emerald-500 text-emerald-400"
                      : "bg-slate-800 border-slate-700 hover:border-slate-600 text-slate-300"
                  }`}
                >
                  {/* CORRECCIÓN VISUAL: Usamos .ID en el find */}
                  <span className="font-bold text-lg">
                    {matches.find((m) => m.ID === selectedMatchId)?.away_team}
                  </span>
                  <span className="text-xs bg-slate-950 px-2 py-1 rounded text-slate-400">
                    x{matches.find((m) => m.ID === selectedMatchId)?.away_odds}
                  </span>
                </button>
              </div>
            </div>
          )}

          {/* DATOS FINANCIEROS */}
          <div className="grid grid-cols-2 gap-4">
            <Input
              label="Importe ($)"
              type="number"
              step="0.01"
              placeholder="0.00"
              error={errors.stake_units?.message}
              {...register("stake_units", { valueAsNumber: true })}
            />

            <div className="space-y-2">
              <label className="text-sm font-medium text-slate-300 ml-1">
                Cuota (Odds)
              </label>
              <div className="relative">
                <input
                  type="number"
                  step="0.01"
                  readOnly={!!selectedMatchId}
                  className={`flex h-10 w-full rounded-md border border-slate-700 bg-slate-800 px-3 py-2 text-sm text-slate-100 outline-none ${
                    selectedMatchId
                      ? "opacity-70 cursor-not-allowed"
                      : "focus:ring-2 focus:ring-emerald-500"
                  }`}
                  {...register("odds", { valueAsNumber: true })}
                />
              </div>
            </div>
          </div>

          <div className="bg-slate-950/50 p-4 rounded-lg border border-slate-800 flex justify-between items-center">
            <span className="text-slate-400 text-sm">Ganancia Potencial:</span>
            <span className="text-emerald-400 font-bold text-lg">
              $
              {((Number(stakeValue) || 0) * (Number(oddsValue) || 0)).toFixed(
                2
              )}
            </span>
          </div>

          {!selectedMatchId && (
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
              </select>
            </div>
          )}

          <div className="pt-2 flex gap-3">
            <Button
              type="button"
              variant="outline"
              onClick={onClose}
              className="w-full"
            >
              Cancelar
            </Button>
            <Button
              type="submit"
              isLoading={isSubmitting}
              className="w-full bg-emerald-600 hover:bg-emerald-700 text-white"
            >
              Confirmar Apuesta
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
};
